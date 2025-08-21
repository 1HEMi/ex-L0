
CREATE TABLE IF NOT EXISTS orders (
  order_uid         TEXT PRIMARY KEY,
  track_number      TEXT NOT NULL,
  entry             TEXT NOT NULL,
  locale            TEXT,
  internal_signature TEXT,
  customer_id       TEXT,
  delivery_service  TEXT,
  shardkey          TEXT,
  sm_id             INTEGER,
  date_created      TIMESTAMPTZ NOT NULL,
  oof_shard         TEXT
);

CREATE TABLE IF NOT EXISTS delivery (
  order_uid  TEXT PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
  name       TEXT,
  phone      TEXT,
  zip        TEXT,
  city       TEXT,
  address    TEXT,
  region     TEXT,
  email      TEXT
);

CREATE TABLE IF NOT EXISTS payment (
  order_uid     TEXT PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
  transaction   TEXT UNIQUE NOT NULL,
  request_id    TEXT,
  currency      TEXT,
  provider      TEXT,
  amount        INTEGER NOT NULL,
  payment_dt    TIMESTAMPTZ NOT NULL,
  bank          TEXT,
  delivery_cost INTEGER,
  goods_total   INTEGER,
  custom_fee    INTEGER
);

CREATE TABLE IF NOT EXISTS items (
  id           BIGSERIAL PRIMARY KEY,
  order_uid    TEXT NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
  chrt_id      INTEGER,
  track_number TEXT,
  price        INTEGER,
  rid          TEXT,
  name         TEXT,
  sale         INTEGER,
  size         TEXT,
  total_price  INTEGER,
  nm_id        INTEGER,
  brand        TEXT,
  status       INTEGER
);

CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);
CREATE INDEX IF NOT EXISTS idx_orders_date_created ON orders(date_created DESC);
