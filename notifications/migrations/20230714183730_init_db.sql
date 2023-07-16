-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM (
    'new',
    'payed',
    'failed',
    'awaiting payment',
    'cancelled'
);

CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    status order_status NOT NULL,
    message CHAR(255) NOT NULL,
    created_date TIMESTAMP DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message;
-- +goose StatementEnd
