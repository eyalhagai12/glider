-- +goose Up
-- +goose StatementBegin
CREATE TABLE nodes (
    id UUID PRIMARY KEY,
    deployment_url VARCHAR(8000) NOT NULL,
    metrics_url VARCHAR(8000) NOT NULL,
    health_url VARCHAR(8000) NOT NULL,
    connection_url VARCHAR(8000) NOT NULL,
    status VARCHAR(8000) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE nodes;
-- +goose StatementEnd
