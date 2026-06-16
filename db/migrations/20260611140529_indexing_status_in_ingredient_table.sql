-- +goose NO TRANSACTION

-- +goose Up
CREATE INDEX CONCURRENTLY idx_ingredients_status 
ON ingredients(status);

-- +goose Down
DROP INDEX CONCURRENTLY idx_ingredients_status;