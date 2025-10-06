from stock.db import AsyncQuerier

from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy import text

DATABASE_URL = "sqlite+aiosqlite:///./stock.db"

QUERY = """
CREATE TABLE IF NOT EXISTS stock (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item TEXT NOT NULL UNIQUE,
    available_number INTEGER NOT NULL DEFAULT 0 CHECK(available_number >= 0),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
"""

engine = create_async_engine(
    DATABASE_URL
)

async def get_querier(*args, **kwargs) -> AsyncQuerier:
    conn = await engine.connect()
    await conn.execute(text(QUERY))
    querier = AsyncQuerier(conn)
    await querier.insert_seed_data()
    return querier
