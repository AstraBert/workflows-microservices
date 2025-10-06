import time
import json
import asyncio

from typing import Annotated, Union
from stock.db import AsyncQuerier
from .events import InputEvent, OutputEvent, UpdateStockEvent
from .resources import get_querier
from workflows import Workflow, step, Context
from workflows.resource import Resource
from kafka import KafkaConsumer, KafkaProducer

PRICE_TO_ITEM = {"25": "Mug", "50": "Socks", "100": "T-Shirt"}

class StockWorkflow(Workflow):
    @step
    async def extract_payment(self, ev: InputEvent, ctx: Context, querier: Annotated[AsyncQuerier, Resource(get_querier)]) -> Union[UpdateStockEvent, OutputEvent]:
        async with ctx.store.edit_state() as state:
            state.order_id = ev.data.get("orderId", "")
        item = PRICE_TO_ITEM.get(ev.data.get("amount", "25"), "Mug")
        print(f"Chosen item: {item}")
        stock_item = await querier.get_stock_item(item=item)
        if not stock_item:
            print(f"It was not possible to retrieve data about {item} from the database")
            return OutputEvent(success=False, orderId=state.order_id)
        else:
            number = stock_item.available_number - 1
            print(f"Updating available number of {item} to {number}")
            return UpdateStockEvent(item=item, available_number=number)
    @step
    async def process_payment(self, ev: UpdateStockEvent, querier: Annotated[AsyncQuerier, Resource(get_querier)], ctx: Context) -> OutputEvent:
        state = await ctx.store.get_state()
        try:
            await querier.update_stock_item(item=ev.item, available_number=ev.available_number)
            print(f"Available number of {ev.item} updated successfully")
            return OutputEvent(success=True, orderId=state.order_id)
        except Exception:
            print(f"Failed to update available number of {ev.item}")
            return OutputEvent(success=False, orderId=state.order_id)

async def run_workflow():
    consumer = KafkaConsumer(
        'orders',
        bootstrap_servers='kafka:9092',
        auto_offset_reset='earliest',
        enable_auto_commit=True, 
        group_id='order-processor-1',
        value_deserializer=lambda x: json.loads(x.decode('utf-8')),
        consumer_timeout_ms=1000
    )
    
    producer = KafkaProducer(
        bootstrap_servers=['kafka:9092'],
        value_serializer=lambda v: json.dumps(v).encode('utf-8')
    )
    
    wf = StockWorkflow()
    
    print("Starting to consume messages from 'orders' topic...")
    
    try:
        while True:
            # Poll for messages with timeout
            message_batch = consumer.poll(timeout_ms=1000)
            
            for _, messages in message_batch.items():
                for message in messages:
                    print(f"Received message: {message.value}")
                    try:
                        output = await wf.run(start_event=InputEvent(data=message.value))
                        producer.send(topic="order-status", value=output.model_dump())
                        producer.flush()  # Ensure message is sent
                    except Exception as e:
                        print(f"Error processing message: {e}")
            
            # Small sleep to prevent CPU spinning
            await asyncio.sleep(0.1)
            
    except KeyboardInterrupt:
        print("Shutting down consumer...")
    finally:
        consumer.close()
        producer.close()

def main():
    asyncio.run(run_workflow())