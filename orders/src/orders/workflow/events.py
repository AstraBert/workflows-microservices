from orders.db import Order

from pydantic import ConfigDict
from workflows.events import Event, StartEvent, StopEvent
from typing import Any, Literal

class InputEvent(StartEvent):
    data: dict[str, Any]

class OrderEvent(Event):
    order: Order
    model_config = ConfigDict(arbitrary_types_allowed=True)

class OutputEvent(Event):
    type: Literal["order"] = "order"
    success: bool
    orderId: str