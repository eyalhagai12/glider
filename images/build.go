package images

import (
	"context"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func BuildImage(name string, dirPath string, relativeDockerPath string, tag string) (types.ImageBuildResponse, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return types.ImageBuildResponse{}, err
	}

	dir, err := os.Open(path.Join(dirPath, relativeDockerPath))
	if err != nil {
		return types.ImageBuildResponse{}, err
	}
	defer dir.Close()

	imageBuildResponse, err := cli.ImageBuild(ctx, dir, types.ImageBuildOptions{
		Tags:        []string{name + ":" + tag},
		Dockerfile:  relativeDockerPath,
		Remove:      true,
		ForceRemove: true,
		NoCache:     true,
		BuildArgs: map[string]*string{
			"BUILD_NAME": &name,
		},
	})
	if err != nil {
		return types.ImageBuildResponse{}, err
	}
	defer imageBuildResponse.Body.Close()

	return imageBuildResponse, nil
}
