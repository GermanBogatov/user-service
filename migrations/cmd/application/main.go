package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/GermanBogatov/user-service/migrations/configs"
	psqlMigration "github.com/GermanBogatov/user-service/migrations/psqlmigrations"
	"github.com/GermanBogatov/user-service/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"log"
	"os"
)

func init() {
	systemName := os.Getenv("USER_SERVICE_MIGRATIONS_SYSTEM_NAME")
	if systemName == "" {
		systemName = configs.Namespace
	}

	serviceEnv := os.Getenv("USER_SERVICE_MIGRATIONS_ENV")
	if serviceEnv == "" {
		serviceEnv = "dev"
	}
	logLevel := os.Getenv("USER_SERVICE_MIGRATIONS_LOG_LEVEL")
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		migrate, cfgPath string
		version          int64
	)

	logging.Info("parsing flags...")
	flag.StringVar(&cfgPath, "configPath", "", "path to config file")
	flag.StringVar(&migrate, "migrate", "", "migration up/down/up-to/down-to/up-one")
	flag.Int64Var(&version, "version", 0, "version migration")
	flag.Parse()

	if migrate == "" {
		logging.Fatalf("flag `migrate` is empty")
	}

	// если задан флаг -configPath, то пытаемся считать конфиг из переданного файла, иначе берем из окружения
	var cfg *configs.Config
	var err error
	switch cfgPath {
	case "":
		cfg, err = configs.NewEnvConfig()
		if err != nil {
			logging.Fatalf("Config error: %s", err)
		}
	default:
		// path to example.env    example: "configs/example.env"
		cfg, err = configs.NewEnvConfigFromFile(cfgPath)
		if err != nil {
			cfg, err = configs.NewEnvConfig()
			if err != nil {
				logging.Fatalf("Config error: %s", err)
			}
		}
	}

	// проверяем url postgres
	_, err = pgxpool.ParseConfig(cfg.Postgres.URL)
	if err != nil {
		logging.Fatalf("error parse config: %s", err)
	}

	logging.Info("opening db...")
	stdlib.GetDefaultDriver()
	db, err := goose.OpenDBWithDriver("pgx", cfg.Postgres.URL)
	if err != nil {
		logging.Fatalf("error open db with driver: %s", err)
	}

	defer func() {
		errClose := db.Close()
		if errClose != nil {
			logging.Errorf("error close db: %s", err)
		}
	}()

	goose.SetBaseFS(&psqlMigration.MigratePostgres)

	err = goose.SetDialect("postgres")
	if err != nil {
		logging.Fatalf("error set dialect: %s", err)
	}

	logging.Info("start migration...")
	err = migrationSelection(ctx, migrate, version, db)
	if err != nil {
		logging.Fatalf("migration selection: %s", err)
	}

}

func migrationSelection(ctx context.Context, migrate string, version int64, db *sql.DB) error {
	var err error

	switch migrate {
	case "up":
		err = migrateUp(ctx, db)
		if err != nil {
			return err
		}
	case "up-one":
		err = migrateUpByOne(ctx, db)
		if err != nil {
			return err
		}

	case "up-to":
		if version != 0 {
			err = migrateUpTo(ctx, db, version)
			if err != nil {
				return err
			}
		} else {
			return errors.New("required flag version")
		}

	case "down":
		err = migrateDown(ctx, db)
		if err != nil {
			return err
		}

	case "down-to":
		if version != 0 {
			err = migrateDownTo(ctx, db, version)
			if err != nil {
				return err
			}
		} else {
			return errors.New("required flag version")
		}

	default:
		return errors.New("invalid flag: migrate")
	}

	return nil
}

// migrateUp - накатывает все доступные миграции
func migrateUp(ctx context.Context, db *sql.DB) error {
	err := goose.UpContext(ctx, db, ".")
	if err != nil {
		return errors.Wrap(err, "migrate up")
	}
	return nil
}

// migrateUp - накатывает все доступные миграции
func migrateUpTo(ctx context.Context, db *sql.DB, version int64) error {
	err := goose.UpToContext(ctx, db, ".", version)
	if err != nil {
		return errors.Wrap(err, "migrate up to")
	}
	return nil
}

// migrateDown - откатывает миграцию на прошлую (-1 итерация)
func migrateDown(ctx context.Context, db *sql.DB) error {
	err := goose.DownContext(ctx, db, ".")
	if err != nil {
		return errors.Wrap(err, "migrate down")
	}
	return nil
}

// migrateDownTo - откатывает миграцию к определенной версии
func migrateDownTo(ctx context.Context, db *sql.DB, version int64) error {
	err := goose.DownToContext(ctx, db, ".", version)
	if err != nil {
		return errors.Wrap(err, "migrate down to")
	}
	return nil
}

// migrateUpByOne - накатывает одну миграцию (+1 итерация)
func migrateUpByOne(ctx context.Context, db *sql.DB) error {
	err := goose.UpByOneContext(ctx, db, ".")
	if err != nil {
		return errors.Wrap(err, "migrate up by one")
	}

	return nil
}
