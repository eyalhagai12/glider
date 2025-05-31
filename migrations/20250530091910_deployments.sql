-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS deployments (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    version VARCHAR(100) NOT NULL,
    environment VARCHAR(100) NOT NULL,
    project_id UUID NOT NULL,
    status VARCHAR(100) NOT NULL,
    image_id UUID NOT NULL,
    replicas INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (project_id) REFERENCES projects (id),
    FOREIGN KEY (image_id) REFERENCES images (ID)
);

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    deployment_id UUID NOT NULL,
    is_system BOOLEAN NOT NULL,
    FOREIGN KEY (deployment_id) REFERENCES deployments (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS deployments;
-- +goose StatementEnd
