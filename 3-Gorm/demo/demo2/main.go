package main

import (
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
		Logger: logger.Default.LogMode(logger.Info), // Silent/Error/Warn/Info
	})

	// if err != nil {
	// 	log.Fatalf("failed to connect database: %v", err)
	// }

	// err = db.AutoMigrate(&User{}, &Company{}, &Profile{}, &Language{})
	// if err != nil {
	// 	log.Fatalf("failed to auto migrate tables: %v", err)
	// }

	// user := User{
	// 	Name:  "spy",
	// 	Email: "24397399@qq.com",
	// 	Age:   18,
	// 	Profile: Profile{
	// 		Img: "http://example.com/spy.jpg",
	// 		Bio: "I am spy.",
	// 	},
	// 	Company: Company{
	// 		Name: "内江供电公司",
	// 	},
	// 	Languages: []Language{
	// 		{Code: "en", Name: "English"}, // 多对多关系
	// 		{Code: "zh", Name: "中文"},      // 多对多关系
	// 	},
	// }

	// fmt.Println("开始创建用户")

	// db.Create(&user)

	// 创建记录
	// company := &Company{Name: "达州供电公司"}
	// result := db.Create(&company)
	// if result.Error != nil {
	// 	fmt.Println(result.Error)
	// }

	// user := User{Name: "李四", Email: "2439739932@qq.com", Age: 18, CompanyId: 1}
	// result := db.Create(&user)
	// if result.Error != nil {
	// 	fmt.Println(result.Error)
	// }

	var user User
	db.Model(&User{}).Where("Id=?", 4).Preload("Profile").Preload("Languages").Preload("Company").First(&user)

	languages := user.Languages
	user.Languages = languages[:len(languages)-1] // 删除最后一个语言

	// user.Company.Name = "成都供电公司" // 修改公司名称
	user.Age = 16
	// db.Save(&user) // 更新用户信息，包括关联的公司信息

	db.Model(&user).Association("Languages").Replace(languages[:len(languages)-1]) // 更新关联的语言
	// var company Company
	// db.Where("id=?", 4).Preload("Users").Preload("Users.Profile").First(&company)
	// fmt.Println(user)
	// fmt.Println(company)
}
