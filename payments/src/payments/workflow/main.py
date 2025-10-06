import time
import json
import asyncio

from typing import Annotated
from payments.db import Payment, AsyncQuerier
from .events import InputEvent, OutputEvent, PaymentEvent
from .resources import get_querier
from workflows import Workflow, step, Context
from workflows.resource import Resource
from kafka import KafkaConsumer, KafkaProducer

class PaymentWorkflow(Workflow):
    @step
    async def extract_payment(self, ev: InputEvent, ctx: Context) -> PaymentEvent:
        pay = Payment(payment_id=0, payment_time=time.time(), status="completed", method=ev.data.get("paymentMethod", ""), amount=float(ev.data.get("amount", "25.00")), user_id=ev.data.get("userId", ""))
        async with ctx.store.edit_state() as state:
            state.order_id = ev.data.get("orderId", "")
        return PaymentEvent(payment=pay)
    @step
    async def process_payment(self, ev: PaymentEvent, querier: Annotated[AsyncQuerier, Resource(get_querier)], ctx: Context) -> OutputEvent:
        state = await ctx.store.get_state()
        try:
            res = await querier.create_payment(user_id=ev.payment.user_id, status=ev.payment.status, method=ev.payment.method, amount=ev.payment.amount)
            if res is None:
                return OutputEvent(success=False, orderId=state.order_id)
            else:
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