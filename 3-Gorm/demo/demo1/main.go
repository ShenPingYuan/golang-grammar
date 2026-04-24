package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
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
		DSN:                       dns,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		// ========== 加上这段：打印 SQL ==========
		Logger: logger.Default.LogMode(logger.Info), // Silent/Error/Warn/Info
		// ======================================
		
	})

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

	// var users = []User{
	// 	{Name: "Jinzhu", Birthday: time.Now()},
	// 	{Name: "Jackson", Birthday: time.Now()},
	// }
	// result := db.Create(&users)
	// if result.Error != nil {
	// 	log.Fatalln(err)
	// }
	// db.Session(&gorm.Session{SkipHooks: true}).Create(&users)
	// fmt.Println("users:", users)

	var user User
	// ctx := context.Background()
	result := db.Where("ID=?", 11).First(&user)
	// user, err := gorm.G[User](db).Where("ID=?", 10).Take(ctx)
	err = result.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println(err)
	} else {
		fmt.Printf("用户信息: %+v\n", user)
	}

	var ctx = context.Background()
	user, _ = gorm.G[User](db).First(ctx)

	// u1, err := gorm.G[User](db).Find(context.Background())
	// fmt.Println("u1:", u1)
	var u2 User
	_ = db.Find(&u2)
	fmt.Println("u2:", u2)

	u3, _ := gorm.G[User](db).Where("id=?", 10).First(ctx)
	fmt.Println("u3:", u3)

	users1, _ := gorm.G[User](db).Where("id IN ?", []int{1, 3, 5}).Find(ctx)
	fmt.Println("users1:", users1)

	var user3 User
	db.Where("id=? or name=?", 11, "Jinzhu").Where("name=?", "Jinzhu").Select("id", "Name", "Created_At").Order("Created_At desc").Distinct("id", "Name", "Created_At").First(&user3)
	fmt.Printf("user3: %+v\n", user3)

	var r []struct {
		Name  string
		Count int64
	}
	db.Model(&User{}).Distinct("name", "age").Group("name").Select("name", "count(*) as count").Find(&r)
	fmt.Println(r)

	var u4 User
	db.Clauses(clause.Locking{
		Strength: "update",
	}).Where("id=?", 10).First(&u4)
	fmt.Printf("user3: %+v\n", u4)

	db.Model(&User{}).Where("id=?", 1).Update("name", "Michael")

	email := "michael@example.com"
	u4.Email = &email
	db.Save(&u4)

	db.Model(&User{}).Where("id=?", 1).Updates(map[string]any{"name": "ms123", "age": 30})
	db.Model(&User{}).Where("id=?", 1).Select("*").Omit("birthday", "created_at", "id").Updates(User{})

	db.Model(&User{}).Where("id=?", 1).Select("age").Updates(map[string]any{"name": "ms123", "age": 30})

	// db.Model(&User{}).Where("id=?", 10).Update("name", gorm.Expr("CONCAT(name, ':', email)"))

	r1, err := gorm.G[Result](db).Raw("select * from users where Id=?", 5).Find(ctx)
	fmt.Println(r1)

	var r2 Result
	db.Raw("select * from users where Id=@id", sql.Named("id", 1)).Find(&r2)
	fmt.Println(r2)
	r3 := db.Exec("update users set name=? where id=?", "周玲", 1)
	fmt.Println(r3.RowsAffected)

	var u5 User
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&u5, 1).Statement
	fmt.Println(stmt.SQL.String())
	fmt.Println(stmt.Vars...)

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&User{}).Where("id = ?", 100).Limit(10).Order("age desc").Find(&[]User{})
	})
	fmt.Println(sql)

	// var uname string
	// var age int
	var r4 Result
	rows, err := db.Table("users").Where("id<3").Select("*").Rows()
	defer rows.Close()
	for rows.Next() {
		// rows.Scan(&r4.ID, &r4.Name)
		db.ScanRows(rows, &r4)
		fmt.Println(r4)
	}
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("\n更新前")
	return nil
}

type Result struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
