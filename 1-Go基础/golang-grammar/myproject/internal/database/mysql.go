package database

import (
	"database/sql"
	"log/slog"
	"time"
)

type MySQL struct {
	DB *sql.DB
}

// NewMySQL 创建 MySQL 连接
// 需要导入驱动: _ "github.com/go-sql-driver/mysql"
func NewMySQL(dsn string) (*MySQL, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	slog.Info("MySQL connected")
	return &MySQL{DB: db}, nil
}

func (m *MySQL) Close() error {
	return m.DB.Close()
}