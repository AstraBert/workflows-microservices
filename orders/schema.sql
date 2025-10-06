-- Orders table
CREATE TABLE orders (
    order_id TEXT UNIQUE NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    order_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL CHECK(status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),
    shipping_address TEXT NOT NULL
);
