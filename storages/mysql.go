package storages

import (
	"database/sql"
	"fmt"

	"github.com/g-villarinho/link-fizz-api/config"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBHost,
		config.Env.DBPort,
		config.Env.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
