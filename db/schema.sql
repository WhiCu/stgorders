
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
  order_id bigint references orders (id) not null,
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
  order_id bigint references orders (id) not null,
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
  order_id bigint references orders (id) not null,
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