-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS cart_item (
    cart_id BIGINT NOT NULL,
    sku BIGINT NOT NULL,
    "count" INT NOT NULL,
    PRIMARY KEY (cart_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_item;
DROP TABLE IF EXISTS cart;
-- +goose StatementEnd
