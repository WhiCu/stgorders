-- =======================
-- ORDER
-- =======================

-- name: ListOrders :many
SELECT * FROM orders;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
    order_uid,
    track_number,
    entry,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id;

-- name: UpdateOrder :exec
UPDATE orders
SET
    track_number = $2,
    entry = $3,
    locale = $4,
    internal_signature = $5,
    customer_id = $6,
    delivery_service = $7,
    shardkey = $8,
    sm_id = $9,
    date_created = $10,
    oof_shard = $11
WHERE id = $1;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- =======================
-- DELIVERY
-- =======================

-- name: ListDeliveries :many
SELECT * FROM delivery;

-- name: GetDeliveryByID :one
SELECT * FROM delivery
WHERE id = $1;

-- name: CreateDelivery :one
INSERT INTO delivery (
    order_id,
    name,
    phone,
    zip,
    city,
    address,
    region,
    email
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id;

-- name: UpdateDelivery :exec
UPDATE delivery
SET
    name = $2,
    phone = $3,
    zip = $4,
    city = $5,
    address = $6,
    region = $7,
    email = $8
WHERE id = $1;

-- name: DeleteDelivery :exec
DELETE FROM delivery
WHERE id = $1;

-- =======================
-- PAYMENT
-- =======================

-- name: ListPayments :many
SELECT * FROM payment;

-- name: GetPaymentByID :one
SELECT * FROM payment
WHERE id = $1;

-- name: CreatePayment :one
INSERT INTO payment (
    order_id,
    transaction,
    request_id,
    currency,
    provider,
    amount,
    payment_dt,
    bank,
    delivery_cost,
    goods_total,
    custom_fee
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id;

-- name: UpdatePayment :exec
UPDATE payment
SET
    transaction = $2,
    request_id = $3,
    currency = $4,
    provider = $5,
    amount = $6,
    payment_dt = $7,
    bank = $8,
    delivery_cost = $9,
    goods_total = $10,
    custom_fee = $11
WHERE id = $1;

-- name: DeletePayment :exec
DELETE FROM payment
WHERE id = $1;

-- =======================
-- ITEMS
-- =======================

-- name: ListItems :many
SELECT * FROM items;

-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1;

-- name: CreateItem :one
INSERT INTO items (
    order_id,
    chrt_id,
    track_number,
    price,
    rid,
    name,
    sale,
    size,
    total_price,
    nm_id,
    brand,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING id;

-- name: UpdateItem :exec
UPDATE items
SET
    chrt_id = $2,
    track_number = $3,
    price = $4,
    rid = $5,
    name = $6,
    sale = $7,
    size = $8,
    total_price = $9,
    nm_id = $10,
    brand = $11,
    status = $12
WHERE id = $1;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = $1;

-- =======================
-- CREATE FULL ORDER
-- =======================

