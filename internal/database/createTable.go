package database

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
