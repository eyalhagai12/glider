package backend

import "context"

type Deployer interface {
	Deploy(ctx context.Context, deployment *Deployment, image *Image) (*Deployment, error)
}
