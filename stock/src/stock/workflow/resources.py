import os
from stock.db import AsyncQuerier

from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy import text

DATABASE_URL=f"postgresql+asyncpg://{os.getenv('POSTGRES_USER')}:{os.getenv('POSTGRES_PASSWORD')}@postgres:5432/{os.getenv('POSTGRES_DB')}"

QUERY = """
CREATE TABLE IF NOT EXISTS stock (
    id SERIAL PRIMARY KEY,
    item TEXT NOT NULL UNIQUE,
    available_number INTEGER NOT NULL DEFAULT 0 CHECK(available_number >= 0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
"""

engine = create_async_engine(
    DATABASE_URL,
)

async def get_querier(*args, **kwargs) -> AsyncQuerier:
    conn = await engine.connect()
    await conn.execute(text(QUERY))
    querier = AsyncQuerier(conn)
    await querier.insert_seed_data()
    return querier
