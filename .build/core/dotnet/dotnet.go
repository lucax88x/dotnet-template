package dotnet

import (
	"io/fs"
	"path"
	"path/filepath"
	"time"

	"dzor/core"

	"dagger.io/dagger"
)

func CreateSdkContainer(ctx core.WrapContext) *dagger.Container {

	sdk := ctx.
		Client.
		Container().
		From(ctx.Config.Images.Sdk)

	if ctx.Config.Debug {
		ctx.Log.Infof("debug mode")

		sdk = sdk.
			WithEnvVariable("BURST_CACHE", time.Now().String())
	}

	return sdk.
		WithEnvVariable("DOTNET_SKIP_FIRST_TIME_EXPERIENCE", "true").
		WithEnvVariable("DOTNET_RUNNING_IN_CONTAINER", "true").
		WithEnvVariable("NUGET_XMLDOC_MODE", "skip")
}

// func cd(client *dagger.Client, isDebug bool) error {
//
//		sdk := CreateSdkContainer(client, isDebug)
//
//		sdk, solutionPath := Restore(client, sdk)
//		sdk = Build(client, sdk, solutionPath)
//		Dockerize(client, sdk)
//
//		return nil
//	}
//
// copies solution and csprojs and then restores, so you have a fixed cache for nuget even if you change the src
func Restore(ctx core.WrapContext, sdk *dagger.Container) (*dagger.Container, string) {
	ctx.Log.Infof("RESTORING")

	host := ctx.
		Client.
		Host()

	nugetCache := ctx.
		Client.CacheVolume("nuget")

	sdk, solutionPath := copySolution(ctx, host, sdk, "/build")
	sdk = copyProjects(ctx, host, sdk, "/build")

	sdk = sdk.
		WithWorkdir("/build")

	sdk = sdk.
		WithMountedCache("/root/.nuget", nugetCache).
		WithExec([]string{"dotnet", "restore", solutionPath})

	captureAndLogStdout(ctx, sdk)

	return sdk, solutionPath
}

// builds
func Build(ctx core.WrapContext, sdk *dagger.Container, solutionPath string) *dagger.Container {
	ctx.Log.Infof("BUILDING")

	host := ctx.
		Client.
		Host()

	sdk = sdk.
		WithDirectory("/build/src", host.Directory("src", withIgnored())).
		WithWorkdir("/build").
		WithExec([]string{"dotnet", "build", "--no-restore", solutionPath})

	captureAndLogStdout(ctx, sdk)

	return sdk
}

// runs unit tests, will not stop if fails, so we can capture testresults
func RunUnitTests(ctx core.WrapContext, sdk *dagger.Container, solutionPath string) *dagger.Container {
	ctx.Log.Infof("RUN UNIT TESTS")

	sdk = sdk.
		WithExec([]string{"dotnet",
			"test",
			"--filter", "Category!=Integration",
			"--logger", "trx",
			"--no-restore",
			"--no-build",
			solutionPath})
	// --logger trx \
	// --logger "console;verbosity=quiet" \
	// --verbosity normal \
	// --no-build --no-restore

	captureAndLogStderr(ctx, sdk)

	return sdk
}

// runs integration tests, will not stop if fails, so we can capture testresults
func RunIntegrationTests(ctx core.WrapContext, sdk *dagger.Container, solutionPath string) *dagger.Container {
	ctx.Log.Infof("RUN INTEGRATION TESTS")

	compose := prepareCompose(ctx.Client)
	startCompose(ctx, compose)

	sdk = sdk.
		WithExec([]string{"dotnet",
			"test",
			"--filter", "Category=Integration",
			"--logger", "trx",
			"--no-restore",
			"--no-build",
			solutionPath})

	captureAndLogStderr(ctx, sdk)

	stopCompose(ctx, compose)

	return sdk
}

func SaveTestResults(ctx core.WrapContext, sdk *dagger.Container) *dagger.Container {
	sdk = sdk.
		WithExec([]string{"find", ".", "-name", "TestResults"})

	stdout, err := sdk.Stdout(ctx.Context)

	ctx.Log.Infof(stdout)

	if err != nil {
		ctx.Log.Fatalf("cannot find test results %v", err)
	}

	// TODO: need to get multiple lines from stdout!

	testResultsDirectory := sdk.
		Directory("./src/Template.Domain.Tests/TestResults")

	_, exportError := testResultsDirectory.
		Export(ctx.Context, path.Join("./TestResults", "./src/Template.Domain.Tests/TestResults"))

	if exportError != nil {
		ctx.Log.Fatalf("cannot export test results %v", exportError)
	}

	captureAndLogStdout(ctx, sdk)

	return sdk
}

func prepareCompose(client *dagger.Client) *dagger.Container {
	host := client.Host()

	compose := client.Container(). // platform ??
					From("docker:dind")

	socket := client.
		Host().
		UnixSocket("/var/run/docker.sock")

	return compose.
		WithFile("/tests/docker-compose.yml", host.Directory(".", withIgnored()).File("docker-compose.yml")).
		WithWorkdir("/tests").
		WithUnixSocket("/var/run/docker.sock", socket)
}

