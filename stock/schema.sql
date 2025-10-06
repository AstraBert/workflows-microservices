-- Stock table
CREATE TABLE stock (
    id SERIAL PRIMARY KEY,
    item TEXT NOT NULL UNIQUE,
    available_number INTEGER NOT NULL DEFAULT 0 CHECK(available_number >= 0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);