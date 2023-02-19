package dotnet

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"time"

	"dzor/core"
	"dzor/core/docker"

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

	core.CaptureAndLogStdout(ctx, sdk)

	return sdk, solutionPath
}

// builds
func Build(ctx core.WrapContext, sdk *dagger.Container, solutionPath string) *dagger.Container {
	ctx.Log.Infof("BUILDING")

	host := ctx.
		Client.
		Host()

	sdk = sdk.
		WithDirectory("/build/src", host.Directory("src", core.WithIgnored())).
		WithWorkdir("/build").
		WithExec([]string{"dotnet", "build", "--no-restore", solutionPath})

	core.CaptureAndLogStdout(ctx, sdk)

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

	core.CaptureAndLogStderr(ctx, sdk)

	return sdk
}

// runs integration tests, will not stop if fails, so we can capture testresults
func RunIntegrationTests(ctx core.WrapContext, sdk *dagger.Container, solutionPath string) *dagger.Container {
	ctx.Log.Infof("RUN INTEGRATION TESTS")

	compose := docker.PrepareCompose(ctx.Client)
	docker.StartCompose(ctx, compose)

	sdk = sdk.
		WithExec([]string{"dotnet",
			"test",
			"--filter", "Category=Integration",
			"--logger", "trx",
			"--no-restore",
			"--no-build",
			solutionPath})

	core.CaptureAndLogStderr(ctx, sdk)

	docker.StopCompose(ctx, compose)

	return sdk
}

func SaveTestResults(ctx core.WrapContext, sdk *dagger.Container) *dagger.Container {
	sdk = sdk.
		WithExec([]string{"find", ".", "-name", "TestResults"})

	stdout, err := sdk.Stdout(ctx.Context)

	ctx.Log.Infof(stdout)

	if err != nil {
		ctx.Log.Fatalf("cannot find test results %+v", err)
	}

	// TODO: need to get multiple lines from stdout!

	testResultsDirectory := sdk.
		Directory("./src/Template.Domain.Tests/TestResults")

	_, exportError := testResultsDirectory.
		Export(ctx.Context, path.Join("./TestResults", "./src/Template.Domain.Tests/TestResults"))

	if exportError != nil {
		ctx.Log.Fatalf("cannot export test results %+v", exportError)
	}

	core.CaptureAndLogStdout(ctx, sdk)

	return sdk
}

// Dockerize
func Dockerize(ctx core.WrapContext, sdk *dagger.Container, imageTag string) {
	// https: //docs.dagger.io/205271/replace-dockerfile
	// https://gist.github.com/gmlewis/536345ad27c6986e41ae8ff7f5c0f7ff

	// TODO: might be worth splitting publish & docker publish so we publish only if all of them published?
	for _, dockerProject := range ctx.Config.Docker.Projects {
		ctx.Log.Infof("Building %s", dockerProject.Name)

		publishPath := path.Join("/publish", dockerProject.Name)

		var publishAddress string
		if len(ctx.Config.Docker.Registry) == 0 {
			publishAddress = dockerProject.Name
		} else {
			publishAddress = path.Join(ctx.Config.Docker.Registry, dockerProject.Name)
		}

		publishAddress = fmt.Sprintf("%s:%s", publishAddress, imageTag)

		sdk = sdk.
			WithExec([]string{
				"dotnet", "publish", dockerProject.Path,
				"-c", "release",
				"-o", publishPath})

		core.CaptureAndLogStdout(ctx, sdk)

		entrypointDirectory := "/app"
		runtime := ctx.
			Client.
			Container().
			From(ctx.Config.Images.Runtime).
			WithDirectory(entrypointDirectory, sdk.Directory(publishPath)).
			WithEnvVariable("ASPNETCORE_ENVIRONMENT", "PRODUCTION")

		core.CaptureAndLogStdout(ctx, runtime)

		publishRef, err := runtime.
			Publish(ctx.Context, publishAddress)

		if err != nil {
			ctx.Log.Fatalf("cannot publish %s with %+v", publishAddress, err)
		}

		ctx.Log.Infof("Successfully published %s image to %s - ref: %s", dockerProject.Name, publishAddress, publishRef)
	}
}

func copySolution(ctx core.WrapContext, host *dagger.Host, container *dagger.Container, containerPath string) (*dagger.Container, string) {
	sourceDirectory := host.Directory(".", core.WithIgnored())

	container, solutions := findAndCopyFromHost(ctx, sourceDirectory, container, containerPath, ".", ".sln")
	container = copyFileFromHost(sourceDirectory, container, containerPath, path.Join("src", "global.json"))
	container = copyFileFromHost(sourceDirectory, container, containerPath, path.Join("src", "Directory.Build.props"))

	return container, solutions[0]
}

func copyProjects(ctx core.WrapContext, host *dagger.Host, container *dagger.Container, containerPath string) *dagger.Container {
	sourceDirectory := host.Directory(".", core.WithIgnored())

	container, _ = findAndCopyFromHost(ctx, sourceDirectory, container, containerPath, ".", ".csproj")

	return container
}

func findAndCopyFromHost(ctx core.WrapContext, sourceDirectory *dagger.Directory, container *dagger.Container, containerPath string, sourceRoot string, ext string) (*dagger.Container, []string) {
	files, err := findFilesFromHost(sourceRoot, ext)

	if err != nil || len(files) == 0 {
		ctx.Log.Fatalf("cannot find with %s with %+v", ext, err)
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