func startCompose(ctx core.WrapContext, compose *dagger.Container) {
	compose = compose.
		WithEnvVariable("BURST_CACHE", time.Now().String()).
		WithExec([]string{"docker", "compose", "up", "-d"})

	captureAndLogStdout(ctx, compose)
}

func stopCompose(ctx core.WrapContext, compose *dagger.Container) {
	compose = compose.
		WithEnvVariable("BURST_CACHE", time.Now().String()).
		WithExec([]string{"docker", "compose", "down"})

	captureAndLogStdout(ctx, compose)
}

// Dockerize
func Dockerize(ctx core.WrapContext, sdk *dagger.Container) {
	// https: //docs.dagger.io/205271/replace-dockerfile
	// https://gist.github.com/gmlewis/536345ad27c6986e41ae8ff7f5c0f7ff

	// TODO: might be worth splitting publish & docker publish so we publish only if all of them published?
	for _, dockerProject := range ctx.Config.Docker.Projects {
		ctx.Log.Infof("Building %s", dockerProject.Name)

		publishPath := path.Join("/publish", dockerProject.Name)
		publishAddress := path.Join(ctx.Config.Docker.Registry, dockerProject.Name)

		sdk = sdk.
			WithExec([]string{
				"dotnet", "publish", dockerProject.Path,
				"-c", "release",
				"-o", publishPath})

		captureAndLogStdout(ctx, sdk)

		publishRef, err := ctx.Client.Container().
			From(ctx.Config.Images.Runtime).
			WithDirectory("/app", sdk.Directory(publishPath)).
			WithEnvVariable("ASPNETCORE_ENVIRONMENT", "PRODUCTION").
			// WithEntrypoint("/app").
			Publish(ctx.Context, publishAddress)

		if err != nil {
			ctx.Log.Fatalf("cannot publish %s with %s", publishAddress, err)
		}

		ctx.Log.Infof("Successfully published %s image to %v - ref: %v", dockerProject.Name, publishAddress, publishRef)
	}

}

func copySolution(ctx core.WrapContext, host *dagger.Host, container *dagger.Container, containerPath string) (*dagger.Container, string) {
	sourceDirectory := host.Directory(".", withIgnored())

	container, solutions := findAndCopyFromHost(ctx, sourceDirectory, container, containerPath, ".", ".sln")
	container = copyFileFromHost(sourceDirectory, container, containerPath, path.Join("src", "global.json"))
	container = copyFileFromHost(sourceDirectory, container, containerPath, path.Join("src", "Directory.Build.props"))

	return container, solutions[0]
}

func copyProjects(ctx core.WrapContext, host *dagger.Host, container *dagger.Container, containerPath string) *dagger.Container {
	sourceDirectory := host.Directory(".", withIgnored())

	container, _ = findAndCopyFromHost(ctx, sourceDirectory, container, containerPath, ".", ".csproj")

	return container
}

func findAndCopyFromHost(ctx core.WrapContext, sourceDirectory *dagger.Directory, container *dagger.Container, containerPath string, sourceRoot string, ext string) (*dagger.Container, []string) {
	files, err := findFilesFromHost(sourceRoot, ext)

	if err != nil || len(files) == 0 {
		ctx.Log.Fatalf("cannot find with %s", ext)
	}

	for _, file := range files {
		container = copyFileFromHost(sourceDirectory, container, containerPath, file)
	}

	return container, files
}

func findFilesFromHost(root string, ext string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(dir.Name()) == ext {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func copyFileFromHost(sourceDirectory *dagger.Directory, container *dagger.Container, root string, path string) *dagger.Container {
	containerPath := filepath.Join(root, path)

	return container.
		WithFile(containerPath, sourceDirectory.File(path))
}

// this will capture stdout only, so if you get error from task it will fail,
// if you need to capture error use the other one
func captureAndLogStdout(ctx core.WrapContext, container *dagger.Container) {
	stdout, err := container.Stdout(ctx.Context)

	if err != nil {
		ctx.Log.Fatal(err)
	}

	ctx.Log.Infof("%s", stdout)
}

// this will capture stderr, so it will NOT stop if your task fails
func captureAndLogStderr(ctx core.WrapContext, container *dagger.Container) {
	exitCode, err := container.ExitCode(ctx.Context)
	if err != nil {
		ctx.Log.Infof("failed with exitCode %s and error %v", exitCode, err)
	}
}

func withIgnored() dagger.HostDirectoryOpts {
	return dagger.HostDirectoryOpts{
		Exclude: []string{"**/bin", "**/obj", "**/node_modules", "**/.git", "**/.idea", "**/.vscode", "**/.vs", "**/TestResults"},
	}
}
