-- =======================
-- Создание таблиц
-- =======================

create table orders (
  id bigint primary key generated always as identity,
  order_uid text unique not null,
  track_number text not null,
  entry text not null,
  locale text not null,
  internal_signature text not null,
  customer_id text not null,
  delivery_service text not null,
  shardkey text not null,
  sm_id int not null,
  date_created timestamp not null,
  oof_shard text not null
);

create table delivery (
  id bigint primary key generated always as identity,
  name text not null,
  phone text not null,
  zip text not null,
  city text not null,
  address text not null,
  region text not null,
  email text not null
);

create table payment (
  id bigint primary key generated always as identity,
  transaction text not null,
  request_id text not null,
  currency text not null,
  provider text not null,
  amount int not null,
  payment_dt bigint not null,
  bank text not null,
  delivery_cost int not null,
  goods_total int not null,
  custom_fee int not null
);

create table items (
  id bigint primary key generated always as identity,
  chrt_id bigint not null,
  track_number text not null,
  price int not null,
  rid text not null,
  name text not null,
  sale int not null,
  size text not null,
  total_price int not null,
  nm_id bigint not null,
  brand text not null,
  status int not null
);

-- =======================
-- Тестовые данные
-- =======================

insert into orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
values
('order_001', 'TRACK123', 'web', 'en', 'sig1', 'cust001', 'dhl', '1', 100, now(), 'shardA'),
('order_002', 'TRACK456', 'app', 'ru', 'sig2', 'cust002', 'fedex', '2', 200, now(), 'shardB');

insert into delivery (name, phone, zip, city, address, region, email)
values
('John Doe', '+123456789', '12345', 'New York', '123 Main St', 'NY', 'john@example.com'),
('Иван Иванов', '+79991234567', '101000', 'Москва', 'ул. Ленина, 5', 'Московская область', 'ivan@example.com');

insert into payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
values
('order_001', 'req001', 'USD', 'paypal', 1500, 1692000000, 'Chase', 100, 1400, 0),
('order_002', 'req002', 'RUB', 'yoomoney', 2000, 1692100000, 'Sberbank', 200, 1800, 0);

insert into items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
values
(111, 'TRACK123', 500, 'RID001', 'T-Shirt', 10, 'M', 450, 10001, 'Nike', 1),
(112, 'TRACK123', 1000, 'RID002', 'Shoes', 5, '42', 950, 10002, 'Adidas', 1),
(113, 'TRACK456', 2000, 'RID003', 'Куртка', 15, 'L', 1700, 10003, 'Puma', 1);
