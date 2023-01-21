package main

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

	"os"

	"dagger.io/dagger"
	"github.com/fatih/color"
)

func main() {
	var task string

	if len(os.Args) >= 2 {
		task = os.Args[1]
	} else {
		task = "ci"
	}

	log_info("running task %s", task)

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
			task_error = build(ctx, client)
		}
	default:
		{
			log_error("unrecognized task")
		}
	}

	if task_error != nil {
		log_error("%e", task_error)
	}
}

const sdk_docker_image = "mcr.microsoft.com/dotnet/sdk:7.0"
const solution = "src/Template-Solution.sln"

func build(ctx context.Context, client *dagger.Client) error {
	log_info("building")

	// host_src := client.
	// 	Host().
	// 	Directory("./src")

	outputs := client.
		Directory()

	sdk := client.
		Container().
		From(sdk_docker_image)

	projects_dir, err := copyProjects(outputs)

	if err != nil {
		return err
	}

	sdk = sdk.
		WithMountedDirectory("/src", projects_dir).
		WithWorkdir("/src")

	sdk = sdk.
		WithExec([]string{"ls"})

	sdk = sdk.
		WithExec([]string{"ls", "src"})

	sdk = sdk.
		WithExec([]string{"dotnet", "restore", solution})

	// sdk = sdk.
	// 	WithExec([]string{"dotnet", "build", "--no-restore", solution})

	// get reference to build output directory in container

	output := sdk.
		Directory(solution)

	// output := client.
	// 	Directory()

	// write contents of container build/ directory to the host
	_, export_error := output.
		Export(ctx, solution)

	if export_error != nil {
		return export_error
	}

	return nil
}

func copyProjects(destinationDirectory *dagger.Directory) (*dagger.Directory, error) {
	// copy solution
	// projects_dir.WithFile(solution, )

	files, err := find_files(".", ".csproj")

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)

		if err != nil {
			panic(fmt.Sprintf("cannot read %s", file))
		}

		log_info("copying %s project", file)

		destinationDirectory = destinationDirectory.WithNewFile(file, string(content))
	}

	if err != nil {
		return nil, err
	}

	return destinationDirectory, nil
}

func find_files(root string, ext string) ([]string, error) {
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

func log_error(message string, args ...any) {
	log(color.New(color.FgBlue), message, args)
}

func log_info(message string, args ...any) {
	log(color.New(color.FgBlue), message, args)
}

func log(logger *color.Color, message string, args ...any) {
	// fmt.Println(message)
	// fmt.Println(len(args))
	// fmt.Println(args[0])

	if args != nil && len(args) > 0 {
		logger.Printf(message, args)
	} else {
		logger.Print(message)
	}

	logger.Println()
}
