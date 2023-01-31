package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"dagger.io/dagger"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLogger() *zap.Logger {
	encoder := zap.NewDevelopmentEncoderConfig()
	encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

	return logger
}

var logger = createLogger()

var log = logger.Sugar()
var ctx = context.Background()

func main() {
	defer logger.Sync()

	var task string
	var isDebug = false

	if len(os.Args) >= 2 {
		task = os.Args[1]
		lastArg := os.Args[len(os.Args)-1]

		if lastArg == "--debug" {
			isDebug = true
		}

	} else {
		task = "ci"
	}

	log.Infof("running task %s", task)

	var task_error error

	client, client_error := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))

	if client_error != nil || client == nil {
		log.Fatalf(fmt.Sprintf("cannot start dagger client! %e", client_error))
	}

	defer client.
		Close()

	switch task {
	case "ci":
		{
			task_error = ci(client, isDebug)
		}
	default:
		{
			log.Infof("unrecognized task")
		}
	}

	if task_error != nil {
		log.Fatalf("%e", task_error)
	}
}

const dotnetSdkDockerImage = "mcr.microsoft.com/dotnet/sdk:7.0"
const outputDirectory = "build/"

func ci(client *dagger.Client, isDebug bool) error {
	now := time.Now()

	log.Infof("CI")

	sdk := client.
		Container().
		From(dotnetSdkDockerImage)

	if isDebug {
		log.Infof("debug mode")

		sdk = sdk.
			WithEnvVariable("BURST_CACHE", now.String())
	}

	sdk = sdk.
		WithEnvVariable("DOTNET_SKIP_FIRST_TIME_EXPERIENCE", "true").
		WithEnvVariable("DOTNET_RUNNING_IN_CONTAINER", "true").
		WithEnvVariable("NUGET_XMLDOC_MODE", "skip")

	sdk, solutionPath := restore(client, sdk)
	sdk = build(client, sdk, solutionPath)
	sdk = runUnitTests(client, sdk, solutionPath)
	sdk = runIntegrationTests(client, sdk, solutionPath)

	captureStdout(sdk)

	output := sdk.
		Directory(outputDirectory)

	_, export_error := output.
		Export(ctx, outputDirectory)

	if export_error != nil {
		return export_error
	}

	return nil
}

// copies solution and csprojs and then restores, so you have a fixed cache for nuget even if you change the src
func restore(client *dagger.Client, sdk *dagger.Container) (*dagger.Container, string) {
	log.Infof("RESTORING")

	host := client.
		Host()

	// solutionPath := "src/Template-Solution.sln"
	sdk, solutionPath := copySolution(host, sdk, "/build")
	sdk = copyProjects(host, sdk, "/build")

	sdk = sdk.
		WithWorkdir("/build")

	sdk = sdk.
		WithExec([]string{"dotnet", "restore", solutionPath})

	captureStdout(sdk)

	return sdk, solutionPath
}

// builds
func build(client *dagger.Client, sdk *dagger.Container, solutionPath string) *dagger.Container {
	log.Infof("BUILDING")

	host := client.
		Host()

	sdk = sdk.
		WithDirectory("/build/src", host.Directory("src", withIgnored())).
		WithWorkdir("/build").
		WithExec([]string{"dotnet", "build", "--no-restore", solutionPath})

	captureStdout(sdk)

	return sdk
}

// runs unit tests
func runUnitTests(client *dagger.Client, sdk *dagger.Container, solutionPath string) *dagger.Container {
	log.Infof("RUN UNIT TESTS")

	sdk = sdk.
		WithExec([]string{"dotnet",
			"test",
			"--filter", "Category!=Integration",
			"--no-restore",
			"--no-build",
			solutionPath})
	// --logger trx \
	// --logger "console;verbosity=quiet" \
	// --verbosity normal \
	// --no-build --no-restore

	captureStdout(sdk)

	return sdk
}

// runs integration tests
func runIntegrationTests(client *dagger.Client, sdk *dagger.Container, solutionPath string) *dagger.Container {
	log.Infof("RUN INTEGRATION TESTS")

	compose := prepareCompose(client)
	startCompose(compose)

	// sdk = client.
	// 	Container().
	// 	From(dotnetSdkDockerImage).
	// 	WithMountedDirectory("/src", client.Host().Directory("./src")).
	// 	WithWorkdir("/src").

	sdk = sdk.
		WithExec([]string{"dotnet",
			"test",
			"--filter", "Category=Integration",
			"--no-restore",
			"--no-build",
			solutionPath})

	captureStdout(sdk)

	stopCompose(compose)

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

func startCompose(compose *dagger.Container) {
	compose = compose.
		WithExec([]string{"docker", "compose", "up", "-d"})

	captureStdout(compose)
}

func stopCompose(compose *dagger.Container) {
	compose = compose.
		WithExec([]string{"docker", "compose", "down"})

	captureStdout(compose)
}

// dockerize
func dockerize() {
	// https: //docs.dagger.io/205271/replace-dockerfile
}

func copySolution(host *dagger.Host, container *dagger.Container, containerPath string) (*dagger.Container, string) {
	sourceDirectory := host.Directory(".", withIgnored())

	container, solutions := findAndCopyFromHost(sourceDirectory, container, containerPath, ".", ".sln")
	container = copyFileFromHost(sourceDirectory, container, containerPath, path.Join("src", "global.json"))
	container = copyFileFromHost(sourceDirectory, container, containerPath, path.Join("src", "Directory.Build.props"))

	return container, solutions[0]
}

func copyProjects(host *dagger.Host, container *dagger.Container, containerPath string) *dagger.Container {
	sourceDirectory := host.Directory(".", withIgnored())

	container, _ = findAndCopyFromHost(sourceDirectory, container, containerPath, ".", ".csproj")

	return container
}

func findAndCopyFromHost(sourceDirectory *dagger.Directory, container *dagger.Container, containerPath string, sourceRoot string, ext string) (*dagger.Container, []string) {
	files, err := findFiles(sourceRoot, ext)

	if err != nil || len(files) == 0 {
		log.Fatalf("cannot find with %s", ext)
	}

	for _, file := range files {
		container = copyFileFromHost(sourceDirectory, container, containerPath, file)
	}

	return container, files
}

func findFiles(root string, ext string) ([]string, error) {
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

	log.Infof("copying %s to %s", path, containerPath)

	return container.
		WithFile(containerPath, sourceDirectory.File(path))
}

func captureStdout(container *dagger.Container) {
	stdout, err := container.Stdout(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%s", stdout)
}

func withIgnored() dagger.HostDirectoryOpts {
	return dagger.HostDirectoryOpts{
		Exclude: []string{"**/bin", "**/obj", "**/node_modules", "**/.git", "**/.idea", "**/.vscode", "**/.vs"},
	}
}
