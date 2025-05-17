package images

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func BuildImage(ctx context.Context, cli *client.Client, name string, namespace string, dirPath string, relativeDockerPath string, tag string) (types.ImageBuildResponse, error) {
	tarArchive, err := archive.TarWithOptions(dirPath, &archive.TarOptions{})
	if err != nil {
		return types.ImageBuildResponse{}, err
	}
	defer tarArchive.Close()

	imageName := registryURL + "/" + namespace + "/" + name + ":" + tag
	imageBuildResponse, err := cli.ImageBuild(ctx, tarArchive, types.ImageBuildOptions{
		Tags:        []string{imageName},
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
	io.Copy(os.Stdout, imageBuildResponse.Body)
	defer imageBuildResponse.Body.Close()

	return imageBuildResponse, nil
}

// func BuildImage(name string, dirPath string, relativeDockerPath string, tag string) (types.ImageBuildResponse, error) {
// 	ctx := context.Background()
// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	if err != nil {
// 		return types.ImageBuildResponse{}, err
// 	}

// 	dir, err := os.Open(path.Join(dirPath, relativeDockerPath))
// 	if err != nil {
// 		return types.ImageBuildResponse{}, fmt.Errorf("failed to open dockerfile: %v", err)
// 	}
// 	defer dir.Close()

// 	imageBuildResponse, err := cli.ImageBuild(ctx, dir, types.ImageBuildOptions{
// 		Tags:        []string{name + ":" + tag},
// 		Dockerfile:  relativeDockerPath,
// 		Remove:      true,
// 		ForceRemove: true,
// 		NoCache:     true,
// 		BuildArgs: map[string]*string{
// 			"BUILD_NAME": &name,
// 		},
// 	})
// 	if err != nil {
// 		return types.ImageBuildResponse{}, fmt.Errorf("failed to build image: %v", err)
// 	}
// 	defer imageBuildResponse.Body.Close()

// 	return imageBuildResponse, nil
// }
