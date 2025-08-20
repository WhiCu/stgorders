-- =======================
-- ORDER
-- =======================

-- name: ListOrders :many
SELECT * FROM orders;

-- name: GetOrderByUID :one
SELECT *
FROM orders
WHERE order_uid = $1;

-- name: GetLastOrders :many
SELECT *
FROM orders
ORDER BY id DESC
LIMIT $1;


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

-- =======================
-- DELIVERY
-- =======================

-- name: ListDeliveries :many
SELECT * FROM delivery;

-- name: GetDeliveryByOrderUID :one
SELECT *
FROM delivery
WHERE order_id = $1;

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

-- =======================
-- PAYMENT
-- =======================

-- name: ListPayments :many
SELECT * FROM payment;

-- name: GetPaymentByTransaction :one
SELECT *
FROM payment
WHERE transaction = $1;


-- name: CreatePayment :one
INSERT INTO payment (
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
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id;

-- =======================
-- ITEMS
-- =======================

-- name: ListItems :many
SELECT * FROM items;

-- name: GetItemsByTrackNumber :many
SELECT *
FROM items
WHERE track_number = $1;

-- name: CreateItem :one
INSERT INTO items (
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
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id;
