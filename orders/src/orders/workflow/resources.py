from orders.db import AsyncQuerier

from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy import text

DATABASE_URL = "sqlite+aiosqlite:///./orders.db"

QUERY = """
CREATE TABLE IF NOT EXISTS orders (
    order_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    order_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL CHECK(status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),
    shipping_address TEXT NOT NULL
);
"""

engine = create_async_engine(
    DATABASE_URL
)

async def get_querier(*args, **kwargs) -> AsyncQuerier:
    conn = await engine.connect()
    await conn.execute(text(QUERY))
    return AsyncQuerier(conn)
