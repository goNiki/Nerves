package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/goNiki/Nerves/internal/config"
	"github.com/goNiki/Nerves/internal/infrastructure/database"
	"github.com/goNiki/Nerves/internal/infrastructure/logger"
	"github.com/goNiki/Nerves/internal/infrastructure/migrator"
)

const MigDir = "./migrations"

var configPath = "./internal/config/.env"

func main() {

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	l := logger.InitLogger(cfg.Log)

	l.Log.Info("Конфиг иннициализирован")

	postgres, err := database.InitDatabase(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
		return
	}
	l.Log.Info("Подключение к базе данных выполнено")

	defer postgres.DB.Close()
	migrator := migrator.NewMigrator(postgres.DB, MigDir)
	Migration(l, migrator)

}

func Migration(l *logger.Logger, migrator *migrator.Migrator) {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		var choice int64

		fmt.Println("Выберите действие: \n1 - Создать файл для миграции \n2 - Проверить статус миграций \n3 - Миграций UP \n4- Миграций UpTo \n5 - Миграция Down \n6 - Миграция DownTo \n7 - выйти")
		if scanner.Scan(); scanner.Err() != nil {
			fmt.Println("Ошибка ввода")
			continue
		}
		if _, err := fmt.Sscanf(scanner.Text(), "%d", &choice); err != nil || choice < 1 || choice > 7 {
			fmt.Println("Ошибка ввода. Не соответствует условиям")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("Введите название файла ")
			scanner.Scan()
			name := scanner.Text()
			if err := migrator.Create(name, "sql"); err != nil {
				l.Log.Error("Create file: ", slog.String("Error", err.Error()))
			}
			l.Log.Info("migration created successfully")

		case 2:
			if err := migrator.Status(); err != nil {
				l.Log.Error("Check status: ", slog.String("Error", err.Error()))
			}
			l.Log.Info("Status migrations apply")
		case 3:
			if err := migrator.Up(); err != nil {
				l.Log.Error("Migration UP:", slog.String("Error", err.Error()))
			} else {
				l.Log.Info("UP Migration apply")
			}
		case 4:
			fmt.Println("Введите версию миграции:")
			scanner.Scan()
			var name int64
			if _, err := fmt.Sscanf(scanner.Text(), "%d", &name); err != nil {
				fmt.Println("Ошибка при вводе версии миграции")
				continue
			}
			if err := migrator.UpTo(name); err != nil {
				l.Log.Error("Migration UPto: ", slog.String("Error", err.Error()))
			}
			l.Log.Info("UpTo migration apply")
		case 5:
			if err := migrator.Down(); err != nil {
				l.Log.Error("Migration Down", slog.String("Error", err.Error()))
			}
			l.Log.Info("DOWN migration apply")
		case 6:
			fmt.Println("Введите версию миграции:")
			scanner.Scan()
			var name int64
			if _, err := fmt.Sscanf(scanner.Text(), "%d", &name); err != nil {
				fmt.Println("Ошибка при вводе версии миграции")
				continue
			}

			if err := migrator.DownTo(name); err != nil {
				l.Log.Error("Migration DownTO", slog.String("Error", err.Error()))
			}
			l.Log.Info("DownTo migration apply")
		case 7:
			return
		}
	}
}
