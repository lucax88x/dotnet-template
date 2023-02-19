package git

import (
	"fmt"
	"os"
	"path"
	"time"

	"dzor/core"

	"dagger.io/dagger"
	"gopkg.in/yaml.v2"
)

// https://gist.github.com/gmlewis/680621bc9ed2477e6cfa5832fcb7194e
// pushes new version on gitops, this assumes a SINGLE gitops repository for now
// if ssh doesnt work, consider using uselocalssh
func PatchGitOps(ctx core.WrapContext, sdk *dagger.Container, imageTag string, commitMessage string) {

	imageTagData, err := serializeToYaml(ctx, imageTag)

	if err != nil {
		ctx.Log.Fatalf("Error while Marshaling. %+v", err)
	}

	gitContainer := ctx.
		Client.
		Container().
		From("alpine/git").
		WithEntrypoint([]string{}).
		WithEnvVariable("BURST_CACHE", time.Now().String())

	if ctx.Config.GitOps.UseLocalSsh {

		homeDir, err := os.UserHomeDir()

		if err != nil {
			ctx.Log.Fatal(err)
		}

		sshDir := ctx.
			Client.
			Host().
			Directory(path.Join(homeDir, ".ssh"))

		gitContainer = gitContainer.
			WithMountedDirectory("/root/.ssh", sshDir)

	} else {
		sshAgentPath := os.Getenv("SSH_AUTH_SOCK")

		sshSocket := ctx.
			Client.
			Host().
			UnixSocket(sshAgentPath)

		gitContainer = gitContainer.
			WithUnixSocket("/default.ssh", sshSocket).
			WithEnvVariable("SSH_AUTH_SOCK", "/default.ssh").
			WithExec([]string{"mkdir", "-p", "/root/.ssh"}).
			WithExec([]string{"ash", "-c", fmt.Sprintf("ssh-keyscan -t rsa %s >> /root/.ssh/known_hosts", ctx.Config.GitOps.SshHost)})
	}

	gitContainer = gitContainer.
		WithExec([]string{"git", "clone", ctx.Config.GitOps.Url, "/git"}).
		WithWorkdir("/git").
		WithExec([]string{"git", "config", "user.name", ctx.Config.GitOps.Name}).
		WithExec([]string{"git", "config", "user.email", ctx.Config.GitOps.Email}).
		WithNewFile("imageTag.yml", dagger.ContainerWithNewFileOpts{Contents: imageTagData}).
		WithExec([]string{"git", "status"}).
		WithExec([]string{"git", "add", "-A"}).
		WithExec([]string{"git", "status"}).
		WithExec([]string{"git", "commit", "-m", fmt.Sprintf("'%s'", commitMessage)}).
		WithExec([]string{"git", "pull", "--rebase"}).
		WithExec([]string{"git", "push", "origin"}).
		WithExec([]string{"git", "tag", imageTag}).
		WithExec([]string{"git", "push", "origin", imageTag})

	core.CaptureAndLogStderr(ctx, gitContainer)
}

func serializeToYaml(ctx core.WrapContext, imageTag string) (string, error) {
	s1 := imageTagYml{Image: struct{ Tag string }{Tag: imageTag}}

	yamlData, err := yaml.Marshal(&s1)

	if err != nil {
		return "", err
	}

	return string(yamlData), nil
}

type imageTagYml struct {
	Image struct {
		Tag string
	}
}
