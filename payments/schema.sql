-- Payments table
CREATE TABLE payments (
    payment_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    payment_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL CHECK(status IN ('pending', 'completed', 'failed', 'refunded')),
    method TEXT NOT NULL CHECK(method IN ('credit_card', 'debit_card', 'paypal', 'bank_transfer')),
    amount DECIMAL(10, 2)
);
