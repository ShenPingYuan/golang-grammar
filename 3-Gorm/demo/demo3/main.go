package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Name      string
	Age       int
	Email     string
	CompanyId uint
	Company   Company    `gorm:"foreignKey:CompanyId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Profile   Profile    `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Languages []Language `gorm:"many2many:user_languages;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Age <= 0 {
		err = errors.New("age cannot less than 0")
	}
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		tx.Model(u).Update("role", "admin")
	}
	return
}

type Language struct {
	Id   uint
	Code string
	Name string
}

type Profile struct {
	Id     uint `gorm:"primaryKey"`
	UserID uint `gorm:"unique"`
	Bio    string
	Img    string
}

type Company struct {
	Id    uint
	Name  string
	Users []User `gorm:"foreignKey:CompanyId"`
}

func main() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	// 连接数据库、初始化 GORM 等操作...
	dns := "root:root123@tcp(localhost:3306)/gorm_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dns,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})

	var user User = User{}
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
	fmt.Println(stmt.SQL.String())
	fmt.Println(stmt.Vars)
	finalSql := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	fmt.Println(finalSql)

	db.Session(&gorm.Session{})

	ctx := context.Background()
	db.Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[User](tx).First(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	// tx := db.Begin()
	// tx.Rollback()
	// tx.Commit()
}
