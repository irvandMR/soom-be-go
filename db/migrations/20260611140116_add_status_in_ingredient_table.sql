-- +goose Up
ALTER TABLE ingredients 
ADD COLUMN status VARCHAR(50) NULL ;

-- +goose Down
ALTER TABLE ingredients 
DROP COLUMN status;