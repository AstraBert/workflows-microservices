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

class PaymentWorkflow(Workflow):
    @step
    async def extract_payment(self, ev: InputEvent, ctx: Context, querier: Annotated[AsyncQuerier, Resource(get_querier)]) -> Union[UpdateStockEvent, OutputEvent]:
        async with ctx.store.edit_state() as state:
            state.order_id = ev.data.get("orderId", "")
        item = PRICE_TO_ITEM.get(ev.data.get("amount", "25"), "Mug")
        stock_item = await querier.get_stock_item(item=item)
        if not stock_item:
            return OutputEvent(success=False, orderId=state.order_id)
        else:
            number = stock_item.available_number - 1
            return UpdateStockEvent(item=item, available_number=number)
    @step
    async def process_payment(self, ev: UpdateStockEvent, querier: Annotated[AsyncQuerier, Resource(get_querier)], ctx: Context) -> OutputEvent:
        state = await ctx.store.get_state()
        try:
            await querier.update_stock_item(item=ev.item, available_number=ev.available_number)
            return OutputEvent(success=True, orderId=state.order_id)
        except Exception as e:
            return OutputEvent(success=False, orderId=state.order_id)

async def run_workflow():
    consumer = KafkaConsumer(
        bootstrap_servers='localhost:9092',
        auto_offset_reset='earliest',
        value_deserializer=lambda x: json.loads(x.decode('utf-8'))
    )
    producer = KafkaProducer(
        bootstrap_servers=['localhost:9092'],
        value_serializer=lambda v: json.dumps(v).encode('utf-8')  # Serialize data to JSON
    )
    consumer.subscribe(["orders"])
    wf = PaymentWorkflow()
    while True:
        for message in consumer:
            output = await wf.run(start_event=InputEvent(data=message.value))
            producer.send(topic="order-status", value=output.model_dump())

def main():
    asyncio.run(run_workflow())