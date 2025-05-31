package backend

import "context"

type Deployer interface {
	Deploy(context.Context, *Deployment, *Image, DeploymentMetadata) (*Deployment, error)
}
