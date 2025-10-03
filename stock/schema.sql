-- Stock table
CREATE TABLE stock (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item TEXT NOT NULL UNIQUE,
    available_number INTEGER NOT NULL DEFAULT 0 CHECK(available_number >= 0),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);