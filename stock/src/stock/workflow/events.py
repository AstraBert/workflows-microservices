from stock.db import Stock

from pydantic import ConfigDict
from workflows.events import Event, StartEvent, StopEvent
from typing import Any, Literal

class InputEvent(StartEvent):
    data: dict[str, Any]

class UpdateStockEvent(Event):
    item: str
    available_number: int
    model_config = ConfigDict(arbitrary_types_allowed=True)

class OutputEvent(StopEvent):
    type: Literal["stock"] = "stock"
    success: bool
    orderId: str