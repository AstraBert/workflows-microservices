import os
from orders.db import AsyncQuerier

from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy import text

DATABASE_URL=f"postgresql+asyncpg://{os.getenv('POSTGRES_USER')}:{os.getenv('POSTGRES_PASSWORD')}@postgres:5432/{os.getenv('POSTGRES_DB')}"

QUERY = """
CREATE TABLE IF NOT EXISTS orders (
    order_id TEXT UNIQUE NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    order_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
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
