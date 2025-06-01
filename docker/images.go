package docker

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	backend "glider"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/moby/go-archive"
)

type DockerImageService struct {
	cli    *client.Client
	db     *sql.DB
	logger *slog.Logger
}

func NewDockerImageService(cli *client.Client, db *sql.DB, logger *slog.Logger) *DockerImageService {
	return &DockerImageService{
		cli:    cli,
		db:     db,
		logger: logger,
	}
}

func (s *DockerImageService) BuildImage(ctx context.Context, img *backend.Image, path string) (*backend.Image, error) {
	tar, err := archive.TarWithOptions(path, &archive.TarOptions{})
	if err != nil {
		return nil, err
	}
	defer tar.Close()

	buildResp, err := s.cli.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Tags: []string{img.ImagePath()},
	})
	if err != nil {
		return nil, err
	}
	// io.Copy(os.Stdout, buildResp.Body)
	defer buildResp.Body.Close()

	regAuth, err := encodeAuth("username", "password")
	if err != nil {
		return nil, err
	}

	pushResp, err := s.cli.ImagePush(ctx, img.ImagePath(), image.PushOptions{
		RegistryAuth: regAuth,
	})
	if err != nil {
		return nil, err
	}
	// io.Copy(os.Stdout, pushResp)
	defer pushResp.Close()

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO images (id, name, version, registry_url)
		VALUES ($1, $2, $3, $4)
	`, img.ID, img.Name, img.Version, img.RegistryURL)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (s *DockerImageService) PullImage(ctx context.Context, img *backend.Image) (*backend.Image, error) {
	logger := s.logger.With("image_path", img.ImagePath())

	auth, err := encodeAuth("username", "password")
	if err != nil {
		return nil, err
	}

	logger.Debug("pulling image", slog.String("path", img.ImagePath()))
	pullResp, err := s.cli.ImagePull(ctx, img.ImagePath(), image.PullOptions{
		RegistryAuth: auth,
	})
	if err != nil {
		return nil, err
	}
	// io.Copy(os.Stdout, pullResp)
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
