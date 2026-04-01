package database

import "log/slog"

type Mongo struct {
	URI      string
	Database string
}

// NewMongo 创建 MongoDB 连接
// 生产环境请使用 go.mongodb.org/mongo-driver
func NewMongo(uri, database string) *Mongo {
	slog.Info("MongoDB configured", "uri", uri, "database", database)
	return &Mongo{URI: uri, Database: database}
}