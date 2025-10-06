-- name: UpdateStockItem :exec
UPDATE stock
set available_number = ?
WHERE item = ?;

-- name: GetStockItem :one
SELECT * FROM stock
WHERE item = ? LIMIT 1;

-- name: InsertSeedData :exec
INSERT OR IGNORE INTO stock (item, available_number) VALUES
    ('T-Shirt', 100),
    ('Socks', 100),
    ('Mug', 100);