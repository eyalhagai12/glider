package docker

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	backend "glider"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/moby/go-archive"
)

type DockerImageService struct {
	cli *client.Client
	db  *sql.DB
}

func NewDockerImageService(cli *client.Client, db *sql.DB) *DockerImageService {
	return &DockerImageService{
		cli: cli,
		db:  db,
	}
}

func (s *DockerImageService) BuildImage(ctx context.Context, img *backend.Image, path string) (*backend.Image, error) {
	tar, err := archive.TarWithOptions(path, &archive.TarOptions{})
	if err != nil {
		return nil, err
	}
	defer tar.Close()

	buildResp, err := s.cli.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Tags: []string{img.Version},
	})
	if err != nil {
		return nil, err
	}
	defer buildResp.Body.Close()

	regAuth, err := encodeAuth("username", "password")
	if err != nil {
		return nil, err
	}

	pullResp, err := s.cli.ImagePush(ctx, img.ImagePath(), image.PushOptions{
		RegistryAuth: regAuth,
	})
	if err != nil {
		return nil, err
	}
	defer pullResp.Close()

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO images (id, name, version, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`, img.ID, img.Name, img.Version)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (s *DockerImageService) PullImage(ctx context.Context, img *backend.Image) (*backend.Image, error) {
	auth, err := encodeAuth("username", "password")
	if err != nil {
		return nil, err
	}

	pullResp, err := s.cli.ImagePull(ctx, img.ImagePath(), image.PullOptions{
		RegistryAuth: auth,
	})
	if err != nil {
		return nil, err
	}
	defer pullResp.Close()

	return img, nil
}

func (s *DockerImageService) GetByID(ctx context.Context, id uuid.UUID) (*backend.Image, error) {
	var img backend.Image
	_, err := s.db.QueryContext(
		ctx,
		"SELECT id, name, version, RegistryURL FROM images",
		&img.ID, &img.Name, &img.Version, &img.RegistryURL,
	)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func encodeAuth(username, password string) (string, error) {
	authConfig := map[string]string{
		"username": username,
		"password": password,
	}

	authBytes, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(authBytes), nil
}
