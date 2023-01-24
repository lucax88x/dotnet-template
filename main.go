package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"dagger.io/dagger"
)

func main() {
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

	log.Printf("running task %s", task)

	var task_error error
	ctx := context.Background()

	client, client_error := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))

	if client_error != nil || client == nil {
		panic(fmt.Sprintf("cannot start dagger client! %e", client_error))
	}

	defer client.
		Close()

	switch task {
	case "ci":
		{
			task_error = ci(ctx, client, isDebug)
		}
	default:
		{
			log.Printf("unrecognized task")
		}
	}

	if task_error != nil {
		log.Printf("%e", task_error)
	}
}

const dotnetSdkDockerImage = "mcr.microsoft.com/dotnet/sdk:7.0"
const outputDirectory = "build/"

func ci(ctx context.Context, client *dagger.Client, isDebug bool) error {
	now := time.Now()

	log.Printf("CI")

	outputs := client.
		Directory().
		WithoutDirectory("node_modules").
		WithoutDirectory("bin").
		WithoutDirectory("obj")

	sdk := client.
		Container().
		From(dotnetSdkDockerImage)

	if isDebug {
		log.Printf("debug mode")

		sdk = sdk.
			WithEnvVariable("BURST_CACHE", now.String())
	}

	sdk, solutionPath := restore(client, sdk, outputs)
	// sdk = build(client, sdk, outputs, solutionPath)
	runIntegrationTests(client, sdk, solutionPath)

	output := sdk.
		Directory(outputDirectory)

	_, export_error := output.
		Export(ctx, outputDirectory)

	if export_error != nil {
		return export_error
	}

	return nil
}

// copies solution and csprojs and then restores, so you have a fixed cache even if you change the src
func restore(client *dagger.Client, sdk *dagger.Container, outputs *dagger.Directory) (*dagger.Container, string) {

	host := client.Host()

	outputs, solutionPath := copySolution(host, outputs)
	outputs = copyProjects(host, outputs)

	sdk = sdk.
		WithMountedDirectory("/build", outputs).
		WithWorkdir("build")
	//
	// sdk = sdk.
	// 	WithExec([]string{"ls"})
	//
	// sdk = sdk.
	// 	WithExec([]string{"dotnet", "restore", solutionPath})

	return sdk, solutionPath
}

// builds
func build(client *dagger.Client, sdk *dagger.Container, outputs *dagger.Directory, solutionPath string) *dagger.Container {
	host := client.Host()

	// recopy whole host but nuget should be cached
	sdk = sdk.
		WithMountedDirectory("/build/src", host.Directory("src")).
		WithWorkdir("/build")

	sdk = sdk.
		WithExec([]string{"ls"})

	sdk = sdk.
		WithExec([]string{"dotnet", "build", "--no-restore", solutionPath})

	return sdk
}

// runs unit tests
func runUnitTests() {

}

// runs integration tests
func runIntegrationTests(client *dagger.Client, sdk *dagger.Container, solutionPath string) {

	runCompose(client)

	sdk = client.
		Container().
		From(dotnetSdkDockerImage)

		// TODO: check if we need to actively log stdout / stderr?
	sdk.
		WithExec([]string{"dotnet", "test", "--no-restore", "--no-build", solutionPath})
}

// https://discord.com/channels/707636530424053791/1037455051821682718/1066033684278423633
func runCompose(client *dagger.Client) {
	// WithEnvVariable("DOCKER_DEFAULT_PLATFORM", "linux/amd64").
	// WithUnixSocket("/var/run/docker.sock", docker_host).
	// .WithEnvVariable("CACHEBUSTER", datetime.now().strftime("%m/%d/%Y, %H:%M:%S"))
	host := client.Host()

	compose := client.Container(). // platform ??
					From("docker:dind")

	copyFileFromHost(host, compose.Directory("."), ".", "docker-compose.yml")

	compose.WithExec([]string{"dockere", "compose", "up", "-d"})
}

// dockerize
func dockerize() {
	// https: //docs.dagger.io/205271/replace-dockerfile
}

func copySolution(host *dagger.Host, destinationDirectory *dagger.Directory) (*dagger.Directory, string) {
	destinationDirectory, solutions := findAndCopyFromHost(host, destinationDirectory, ".", ".sln")
	destinationDirectory = copyFileFromHost(host, destinationDirectory, ".", path.Join("src", "global.json"))
	destinationDirectory = copyFileFromHost(host, destinationDirectory, ".", path.Join("src", "Directory.Build.props"))

	return destinationDirectory, solutions[0]
}

func copyProjects(host *dagger.Host, destinationDirectory *dagger.Directory) *dagger.Directory {
	destinationDirectory, _ = findAndCopyFromHost(host, destinationDirectory, ".", ".csproj")

	return destinationDirectory
}

func findAndCopyFromHost(host *dagger.Host, destinationDirectory *dagger.Directory, root string, ext string) (*dagger.Directory, []string) {
	files, err := findFiles(root, ext)

	if err != nil || len(files) == 0 {
		panic(fmt.Sprintf("cannot find with %s", ext))
	}

	for _, file := range files {
		destinationDirectory = copyFileFromHost(host, destinationDirectory, root, file)
	}

	return destinationDirectory, files
}

func findFiles(root string, ext string) ([]string, error) {

	// files, err := host.Directory(root, dagger.HostDirectoryOpts{
	// 	Include: []string{ext},
	// }).Entries(ctx)

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

func copyFileFromHost(host *dagger.Host, destinationDirectory *dagger.Directory, root string, path string) *dagger.Directory {
	log.Printf("copying %s", path)

	return destinationDirectory.WithFile(path, host.Directory(root).File(path))
}
