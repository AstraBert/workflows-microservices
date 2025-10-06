-- name: GetStockItem :one
SELECT id, item, available_number, updated_at FROM stock
WHERE item = $1 LIMIT 1;

-- name: UpdateStockItem :exec
UPDATE stock set available_number = $1 WHERE item = $2;

-- name: InsertSeedData :exec
INSERT INTO stock (item, available_number) 
VALUES     
    ('T-Shirt', 100),     
    ('Socks', 100),     
    ('Mug', 100)
ON CONFLICT DO NOTHING;