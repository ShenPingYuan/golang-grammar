package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"myproject/internal/config"
	"myproject/pkg/logger"
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: up or down")
	flag.Parse()

	cfg := config.Load("configs/config.yaml")
	l := logger.New(cfg.Log.Level)

	files, err := filepath.Glob("migrations/*." + *direction + ".sql")
	if err != nil {
		l.Error("failed to read migrations", "err", err)
		os.Exit(1)
	}
	sort.Strings(files)

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			l.Error("failed to read file", "file", f, "err", err)
			os.Exit(1)
		}
		l.Info("would execute migration", "file", f)
		fmt.Println(strings.TrimSpace(string(content)))
		fmt.Println("---")
	}

	// 生产环境应连接真实数据库执行 SQL
	// db := database.NewMySQL(cfg.Database)
	// _, err = db.Exec(string(content))
	l.Info("migration complete", "direction", *direction, "count", len(files))
}