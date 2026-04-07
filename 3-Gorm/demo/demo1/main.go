package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           // Standard field for the primary key
	Name         string         // A regular string field
	Email        *string        `gorm:"default:123456@gmail.com"` // A pointer to a string, allowing for null values
	Age          uint8          //`gorm:"default:18"`               // An unsigned 8-bit integer
	Birthday     time.Time      // A pointer to time.Time, can be null
	MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields
	CreatedAt    time.Time      // Automatically managed by GORM for creation time
	UpdatedAt    time.Time      // Automatically managed by GORM for update time
	ignored      string         // fields that aren't exported are ignored
}

func main() {
	dns := "root:root123@tcp(localhost:3306)/gorm_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dns,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 查看当前连接状态
	stats := sqlDB.Stats()
	fmt.Printf("Open: %d, InUse: %d, Idle: %d\n", stats.OpenConnections, stats.InUse, stats.Idle)

	// ========== 加上这行：自动迁移 ==========
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	// ======================================

	// birthday, _ := time.Parse("2006-01-02", "1998-12-04")
	// user := User{Name: "spy", Age: 18, Birthday: birthday}

	// // err = gorm.G[User](db).Create(context.Background(), &user)
	// result := db.Create(&user)
	// if result.Error != nil {
	// 	log.Fatalln(err)
	// }

	//  users := []*User{
	// 	{Name: "Jinzhu", Age: 18, Birthday: time.Now()},
	// 	{Name: "Jackson", Age: 19, Birthday: time.Now()},
	// }

	// fmt.Printf("创建成功，ID: %d\n", user.ID)
	// fmt.Println(result.RowsAffected)

	// result = db.Create(&users)
	// if result.Error != nil {
	// 	log.Fatalln(err)
	// }

	var users = []User{
		{Name: "Jinzhu", Birthday: time.Now()},
		{Name: "Jackson", Birthday: time.Now()},
	}
	result := db.Create(&users)
	if result.Error != nil {
		log.Fatalln(err)
	}
	db.Session(&gorm.Session{SkipHooks: true}).Create(&users)
	fmt.Println("users:", users)

	var user User
	// ctx := context.Background()
	result = db.Where("ID=?", 11).First(&user)
	// user, err := gorm.G[User](db).Where("ID=?", 10).Take(ctx)
	err = result.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println(err)
	} else {
		fmt.Printf("用户信息: %+v\n", user)

	}
}
