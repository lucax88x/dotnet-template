package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))

	if err != nil {
		return err
	}

	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	// get `dotnet_sdk` image
	dotnet_sdk := client.Container().From("mcr.microsoft.com/dotnet/sdk:7.0")

	// mount cloned repository into `golang` image
	dotnet_sdk = dotnet_sdk.WithMountedDirectory("/src", src).WithWorkdir("/src")

	// define the application build command
	path := "src/Template-Solution.sln"

	dotnet_sdk = dotnet_sdk.WithExec([]string{"dotnet", "restore", path})

  // how do we only copy csprojs so we cache until there?

	dotnet_sdk = dotnet_sdk.WithExec([]string{"dotnet", "build", "--no-restore", path})

	// get reference to build output directory in container
	output := dotnet_sdk.Directory(path)

	// write contents of container build/ directory to the host
	_, err = output.Export(ctx, path)

	if err != nil {
		return err
	}

	return nil
}
