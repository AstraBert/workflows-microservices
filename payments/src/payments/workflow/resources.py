import os
from payments.db import AsyncQuerier

from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy import text

DATABASE_URL=f"postgresql+asyncpg://{os.getenv('POSTGRES_USER')}:{os.getenv('POSTGRES_PASSWORD')}@postgres:5432/{os.getenv('POSTGRES_DB')}"

QUERY = """
CREATE TABLE IF NOT EXISTS payments (
    payment_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    payment_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL CHECK(status IN ('pending', 'completed', 'failed', 'refunded')),
    method TEXT NOT NULL CHECK(method IN ('credit_card', 'debit_card', 'paypal', 'bank_transfer')),
    amount DECIMAL(10, 2)
);
"""

engine = create_async_engine(
    DATABASE_URL
)

async def get_querier(*args, **kwargs) -> AsyncQuerier:
    conn = await engine.connect()
    await conn.execute(text(QUERY))
    return AsyncQuerier(conn)