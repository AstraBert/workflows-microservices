import json
import asyncio

from typing import Annotated
from datetime import datetime
from orders.db import CreateOrderParams, Order, AsyncQuerier
from .events import InputEvent, OutputEvent, OrderEvent
from .resources import get_querier
from workflows import Workflow, step
from workflows.resource import Resource
from kafka import KafkaConsumer, KafkaProducer

class OrderWorkflow(Workflow):
    @step
    async def extract_order(self, ev: InputEvent) -> OrderEvent:
        user_name = ev.data.get("firstName", "") + " " + ev.data.get("lastName", "")
        address = ev.data.get("address", "") + " " + ev.data.get("address2") + ", " + ev.data.get("city", "") + " (" + ev.data.get("zip", "") + "), " + ev.data.get("state", "") + ", " + ev.data.get("country", "")
        ordr = Order(order_id=ev.data.get("orderId", ""), user_id=str(ev.data.get("userId", "")), user_name=user_name, shipping_address=address, order_time=datetime.now(), status="processing", email=ev.data.get("email", ""), phone=ev.data.get("phone", ""))
        print(f"Extracted order:\n{ordr.model_dump_json(indent=4)}")
        return OrderEvent(order=ordr)
    @step
    async def process_order(self, ev: OrderEvent, querier: Annotated[AsyncQuerier, Resource(get_querier)]) -> OutputEvent:
        try:
            params = CreateOrderParams(user_id=str(ev.order.user_id), user_name=ev.order.user_name, email=ev.order.email, phone=ev.order.phone, shipping_address=ev.order.shipping_address, status=ev.order.status, order_id=ev.order.order_id)
            res = await querier.create_order(params)
            if res is None:
                print("Order creation failed because of: no rows returned")
                return OutputEvent(success=False, orderId=ev.order.order_id)
            else:
                print("Successfully created the order")
                return OutputEvent(success=True, orderId=ev.order.order_id)
        except Exception as e:
            print(f"Order creation failed because of: {e}")
            return OutputEvent(success=False, orderId=ev.order.order_id)

async def run_workflow():
    consumer = KafkaConsumer(
        'orders',
        bootstrap_servers='kafka:9092',
        auto_offset_reset='earliest',
        enable_auto_commit=True, 
        group_id='order-processor-3',
        value_deserializer=lambda x: json.loads(x.decode('utf-8')),
        consumer_timeout_ms=1000
    )
    
    producer = KafkaProducer(
        bootstrap_servers=['kafka:9092'],
        value_serializer=lambda v: json.dumps(v).encode('utf-8')
    )
    
    wf = OrderWorkflow()
    
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