package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	dataBase   *sql.DB
	driverName = "postgres"
)

func ConnectionAttempt(dsn string) error {
	var db *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		log.Printf("Попытка подключения к PostgreSQL №%d...\t", i)

		db, err = sql.Open(driverName, dsn)

		if err == nil {
			break
		}

		log.Printf("Ошибка подключения: \"%v\"\n", err.Error())
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Printf("Ошибка подключения к БД: \"%v\"", err.Error())
		return err
	}

	dataBase = db
	log.Printf("Подключение к БД выполнено успешно!")
	return nil
}

func ConnectDB() error {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	err := ConnectionAttempt(dsn)
	if err != nil {
		return err
	}

	err = createTable()
	if err != nil {
		return err
	}

	return nil
}

func createTable() error {
	//Таблица logs
	_, err := dataBase.Exec(`
	create table if not exists logs (
	id bigserial primary key,
	level varchar(10) not null,
	message text,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)
	`)

	return err
}
