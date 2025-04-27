-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE deployments (
    id UUID PRIMARY KEY,
    name VARCHAR(1500) NOT NULL,
    github_repo VARCHAR(8000) NOT NULL,
    github_branch VARCHAR(8000) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    target_replica_count INT NOT NULL,
    replica_count INT NOT NULL,
    status VARCHAR(25) NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE deployments;
-- +goose StatementEnd
