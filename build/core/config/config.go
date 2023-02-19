package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Debug  bool
	Images struct {
		Sdk     string
		Runtime string
	}
	Docker struct {
		Registry string
		Projects []struct {
			Name       string
			Path       string
			Entrypoint string
		}
	}
	GitOps struct {
		UseLocalSsh bool
		SshHost     string
		Url         string
		Email       string
		Name        string
	}
}

func ReadConfig(log *zap.SugaredLogger) Config {
	viper.SetConfigName("build")
	viper.AddConfigPath(".")
	viper.AddConfigPath("build")

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("cannot parse flags %+v", err)
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("cannot read config %+v", err)
	}

	var config Config

	err = viper.Unmarshal(&config)

	if err != nil {
		log.Fatalf("cannot parse config %+v", err)
	}

	log.Infof("parsed to %+v", config)

	return config
}
