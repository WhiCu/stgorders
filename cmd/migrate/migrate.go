package migrate

import (
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDB(dsn string, migrateDir string, log *slog.Logger) error {

	m, err := migrate.New(
		"file://"+migrateDir,
		dsn,
	)
	if err != nil {
		slog.Error("Невозможно создать экземпляр миграций", slog.String("ERR", err.Error()))
		return err
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Error("Невозможно применить миграции", slog.String("ERR", err.Error()))
			return err
		}
		log.Info("Миграции уже были применены")
		return nil
	}

	log.Info("Миграции успешно применены")
	return nil
}
