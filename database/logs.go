package database

import (
	"errors"
	"logger/models"
)

func SendLogs(logs *models.LogsTable) error {
	_, err := dataBase.Exec(`
	insert into logs (
	level,
	message,
	created_at
	)
	values (
	$1,
	$2,
	$3
	)
	`, logs.Level, logs.Message, logs.Created_at)

	if err != nil {
		return errors.New("Не удалось записать данные: " + err.Error())
	}

	return nil
}
