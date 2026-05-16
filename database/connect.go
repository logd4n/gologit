package database

import (
	"database/sql"
	"errors"
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

var (
	DB_TryConnErr = errors.New("Ошибка подключения")
	DB_ConnErr    = errors.New("Ошибка подключения к БД")
	CrTableErr    = errors.New("Не удалось создать таблицу")
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

		log.Printf("%v\n", DB_TryConnErr.Error())
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Printf("%v: \"%v\"\n", DB_ConnErr.Error(), err.Error())
		return DB_ConnErr
	}

	dataBase = db
	log.Printf("Подключение к БД выполнено успешно!\n")
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
	if err != nil {
		return CrTableErr
	}

	return nil
}
