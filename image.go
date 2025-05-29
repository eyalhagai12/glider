package backend

import (
	"context"

	"github.com/google/uuid"
)

type Image struct {
	ID          uuid.UUID
	Name        string
	Version     string
	RegistryURL string
	Path        string
}

func (i *Image) ImageName() string {
	return i.Name + ":" + i.Version
}

func (i *Image) ImagePath() string {
	return i.RegistryURL + "/" + i.ImageName()
}

type ImageService interface {
	BuildImage(ctx context.Context, image *Image) (*Image, error)
	PullImage(ctx context.Context, image *Image) (*Image, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Image, error)
}
