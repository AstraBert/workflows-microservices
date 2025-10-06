from payments.db import Payment

from pydantic import ConfigDict
from workflows.events import Event, StartEvent, StopEvent
from typing import Any, Literal

class InputEvent(StartEvent):
    data: dict[str, Any]

class PaymentEvent(Event):
    payment: Payment
    model_config = ConfigDict(arbitrary_types_allowed=True)

class OutputEvent(StopEvent):
    type: Literal["payment"] = "payment"
    success: bool
    orderId: str