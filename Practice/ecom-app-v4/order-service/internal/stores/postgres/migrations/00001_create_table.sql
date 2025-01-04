-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS orders (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,                          -- UUID of the user placing the order
    product_id UUID NOT NULL,                       -- UUID of the product
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'paid', 'canceled')), -- Status with constraints
    stripe_transaction_id TEXT,                      -- Stripe unique transaction ID
    total_price BIGINT NOT NULL,                    -- Total price in cents (bigint for large amounts)
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;

-- +goose StatementEnd
