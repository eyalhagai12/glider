-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS containers (
    id UUID PRIMARY KEY,
    platform_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    image_id UUID NOT NULL,
    deployment_id UUID NOT NULL,
    Host VARCHAR(255) NOT NULL,
    Port VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS containers;
-- +goose StatementEnd
