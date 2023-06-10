-- +goose Up
-- +goose StatementBegin

CREATE TYPE order_status AS ENUM (
    'new',
    'payed',
    'failed',
    'awaiting payment',
    'cancelled'
);

CREATE TABLE IF NOT EXISTS user_order (
    id SERIAL PRIMARY KEY,
    "user_id" BIGINT NOT NULL,
    status order_status NOT NULL
);

CREATE TABLE IF NOT EXISTS order_item (
    order_id BIGINT NOT NULL,
    sku BIGINT NOT NULL,
    "count" INT NOT NULL,
    PRIMARY KEY (order_id, sku)
);

CREATE TABLE IF NOT EXISTS stock (
    warehouse_id BIGINT NOT NULL,
    sku BIGINT NOT NULL,
    "count" INT NOT NULL,
    PRIMARY KEY (warehouse_id, sku)
);


CREATE TABLE IF NOT EXISTS reservation_stock (
    order_id BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    sku BIGINT NOT NULL,
    "count" INT NOT NULL,
    PRIMARY KEY (order_id, warehouse_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
