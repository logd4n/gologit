package database

import (
	"errors"
	"logger/models"
)

var (
	GetDataErr  = errors.New("Не удалось получить логи из БД")
	ScanRowsErr = errors.New("Не удалось прочитать результат запроса")
	RowsErr     = errors.New("Не удалось обработать результат запроса")
)

func GetLogs(limit, offset int) ([]models.LogsTable, error) {
	var data []models.LogsTable

	rows, err := dataBase.Query(`
	select * from logs
	order by created_at
	desc limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return nil, GetDataErr
	}

	for rows.Next() {
		var logRow models.LogsTable

		err = rows.Scan(
			&logRow.Id,
			&logRow.Level,
			&logRow.Message,
			&logRow.Created_at,
		)
		if err != nil {
			return nil, ScanRowsErr
		}

		data = append(data, logRow)
	}

	if err = rows.Err(); err != nil {
		return nil, RowsErr
	}

	return data, nil
}
