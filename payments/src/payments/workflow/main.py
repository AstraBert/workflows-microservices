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
        print(f"Extracted payment:\n{pay.model_dump_json(indent=4)}")
        async with ctx.store.edit_state() as state:
            state.order_id = ev.data.get("orderId", "")
        return PaymentEvent(payment=pay)
    @step
    async def process_payment(self, ev: PaymentEvent, querier: Annotated[AsyncQuerier, Resource(get_querier)], ctx: Context) -> OutputEvent:
        state = await ctx.store.get_state()
        try:
            res = await querier.create_payment(user_id=ev.payment.user_id, status=ev.payment.status, method=ev.payment.method, amount=ev.payment.amount)
            if res is None:
                print("Payment processing failed")
                return OutputEvent(success=False, orderId=state.order_id)
            else:
                print("Payment processed successfully")
                return OutputEvent(success=True, orderId=state.order_id)
        except Exception:
            print("Payment processing failed")
            return OutputEvent(success=False, orderId=state.order_id)

async def run_workflow():
    consumer = KafkaConsumer(
        'orders',
        bootstrap_servers='kafka:9092',
        auto_offset_reset='earliest',
        enable_auto_commit=True, 
        group_id='order-processor-2',
        value_deserializer=lambda x: json.loads(x.decode('utf-8')),
        consumer_timeout_ms=1000
    )
    
    producer = KafkaProducer(
        bootstrap_servers=['kafka:9092'],
        value_serializer=lambda v: json.dumps(v).encode('utf-8')
    )
    
    wf = PaymentWorkflow()
    
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