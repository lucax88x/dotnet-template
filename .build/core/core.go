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
			fmt.Printf("Error when flushing logger %w", err)
		}
	}()

	client, clientError := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))

	if clientError != nil || client == nil {
		log.Fatalf(fmt.Sprintf("cannot start dagger client! %w", clientError))
	}

	defer client.
		Close()

	wrapContext := WrapContext{Context: ctx, Client: client, Log: log, Config: config}

	taskError := fn(wrapContext)

	if taskError != nil {
		log.Fatal(taskError)
	}
}
