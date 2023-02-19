package core

import (
	"context"
	"fmt"
	"os"

	"dzor/core/config"
	"dzor/core/logger"

	"dagger.io/dagger"
	"go.uber.org/zap"
)

type WrapContext struct {
	Context context.Context
	Client  *dagger.Client
	Log     *zap.SugaredLogger
	Config  config.Config
}

type WrapFunc = func(ctx WrapContext) error

func Wrap(fn WrapFunc) {
	logger := logger.CreateLogger()
	log := logger.Sugar()

	config := config.ReadConfig(log)

	ctx := context.Background()

	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("Error when flushing logger %+v", err)
		}
	}()

	log.Infof("connecting to dagger ...")

	client, clientError := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))

	if clientError != nil || client == nil {
		log.Fatalf("cannot start dagger client! %+v", clientError)
	}

	defer client.
		Close()

	wrapContext := WrapContext{Context: ctx, Client: client, Log: log, Config: config}

	taskError := fn(wrapContext)

	if taskError != nil {
		log.Fatal(taskError)
	}
}

// this will capture stdout only, so if you get error from task it will fail,
// if you need to capture error use the other one
func CaptureAndLogStdout(ctx WrapContext, container *dagger.Container) {
	stdout, err := container.Stdout(ctx.Context)

	if err != nil {
		ctx.Log.Fatal(err)
	}

	ctx.Log.Infof("%v", stdout)
}

// this will capture stderr, so it will NOT stop if your task fails
func CaptureAndLogStderr(ctx WrapContext, container *dagger.Container) {
	exitCode, err := container.ExitCode(ctx.Context)
	if err != nil {
		ctx.Log.Infof("failed with exitCode %v and error %+v", exitCode, err)
	}
}

func WithIgnored() dagger.HostDirectoryOpts {
	return dagger.HostDirectoryOpts{
		Exclude: []string{"**/bin", "**/obj", "**/node_modules", "**/.git", "**/.idea", "**/.vscode", "**/.vs", "**/TestResults"},
	}
}
