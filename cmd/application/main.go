package main

import (
	"context"
	"flag"
	"github.com/GermanBogatov/user-service/internal/application"
	"github.com/GermanBogatov/user-service/internal/config"
	"github.com/GermanBogatov/user-service/pkg/logging"
	"log"
)

func init() {
	systemName := config.SystemName
	serviceEnv := config.ServiceEnv
	logLevel := config.LogLevel

	if systemName == "" {
		systemName = config.Namespace
	}
	if serviceEnv == "" {
		serviceEnv = "dev"
	}
	if logLevel == "" {
		logLevel = "INFO"
	}

	err := logging.InitLogging(&logging.Config{
		SystemName: systemName,
		Env:        serviceEnv,
		Level:      logLevel,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "configPath", "", "path to config file")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// если задан флаг -configPath, то пытаемся считать конфиг из переданного файла, иначе берем из окружения
	var cfg *config.Config
	var err error
	switch cfgPath {
	case "":
		cfg, err = config.NewEnvConfig()
		if err != nil {
			logging.Fatalf("Config error: %s", err)
		}
	default:
		// path to example.env    example: "configs/example.env"
		cfg, err = config.NewEnvConfigFromFile(cfgPath)
		if err != nil {
			cfg, err = config.NewEnvConfig()
			if err != nil {
				logging.Fatalf("Config error: %s", err)
			}
		}
	}

	app, err := application.NewApplication(ctx, cfg)
	if err != nil {
		logging.Fatalf("error create application: %s", err)
	}

	if err := app.Start(ctx); err != nil {
		logging.Fatalf("error start app: %s", err)
	}
}
