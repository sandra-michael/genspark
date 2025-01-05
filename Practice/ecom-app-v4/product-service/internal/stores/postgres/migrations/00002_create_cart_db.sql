-- +goose Up
-- +goose StatementBegin
CREATE TYPE status_enum AS ENUM ('inprogress', 'pending', 'completed');

CREATE TABLE IF NOT EXISTS cart (
    id UUID PRIMARY KEY, -- Unique identifier for the product
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE, -- Foreign key referencing products table
	user_id UUID NOT NULL,
	order_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 1), -- quantity level, must be non-negative
	status status_enum DEFAULT 'inprogress', -- ENUM field with default value
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart;

DROP TYPE status_enum;
-- +goose StatementEnd
