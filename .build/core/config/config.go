package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Debug   bool
	Version string
	Images  struct {
		Sdk     string
		Runtime string
	}
	Docker struct {
		Registry string
		Projects []struct {
			GitOps     string
			Name       string
			Path       string
			Entrypoint string
		}
	}
}

func ReadConfig(log *zap.SugaredLogger) Config {
	viper.SetConfigName("build")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".build")

	// pflag.String("task", "ci", "task to run")
	// pflag.Bool("debug", false, "debug mode")

	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("cannot parse flags %e", err)
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("cannot read config %e", err)
	}

	var config Config

	err = viper.Unmarshal(&config)

	if err != nil {
		log.Fatalf("cannot parse config %e", err)
	}

	log.Infof("parsed to %+v\n", config)

	return config
}
