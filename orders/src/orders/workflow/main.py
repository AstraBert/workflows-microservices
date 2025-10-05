import time
import json
import asyncio

from typing import Annotated
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
        ordr = Order(order_id=ev.data.get("orderId", ""), user_id=ev.data.get("userId", ""), user_name=user_name, shipping_address=address, order_time=time.time(), status="processing", email=ev.data.get("email", ""), phone=ev.data.get("phone", ""))
        return OrderEvent(order=ordr)
    @step
    async def process_order(self, ev: OrderEvent, querier: Annotated[AsyncQuerier, Resource(get_querier)]) -> OutputEvent:
        try:
            params = CreateOrderParams(user_id=ev.order.user_id, user_name=ev.order.user_name, email=ev.order.email, phone=ev.order.phone, shipping_address=ev.order.shipping_address, status=ev.order.status)
            res = await querier.create_order(params)
            if res is None:
                return OutputEvent(success=False, orderId=ev.order.order_id)
            else:
                return OutputEvent(success=True, orderId=ev.order.order_id)
        except Exception as e:
            return OutputEvent(success=False, orderId=ev.order.order_id)

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
    wf = OrderWorkflow()
    while True:
        for message in consumer:
            output = await wf.run(start_event=InputEvent(data=message.value))
            producer.send(topic="order-status", value=output.model_dump())

def main():
    asyncio.run(run_workflow())