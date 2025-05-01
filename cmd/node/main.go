package main

import (
	"context"

	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()

	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	dockerCli.Ping(ctx)
	defer dockerCli.Close()
}
