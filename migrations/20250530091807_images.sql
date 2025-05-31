-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    version VARCHAR(100) NOT NULL,
    RegistryURL VARCHAR(100)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
