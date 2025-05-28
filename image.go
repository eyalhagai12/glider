package backend

import (
	"context"

	"github.com/google/uuid"
)

type Image struct {
	ID      uuid.UUID
	Name    string
	Version string
}

func (i *Image) ToString() string {
	return i.Name + ":" + i.Version
}

type ImageService interface {
	BuildImage(ctx context.Context, image *Image) (*Image, error)
	PullImage(ctx context.Context, image *Image) error
	GetImageByID(ctx context.Context, id uuid.UUID) (*Image, error)
}
