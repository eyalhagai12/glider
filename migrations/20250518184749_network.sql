-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS networks (
    id UUID PRIMARY KEY,
    interface_name VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    port VARCHAR(5) NOT NULL,
    project_id UUID REFERENCES projects(id) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS networks;
-- +goose StatementEnd
