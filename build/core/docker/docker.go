package docker

import (
	"time"

	"dzor/core"

	"dagger.io/dagger"
)

func StartCompose(ctx core.WrapContext, compose *dagger.Container) {
	compose = compose.
		WithEnvVariable("BURST_CACHE", time.Now().String()).
		WithExec([]string{"docker", "compose", "up", "-d"})

	core.CaptureAndLogStdout(ctx, compose)
}

func StopCompose(ctx core.WrapContext, compose *dagger.Container) {
	compose = compose.
		WithEnvVariable("BURST_CACHE", time.Now().String()).
		WithExec([]string{"docker", "compose", "down"})

	core.CaptureAndLogStdout(ctx, compose)
}

func PrepareCompose(client *dagger.Client) *dagger.Container {
	host := client.Host()

	compose := client.Container(). // platform ??
					From("docker:dind")

	dockerSocket := client.
		Host().
		UnixSocket("/var/run/docker.sock")

	return compose.
		WithFile("/tests/docker-compose.yml", host.Directory(".", core.WithIgnored()).File("docker-compose.yml")).
		WithWorkdir("/tests").
		WithUnixSocket("/var/run/docker.sock", dockerSocket)
}
