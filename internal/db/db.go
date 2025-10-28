package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255)	NOT NULL DEFAULT "",
    comment TEXT DEFAULT "",
    repeat VARCHAR(128) DEFAULT ""
 );
 CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init() error {
	var install bool

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}

	DB, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	if install {
		_, err := DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка создания таблицы: %w", err)
		}
	}
	return nil
}
