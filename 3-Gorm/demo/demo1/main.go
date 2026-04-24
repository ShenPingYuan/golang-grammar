package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID           uint                          `gorm:"autoIncrement:true"`                       // Standard field for the primary key
	Name         string                        `gorm:"index;check:name<>'ms'"`                   // A regular string field
	Email        *string                       `gorm:"default:123456@gmail.com;index:idx_email"` // A pointer to a string, allowing for null values
	Age          uint8                         //`gorm:"default:18"`               // An unsigned 8-bit integer
	Birthday     time.Time                     // A pointer to time.Time, can be null
	MemberNumber sql.NullString                // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime                  // Uses sql.NullTime for nullable time fields
	CreatedAt    time.Time                     // Automatically managed by GORM for creation time
	UpdatedAt    time.Time                     // Automatically managed by GORM for update time
	ignored      string                        // fields that aren't exported are ignored
	Addresses    datatypes.JSONType[[]Address] `gorm:"type:json"`
}

type Address struct {
	City string `json:"city"`
	Line string `json:"line"`
}

func main() {
	dns := "root:8rME16k*8a0iLMIP@tcp(192.168.1.63:13306)/gorm_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, _ := gorm.Open(mysql.New(mysql.Config{
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
	db.AutoMigrate(&User{})

	// u := User{
	// 	Name: "spy1",
	// 	Age:  18,
	// 	Addresses: datatypes.NewJSONType([]Address{
	// 		{City: "Shanghai", Line: "Road 1"},
	// 	}),
	// }
	// db.Create(&u)

	// // 更新（整体替换）
	// u.Addresses = datatypes.NewJSONType([]Address{
	// 	{City: "Beijing", Line: "Street 2"},
	// })
	// db.Save(&u)

	var u2 User
	db.Model(&User{}).Where("id=?", 70).First(&u2)
	fmt.Println(u2)
	u2.Addresses.Data()

}

type Post struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title string
	Tags  pq.StringArray `gorm:"type:text[]"`
}

// 用于筛选 id 大于 10 的记录范围
func IdGreaterThan10(db *gorm.DB) *gorm.DB {
	return db.Where("id > ?", 10)
}

type Json json.RawMessage

func (j *Json) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = Json(result)
	return err
}

func (j Json) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
