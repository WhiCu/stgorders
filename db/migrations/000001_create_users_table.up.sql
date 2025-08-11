create table orders (
    id bigint primary key generated always as identity,
    order_uid text unique,
    track_number text,
    entry text,
    locale text,
    internal_signature text,
    customer_id text,
    delivery_service text,
    shardkey text,
    sm_id int,
    date_created timestamp,
    oof_shard text
);

create table delivery (
    id bigint primary key generated always as identity,
    order_id bigint references orders (id),
    name text,
    phone text,
    zip text,
    city text,
    address text,
    region text,
    email text
);

create table payment (
    id bigint primary key generated always as identity,
    order_id bigint references orders (id),
    transaction text,
    request_id text,
    currency text,
    provider text,
    amount int,
    payment_dt bigint,
    bank text,
    delivery_cost int,
    goods_total int,
    custom_fee int
);

create table items (
    id bigint primary key generated always as identity,
    order_id bigint references orders (id),
    chrt_id bigint,
    track_number text,
    price int,
    rid text,
    name text,
    sale int,
    size text,
    total_price int,
    nm_id bigint,
    brand text,
    status int
);