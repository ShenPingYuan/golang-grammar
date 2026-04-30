# GORM 框架完全指南

## 前言

GORM 是 Go 语言最流行的 ORM 框架。功能包括 CRUD、关联关系、事务、钩子、预加载、自动迁移、泛型 API 等。因为官方文章以及说得很详细的了，因此本文主要以代码方式梳理记录 GORM 的核心用法，可作为日后开发速查手册。

## 目录

1. 安装与连接
2. 模型定义
3. 自动迁移
4. 创建（Create）
5. 查询（Query）
6. 高级查询
7. 更新（Update）
8. 删除（Delete）
9. 原生SQL
10. 关联关系
11. 预加载
12. 事务
13. Hook（钩子）
14. Scopes 与链式调用
15. 泛型 API（v1.30+）
16. Session 与上下文
17. 自定义数据类型
18. 性能优化
19. 常用配置
20. 常见易错点

---

## 1. 安装与连接

```go
// 安装（在项目根目录执行）
// go get -u gorm.io/gorm            ← GORM 核心库
// go get -u gorm.io/driver/mysql    ← MySQL 驱动
// go get -u gorm.io/driver/postgres ← PostgreSQL 驱动
// go get -u gorm.io/driver/sqlite   ← SQLite 驱动（纯 Go，无需 CGO）

package main

import (
    "time"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    // ===================== MySQL 连接 =====================
    // DSN 格式: 用户名:密码@tcp(主机:端口)/数据库名?参数
    // charset=utf8mb4    → 支持 emoji 等 4 字节 Unicode
    // parseTime=True     → 将 MySQL 的时间类型解析为 Go 的 time.Time（必须开启）
    // loc=Local          → 使用本地时区
    dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    // ===================== PostgreSQL 连接 =====================
    // sslmode=disable → 本地开发关闭 SSL，生产环境应改为 require
    dsn2 := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"
    db, err = gorm.Open(postgres.Open(dsn2), &gorm.Config{})

    // ===================== SQLite 连接 =====================
    // 文件路径，如果文件不存在会自动创建
    db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

    if err != nil {
        panic("failed to connect database")
    }

    // ===================== 连接池配置 =====================
    // GORM 底层使用 database/sql 的连接池，通过 db.DB() 获取 *sql.DB 进行配置
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数（连接用完后保留多少个不关闭，等待复用）
    sqlDB.SetMaxOpenConns(100)          // 最大打开连接数（同时能有多少个连接在用）
    sqlDB.SetConnMaxLifetime(time.Hour) // 单个连接最大存活时间（超过后被关闭重建，防止 DB 端超时断开）
}
```

---

## 2. 模型定义

```go
import "gorm.io/gorm"

// ===================== gorm.Model 内置字段 =====================
// GORM 提供了一个内置的 gorm.Model 结构体，包含以下 4 个字段：
// type Model struct {
//     ID        uint           `gorm:"primarykey"`  // 主键，自增
//     CreatedAt time.Time      // 创建时间，Create 时自动填充
//     UpdatedAt time.Time      // 更新时间，每次 Save/Update 自动更新
//     DeletedAt gorm.DeletedAt `gorm:"index"`       // 软删除时间，调用 Delete 时填充，不会真正删除记录
// }
// 嵌入 gorm.Model 就自动拥有这 4 个字段，不需要自己写

// ===================== 基本模型示例 =====================
type User struct {
    gorm.Model                          // 内嵌，自动获得 ID, CreatedAt, UpdatedAt, DeletedAt
    Name     string `gorm:"type:varchar(100);not null;index"`
    //                type:varchar(100) → 数据库列类型为 VARCHAR(100)
    //                not null          → 不允许 NULL
    //                index             → 为此字段创建普通索引

    Email    string `gorm:"uniqueIndex;size:128"`
    //                uniqueIndex → 创建唯一索引（不允许重复值）
    //                size:128    → 等同于 type:varchar(128)

    Age      int    `gorm:"default:18"`
    //                default:18 → 数据库默认值为 18（插入时如果没赋值就用 18）

    Active   bool   `gorm:"default:true"`

    Profile  Profile  // Has One 关联：一个 User 拥有一个 Profile
    Orders   []Order  // Has Many 关联：一个 User 拥有多个 Order
}

// ===================== 自定义表名 =====================
// 默认 GORM 会把结构体名转为 snake_case 复数：User → users
// 如果想自定义，实现 TableName() 方法
func (User) TableName() string {
    return "t_users" // 表名变为 t_users
}

// ===================== 常用 struct tag 详解 =====================
type Product struct {
    ID        uint    `gorm:"primaryKey;autoIncrement"`
    //                  primaryKey     → 主键
    //                  autoIncrement  → 自增（MySQL 默认主键就自增，这里是显式声明）

    Code      string  `gorm:"type:varchar(50);uniqueIndex;not null;comment:产品编码"`
    //                  comment:产品编码 → 数据库列注释（DDL 里的 COMMENT '产品编码'）

    Name      string  `gorm:"size:200;not null"`
    //                  size:200 → varchar(200)

    Price     float64 `gorm:"type:decimal(10,2);default:0"`
    //                  decimal(10,2) → 总共 10 位数字，小数点后 2 位

    Stock     int     `gorm:"check:stock >= 0"`
    //                  check:stock >= 0 → 数据库 CHECK 约束，stock 不能小于 0

    Category  string  `gorm:"index:idx_cat_price"`
    SalePrice float64 `gorm:"index:idx_cat_price"`
    //                  两个字段用相同的索引名 idx_cat_price → 组成复合索引
    //                  等价于 CREATE INDEX idx_cat_price ON products(category, sale_price)

    Remark    string  `gorm:"type:text"`
    //                  text → MySQL TEXT 类型（长文本）

    IgnoreMe  string  `gorm:"-"`
    //                  "-" → GORM 完全忽略此字段（不建表、不读写）

    ReadOnly  string  `gorm:"->"`
    //                  "->" → 只读：只从 DB 读取，写入时忽略

    WriteOnly string  `gorm:"<-"`
    //                  "<-" → 只写：创建和更新都可写（但不影响读取，读取也能读到）

    CreateOnly string `gorm:"<-:create"`
    //                  "<-:create" → 仅在 Create 时可写，Update 时忽略

    UpdateOnly string `gorm:"<-:update"`
    //                  "<-:update" → 仅在 Update 时可写，Create 时忽略
}

// ===================== JSON 标签配合使用 =====================
// gorm tag 控制数据库行为，json tag 控制 JSON 序列化行为
type APIUser struct {
    gorm.Model
    Name  string `gorm:"size:100" json:"name"`              // DB 列 varchar(100)，JSON key 为 "name"
    Email string `gorm:"size:128" json:"email,omitempty"`    // omitempty: JSON 序列化时空值不输出
}

// ===================== 复合主键 =====================
// 多个字段共同组成主键
type UserLanguage struct {
    UserID     uint   `gorm:"primaryKey"` // 联合主键的一部分
    LanguageID uint   `gorm:"primaryKey"` // 联合主键的另一部分
    Skill      string
    // 建表：PRIMARY KEY (user_id, language_id)
}
```

---

## 3. 自动迁移

```go
// AutoMigrate 会根据结构体定义自动：
//   ✅ 创建不存在的表
//   ✅ 添加缺失的列
//   ✅ 添加缺失的索引
//   ❌ 不会删除多余的列（防止误删数据）
//   ❌ 不会修改已有列的类型（防止数据丢失）
//   ❌ 不会删除已有索引
// 因此 AutoMigrate 适合开发阶段，生产环境推荐用 golang-migrate / goose 等专用迁移工具

db.AutoMigrate(&User{}, &Product{}, &Order{})
// 依次检查 users、products、orders 三张表，不存在就建，缺字段就加

// ===================== 判断表是否存在 =====================
db.Migrator().HasTable(&User{})   // 通过结构体判断（会调用 TableName()）
db.Migrator().HasTable("users")   // 通过表名字符串判断

// ===================== 手动建表/删表 =====================
db.Migrator().CreateTable(&User{}) // CREATE TABLE users (...)
db.Migrator().DropTable(&User{})   // DROP TABLE users
db.Migrator().DropTable("users")   // 也可以传表名

// ===================== 列操作 =====================
db.Migrator().AddColumn(&User{}, "Age")    // ALTER TABLE users ADD COLUMN age ...
db.Migrator().DropColumn(&User{}, "Age")   // ALTER TABLE users DROP COLUMN age
db.Migrator().AlterColumn(&User{}, "Name") // ALTER TABLE users MODIFY COLUMN name ...（按最新 struct tag）
db.Migrator().HasColumn(&User{}, "Name")   // 判断列是否存在 → true/false

// ===================== 索引操作 =====================
db.Migrator().CreateIndex(&User{}, "Name")         // 为 Name 字段创建索引
db.Migrator().DropIndex(&User{}, "idx_name")        // 删除索引
db.Migrator().HasIndex(&User{}, "idx_name")          // 判断索引是否存在

// ===================== 重命名 =====================
db.Migrator().RenameTable("users", "t_users")          // 重命名表
db.Migrator().RenameColumn(&User{}, "Name", "FullName") // 重命名列
```

---

## 4. 创建（Create）

```go
// ===================== 创建单条记录 =====================
user := User{Name: "Alice", Age: 25, Email: "alice@go.dev"}
result := db.Create(&user)
// 生成 SQL: INSERT INTO users (name, age, email, created_at, updated_at) VALUES ('Alice', 25, 'alice@go.dev', '2024-...', '2024-...')
// 执行后：
//   user.ID             → 自动回填数据库生成的主键（比如 1）
//   user.CreatedAt      → 自动回填创建时间
//   result.Error        → 如果插入失败，这里是 error；成功则为 nil
//   result.RowsAffected → 插入了几行（通常为 1）

// ===================== 指定字段创建 =====================
// 只插入 Name 和 Age，其他字段即使有值也不会插入（使用数据库默认值或零值）
db.Select("Name", "Age").Create(&user)
// 生成 SQL: INSERT INTO users (name, age, created_at, updated_at) VALUES ('Alice', 25, ...)
// Email 字段被忽略，即使 user.Email 有值

// ===================== 忽略字段创建 =====================
// 除 Age 外的所有字段都插入
db.Omit("Age").Create(&user)
// 生成 SQL: INSERT INTO users (name, email, created_at, updated_at) VALUES ('Alice', 'alice@go.dev', ...)
// Age 被忽略，使用数据库默认值（如果 gorm tag 里写了 default:18 就是 18）

// ===================== 批量创建 =====================
users := []User{
    {Name: "Bob", Age: 20},
    {Name: "Carol", Age: 22},
    {Name: "Dave", Age: 30},
}
db.Create(&users)
// 生成 SQL: INSERT INTO users (name, age, ...) VALUES ('Bob',20,...), ('Carol',22,...), ('Dave',30,...)
// 一条 SQL 插入多行
// 执行后 users[0].ID, users[1].ID, users[2].ID 都会回填

// ===================== 分批插入 =====================
// 如果 users 有 10000 条，一次性 INSERT 可能超过 MySQL 的 max_allowed_packet
// CreateInBatches 会自动分成多条 SQL，每条最多插入 100 行
db.CreateInBatches(&users, 100)
// 生成: 第 1 条 INSERT 100 行，第 2 条 INSERT 100 行 ... 直到全部插完

// ===================== 用 Map 创建 =====================
// 不基于结构体，直接用 map[string]interface{} 指定列名和值
// ⚠️ 不会触发 Hook（BeforeCreate 等）
// ⚠️ 不会回填主键
// ⚠️ 不会自动处理关联
db.Model(&User{}).Create(map[string]interface{}{
    "Name": "Eve", "Age": 28,
})
// 生成 SQL: INSERT INTO users (name, age) VALUES ('Eve', 28)

// 批量 Map 创建
db.Model(&User{}).Create([]map[string]interface{}{
    {"Name": "Frank", "Age": 30},
    {"Name": "Grace", "Age": 26},
})
// 生成 SQL: INSERT INTO users (name, age) VALUES ('Frank', 30), ('Grace', 26)

// ===================== Upsert（插入或更新）=====================
import "gorm.io/gorm/clause"

// 场景：插入时如果 email 已存在（唯一索引冲突），则更新 name 和 age
db.Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "email"}},                          // 冲突判断列
    DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),         // 冲突时更新哪些列
}).Create(&users)
// MySQL 生成: INSERT INTO users (...) VALUES (...) ON DUPLICATE KEY UPDATE name=VALUES(name), age=VALUES(age)
// PostgreSQL 生成: INSERT INTO users (...) VALUES (...) ON CONFLICT (email) DO UPDATE SET name=EXCLUDED.name, age=EXCLUDED.age

// 冲突时什么都不做（忽略重复记录）
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
// MySQL 生成: INSERT IGNORE INTO users (...) VALUES (...)
// PostgreSQL 生成: INSERT INTO users (...) VALUES (...) ON CONFLICT DO NOTHING
```

---

## 5. 查询（Query）

```go
var user User
var users []User

// =============================================
//                基础查询
// =============================================

// --- 主键查询 ---
db.First(&user, 1)
// 生成 SQL: SELECT * FROM users WHERE id = 1 ORDER BY id LIMIT 1
// First 会按主键升序排序并取第一条
// ⚠️ 如果找不到记录，返回 gorm.ErrRecordNotFound

db.First(&user, "id = ?", 10)
// 生成 SQL: SELECT * FROM users WHERE id = 10 ORDER BY id LIMIT 1
// 和上面等价，只是用字符串条件写法

// 多主键查询
db.Find(&users, []int{1, 2, 3})
// 生成 SQL: SELECT * FROM users WHERE id IN (1, 2, 3)

// --- Take：取一条记录（不指定排序，由数据库决定顺序）---
db.Take(&user)
// 生成 SQL: SELECT * FROM users LIMIT 1
// 和 First 的区别：First 会加 ORDER BY id，Take 不加排序

// --- Last：取最后一条 ---
db.Last(&user)
// 生成 SQL: SELECT * FROM users ORDER BY id DESC LIMIT 1
// 按主键降序取第一条，即"最后一条"

// --- 查询全部 ---
db.Find(&users)
// 生成 SQL: SELECT * FROM users
// ⚠️ Find 找不到记录时不会报错，只是 users 为空切片，result.RowsAffected = 0

// =============================================
//              Where 条件
// =============================================

// --- 字符串条件 ---
db.Where("name = ?", "Alice").First(&user)
// 生成 SQL: SELECT * FROM users WHERE name = 'Alice' ORDER BY id LIMIT 1

db.Where("name <> ?", "Alice").Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name <> 'Alice'
// <> 等同于 !=

db.Where("name IN ?", []string{"Alice", "Bob"}).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name IN ('Alice', 'Bob')

db.Where("name LIKE ?", "%ali%").Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name LIKE '%ali%'
// 匹配任何包含 "ali" 的名字（不区分大小写取决于 DB 配置）

db.Where("age BETWEEN ? AND ?", 20, 30).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE age BETWEEN 20 AND 30
// 包含边界：age >= 20 AND age <= 30

db.Where("created_at > ?", time.Now().AddDate(0, 0, -7)).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE created_at > '2024-03-08 ...'
// 查询最近 7 天创建的记录

// 多个 Where 会用 AND 连接
db.Where("name = ?", "Alice").Where("age > ?", 18).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name = 'Alice' AND age > 18

// --- Struct 条件 ---
// ⚠️ 重要陷阱：零值字段会被自动忽略！
db.Where(&User{Name: "Alice", Age: 0}).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name = 'Alice'
// Age = 0 是 int 的零值，被 GORM 认为"没有赋值"，所以被忽略
// 如果你确实想查 age = 0 的记录，不要用 struct 条件，用 Map 或字符串

db.Where(&User{Name: "Alice", Age: 25}).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name = 'Alice' AND age = 25
// Age = 25 不是零值，所以正常参与查询

// --- Map 条件 ---
// Map 不存在"零值被忽略"的问题，所有 key 都会参与查询
db.Where(map[string]interface{}{"name": "Alice", "age": 0}).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name = 'Alice' AND age = 0
// Age = 0 也会参与查询条件

// =============================================
//              Select / 排序 / 分页
// =============================================

// --- Select 指定查询列 ---
db.Select("name", "age").Find(&users)
// 生成 SQL: SELECT name, age FROM users
// 只查 name 和 age 两列，其他字段为零值

db.Select("COALESCE(age, 0) as age").Find(&users)
// 生成 SQL: SELECT COALESCE(age, 0) as age FROM users
// COALESCE：如果 age 是 NULL 就返回 0

// --- 排序 ---
db.Order("age desc, name asc").Find(&users)
// 生成 SQL: SELECT * FROM users ORDER BY age DESC, name ASC
// 先按 age 降序，age 相同的再按 name 升序

db.Order("age desc").Order("name").Find(&users)
// 生成 SQL: SELECT * FROM users ORDER BY age DESC, name
// 多次 Order 调用会拼接

// --- 分页 ---
db.Offset(10).Limit(20).Find(&users)
// 生成 SQL: SELECT * FROM users LIMIT 20 OFFSET 10
// 跳过前 10 条，取接下来的 20 条（第 11~30 条）
// 常见用法：page=2, pageSize=20 → Offset((2-1)*20).Limit(20)

// 取消 Limit（传 -1）
db.Offset(10).Limit(-1).Find(&users)
// 生成 SQL: SELECT * FROM users OFFSET 10
// 从第 11 条开始取到末尾

// --- Distinct ---
db.Distinct("name").Find(&users)
// 生成 SQL: SELECT DISTINCT name FROM users
// 去重查询 name

// --- Count ---
var count int64
db.Model(&User{}).Where("age > ?", 18).Count(&count)
// 生成 SQL: SELECT COUNT(*) FROM users WHERE age > 18
// count 的值就是满足条件的记录数
// ⚠️ Count 需要指定 Model 或 Table，否则 GORM 不知道查哪张表

// --- Pluck：查询单列到切片 ---
var names []string
db.Model(&User{}).Pluck("name", &names)
// 生成 SQL: SELECT name FROM users
// names = ["Alice", "Bob", "Carol", ...]
// Pluck 只能查一列，结果直接存到基础类型切片

var ages []int
db.Model(&User{}).Pluck("age", &ages)
// 生成 SQL: SELECT age FROM users

// =============================================
//                Not / Or
// =============================================

db.Not("name = ?", "Alice").Find(&users)
// 生成 SQL: SELECT * FROM users WHERE NOT (name = 'Alice')
// 等价于 name <> 'Alice'

db.Not(map[string]interface{}{"name": []string{"Alice", "Bob"}}).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name NOT IN ('Alice', 'Bob')

db.Where("age > ?", 18).Or("name = ?", "Admin").Find(&users)
// 生成 SQL: SELECT * FROM users WHERE age > 18 OR name = 'Admin'
// 查询 age > 18 的用户，或者 name 为 Admin 的用户（Admin 不受 age 限制）

// 更复杂的 Or 组合
db.Where("age > ? AND active = ?", 18, true).Or("name = ?", "Admin").Find(&users)
// 生成 SQL: SELECT * FROM users WHERE (age > 18 AND active = true) OR name = 'Admin'

// =============================================
//         FirstOrInit / FirstOrCreate
// =============================================

// --- FirstOrInit ---
// 先尝试查找，找到就返回；找不到就在内存中用条件 + Attrs 初始化（不写入数据库）
db.Where(User{Name: "new_user"}).Attrs(User{Age: 20}).FirstOrInit(&user)
// 第 1 步：执行查询
//   生成 SQL: SELECT * FROM users WHERE name = 'new_user' ORDER BY id LIMIT 1
// 第 2 步：
//   如果找到了 → user = 数据库中的记录（忽略 Attrs）
//   如果没找到 → user = User{Name: "new_user", Age: 20}
//     ✅ Name 来自 Where 条件
//     ✅ Age 来自 Attrs
//     ⚠️ 此时 user 只在内存中，数据库里没有这条记录（ID = 0）

// 如果没找到且不使用 Attrs：
db.Where(User{Name: "new_user"}).FirstOrInit(&user)
// 没找到 → user = User{Name: "new_user", Age: 0}
// 只有 Where 条件里的字段被赋值

// --- FirstOrCreate ---
// 先尝试查找，找到就返回；找不到就创建并写入数据库
db.Where(User{Name: "new_user"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
// 第 1 步：执行查询
//   生成 SQL: SELECT * FROM users WHERE name = 'new_user' ORDER BY id LIMIT 1
// 第 2 步：
//   如果找到了 → user = 数据库中的记录（Attrs 被忽略，不修改已有记录）
//   如果没找到 → 执行 INSERT:
//     生成 SQL: INSERT INTO users (name, age, created_at, updated_at) VALUES ('new_user', 20, ...)
//     user = User{ID: <新生成>, Name: "new_user", Age: 20, CreatedAt: ..., UpdatedAt: ...}
//     ✅ Name 来自 Where 条件
//     ✅ Age 来自 Attrs
//     ✅ 记录已写入数据库

// --- Assign ---
// 无论是否找到，都会将 Assign 的值赋给 user
// FirstOrInit + Assign：赋值但不写库
db.Where(User{Name: "new_user"}).Assign(User{Age: 30}).FirstOrInit(&user)
//   找到了 → user = 数据库记录，但 user.Age 被覆盖为 30（仅内存中改，DB 不变）
//   没找到 → user = User{Name: "new_user", Age: 30}（仅内存）

// FirstOrCreate + Assign：赋值且写库
db.Where(User{Name: "new_user"}).Assign(User{Age: 30}).FirstOrCreate(&user)
//   找到了 → 执行 UPDATE，将 Age 更新为 30
//     生成 SQL: UPDATE users SET age = 30, updated_at = '...' WHERE id = <找到的ID>
//   没找到 → 执行 INSERT，Name = "new_user", Age = 30
//     生成 SQL: INSERT INTO users (name, age, ...) VALUES ('new_user', 30, ...)

// Attrs vs Assign 总结：
// Attrs:  仅在"创建"时生效，找到已有记录则忽略
// Assign: 无论"创建"还是"找到"都生效
```

---

## 6. 高级查询

```go
// =============================================
//                 子查询
// =============================================

// 先构造子查询（不会立即执行）
subQuery := db.Select("AVG(age)").Where("name LIKE ?", "A%").Table("users")
// subQuery 对应: SELECT AVG(age) FROM users WHERE name LIKE 'A%'

db.Where("age > (?)", subQuery).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE age > (SELECT AVG(age) FROM users WHERE name LIKE 'A%')
// 查询年龄大于"所有 A 开头用户的平均年龄"的用户

// FROM 子查询：把子查询结果当作临时表
db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age > ?", 18).Find(&users)
// 生成 SQL: SELECT * FROM (SELECT name, age FROM users) as u WHERE age > 18

// =============================================
//              Group / Having
// =============================================

type Result struct {
    Date  string
    Total int64
}
var results []Result

db.Model(&Order{}).
    Select("DATE(created_at) as date, SUM(amount) as total").
    Group("DATE(created_at)").
    Having("SUM(amount) > ?", 100).
    Scan(&results)
// 生成 SQL:
//   SELECT DATE(created_at) as date, SUM(amount) as total
//   FROM orders
//   GROUP BY DATE(created_at)
//   HAVING SUM(amount) > 100
// 按日期分组统计订单总额，只返回总额 > 100 的日期
// 结果存到自定义的 Result 结构体（用 Scan 而非 Find，因为 Result 不对应任何表）

// =============================================
//                  Joins
// =============================================

type UserWithOrder struct {
    UserName   string
    OrderTotal float64
}
var joinResults []UserWithOrder

db.Model(&User{}).
    Select("users.name as user_name, SUM(orders.amount) as order_total").
    Joins("LEFT JOIN orders ON orders.user_id = users.id").
    Group("users.id").
    Scan(&joinResults)
// 生成 SQL:
//   SELECT users.name as user_name, SUM(orders.amount) as order_total
//   FROM users
//   LEFT JOIN orders ON orders.user_id = users.id
//   GROUP BY users.id
// 查询每个用户的订单总额（LEFT JOIN：没有订单的用户也会出现，order_total 为 NULL）

// 带条件的 Joins
db.Model(&User{}).
    Joins("INNER JOIN orders ON orders.user_id = users.id AND orders.amount > ?", 100).
    Find(&users)
// 生成 SQL:
//   SELECT users.* FROM users
//   INNER JOIN orders ON orders.user_id = users.id AND orders.amount > 100
// INNER JOIN：只返回有 amount > 100 订单的用户

// =============================================
//            Scan 到自定义结构体
// =============================================

// 当查询结果不直接对应某个 Model 时，用 Scan 代替 Find
type SimpleUser struct {
    Name string
    Age  int
}
var simpleUsers []SimpleUser

db.Model(&User{}).Select("name", "age").Scan(&simpleUsers)
// 生成 SQL: SELECT name, age FROM users
// 结果映射到 SimpleUser（只有 Name 和 Age 两个字段）

// =============================================
//     FindInBatches：分批处理大量数据
// =============================================

// 场景：需要处理百万级记录，一次全加载到内存会 OOM
// FindInBatches 每次只查 100 条，处理完再查下一批
db.Where("active = ?", true).FindInBatches(&users, 100, func(tx *gorm.DB, batch int) error {
    // batch 是第几批（从 1 开始）
    // users 是当前批次的数据（最多 100 条）
    // tx 是当前批次的 DB 对象（可以用来做 Count 等）
    for _, user := range users {
        // 处理每条记录，比如发邮件、同步数据等
        _ = user
    }
    fmt.Printf("第 %d 批处理了 %d 条\n", batch, tx.RowsAffected)
    // 返回 error 会停止后续批次
    // 返回 nil 继续处理下一批
    return nil
})
// 生成 SQL:
//   第 1 批: SELECT * FROM users WHERE active = true LIMIT 100
//   第 2 批: SELECT * FROM users WHERE active = true LIMIT 100 OFFSET 100
//   第 3 批: SELECT * FROM users WHERE active = true LIMIT 100 OFFSET 200
//   ... 直到没有更多数据

// =============================================
//              Locking（锁）
// =============================================

import "gorm.io/gorm/clause"

// 悲观锁 - 排他锁（写锁）
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
// 生成 SQL: SELECT * FROM users FOR UPDATE
// 其他事务不能读也不能写这些行，直到当前事务提交

// 悲观锁 - 共享锁（读锁）
db.Clauses(clause.Locking{Strength: "SHARE"}).Find(&users)
// 生成 SQL: SELECT * FROM users FOR SHARE
// 其他事务可以读但不能写这些行

// =============================================
//            Optimizer / Index Hints
// =============================================

import "gorm.io/hints"

// 建议数据库使用某个索引
db.Clauses(hints.UseIndex("idx_name")).Find(&users)
// 生成 SQL: SELECT * FROM users USE INDEX (idx_name)

// 强制使用某个索引
db.Clauses(hints.ForceIndex("idx_name")).Where("name = ?", "Alice").Find(&users)
// 生成 SQL: SELECT * FROM users FORCE INDEX (idx_name) WHERE name = 'Alice'

// =============================================
//               命名参数
// =============================================

import "database/sql"

db.Where("name = @name OR age = @age", sql.Named("name", "Alice"), sql.Named("age", 18)).Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name = 'Alice' OR age = 18
// 用 @name、@age 替代 ?，更清晰，适合条件很多的场景

// =============================================
//              Map 结果
// =============================================

// 不想定义结构体，直接用 map 接收结果
var resultMap []map[string]interface{}
db.Model(&User{}).Find(&resultMap)
// 生成 SQL: SELECT * FROM users
// resultMap 类似: [{"id": 1, "name": "Alice", "age": 25}, {"id": 2, ...}, ...]
// 每条记录是一个 map，key 是列名，value 是 interface{}
```

---

## 7. 更新（Update）

```go
// 假设已有: user = User{ID: 1, Name: "Alice", Age: 25, Email: "alice@go.dev"}

// =============================================
//              单字段更新
// =============================================

db.Model(&user).Update("name", "new_name")
// 生成 SQL: UPDATE users SET name = 'new_name', updated_at = '2024-...' WHERE id = 1
// ⚠️ 必须通过 Model 指定是哪条记录（user.ID = 1 → WHERE id = 1）
// ⚠️ updated_at 自动更新

// 条件更新（不依赖主键，用 Where 指定条件）
db.Model(&User{}).Where("age < ?", 18).Update("active", false)
// 生成 SQL: UPDATE users SET active = false, updated_at = '...' WHERE age < 18
// 将所有 age < 18 的用户设为 inactive

// =============================================
//           多字段更新（Struct）
// =============================================

// ⚠️ 重要陷阱：struct 只更新非零值字段！零值字段（0、""、false）会被忽略
db.Model(&user).Updates(User{Name: "hello", Age: 0})
// 生成 SQL: UPDATE users SET name = 'hello', updated_at = '...' WHERE id = 1
// ❌ Age = 0 是 int 零值，被 GORM 忽略，不会更新！
// 如果你想把 Age 更新为 0，不要用 struct，用 Map 或 Select

db.Model(&user).Updates(User{Name: "hello", Age: 30})
// 生成 SQL: UPDATE users SET name = 'hello', age = 30, updated_at = '...' WHERE id = 1
// Age = 30 不是零值，正常更新

// =============================================
//           多字段更新（Map）
// =============================================

// Map 可以更新零值，不存在"忽略"问题
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 0})
// 生成 SQL: UPDATE users SET name = 'hello', age = 0, updated_at = '...' WHERE id = 1
// ✅ name 和 age 都会更新，包括 age = 0

// =============================================
//       Select / Omit 控制更新字段
// =============================================

// Select：只更新指定的字段（即使是零值也会更新）
db.Model(&user).Select("Name", "Age").Updates(User{Name: "new", Age: 0})
// 生成 SQL: UPDATE users SET name = 'new', age = 0, updated_at = '...' WHERE id = 1
// ✅ 强制更新 Name 和 Age，即使 Age = 0

// Select("*")：更新所有字段（包括零值），效果类似 Save
db.Model(&user).Select("*").Updates(User{Name: "new"})
// 生成 SQL: UPDATE users SET name='new', age=0, email='', active=false, updated_at='...' WHERE id = 1
// ⚠️ 所有未赋值的字段都会被更新为零值！慎用

// Omit：跳过指定字段，其余非零值字段都更新
db.Model(&user).Omit("Name").Updates(User{Name: "ignored", Age: 30})
// 生成 SQL: UPDATE users SET age = 30, updated_at = '...' WHERE id = 1
// Name 被忽略，即使传了值

// =============================================
//              SQL 表达式
// =============================================

// 在现有值的基础上做计算（不是直接赋值）
db.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 10))
// 生成 SQL: UPDATE products SET price = price * 2 + 10, updated_at = '...' WHERE id = 1
// 原来 price = 50，执行后 price = 50 * 2 + 10 = 110

db.Model(&product).Updates(map[string]interface{}{
    "price": gorm.Expr("price + ?", 5),   // price = price + 5
    "stock": gorm.Expr("stock - ?", 1),   // stock = stock - 1
})
// 生成 SQL: UPDATE products SET price = price + 5, stock = stock - 1, updated_at = '...' WHERE id = 1
// 适合"扣库存"、"加积分"等场景，避免先 SELECT 再 UPDATE 的并发问题

// =============================================
//              批量更新
// =============================================

db.Model(&User{}).Where("active = ?", false).Updates(map[string]interface{}{"deleted": true})
// 生成 SQL: UPDATE users SET deleted = true, updated_at = '...' WHERE active = false
// 把所有 active=false 的用户标记为 deleted

// =============================================
//      Save：全量保存所有字段（包括零值）
// =============================================

user.Name = ""    // 空字符串
user.Age = 0      // 零值
db.Save(&user)
// 生成 SQL: UPDATE users SET name='', age=0, email='alice@go.dev', active=true, updated_at='...' WHERE id=1
// ⚠️ Save 会更新所有字段（包括零值），相当于 struct 的完整覆盖
// ⚠️ 如果 user.ID = 0（没有主键），Save 会执行 INSERT 而非 UPDATE
// 建议：除非确定要更新所有字段，否则用 Updates 更安全

// =============================================
//  不触发 Hook 和时间追踪的更新
// =============================================

// UpdateColumn / UpdateColumns 不会：
//   ❌ 触发 BeforeUpdate / AfterUpdate 钩子
//   ❌ 自动更新 updated_at
// 适合纯粹的数据修正、数据迁移等场景

db.Model(&User{}).Where("id = ?", 1).UpdateColumn("name", "silent")
// 生成 SQL: UPDATE users SET name = 'silent' WHERE id = 1
// 注意没有 updated_at = '...'

db.Model(&User{}).Where("id = ?", 1).UpdateColumns(map[string]interface{}{"name": "silent", "age": 30})
// 生成 SQL: UPDATE users SET name = 'silent', age = 30 WHERE id = 1
```

---

## 8. 删除（Delete）

```go
// =============================================
//         软删除（模型含 DeletedAt 时自动启用）
// =============================================

// 如果 User 嵌入了 gorm.Model（含 DeletedAt 字段），Delete 不会真正删除行
// 而是将 deleted_at 设为当前时间
db.Delete(&user)
// 生成 SQL: UPDATE users SET deleted_at = '2024-03-15 ...' WHERE id = 1
// 记录还在数据库中，只是被标记为"已删除"

db.Delete(&user, 1)
// 生成 SQL: UPDATE users SET deleted_at = '...' WHERE id = 1
// 通过主键指定要删除的记录

db.Delete(&User{}, []int{1, 2, 3})
// 生成 SQL: UPDATE users SET deleted_at = '...' WHERE id IN (1, 2, 3)
// 批量软删除

// 条件删除
db.Where("age < ?", 10).Delete(&User{})
// 生成 SQL: UPDATE users SET deleted_at = '...' WHERE age < 10 AND deleted_at IS NULL

// ===================== 软删除的查询行为 =====================

// 默认查询自动排除已软删除的记录
db.Find(&users)
// 生成 SQL: SELECT * FROM users WHERE deleted_at IS NULL
// 不会查到已删除的记录

db.Where("name = ?", "Alice").First(&user)
// 生成 SQL: SELECT * FROM users WHERE name = 'Alice' AND deleted_at IS NULL ORDER BY id LIMIT 1

// 查找包含已软删除记录
db.Unscoped().Where("name = ?", "Alice").Find(&users)
// 生成 SQL: SELECT * FROM users WHERE name = 'Alice'
// Unscoped() 去掉 deleted_at IS NULL 条件，能查到已删除的

// ===================== 永久删除（物理删除）=====================

db.Unscoped().Delete(&user)
// 生成 SQL: DELETE FROM users WHERE id = 1
// 真正从数据库删除这一行，不可恢复

// ===================== 防止误删全表 =====================

// ⚠️ 不带条件的删除会被 GORM 阻止，返回 ErrMissingWhereClause
// db.Delete(&User{})  → 报错！

// 如果确实需要删除全部记录：
db.Where("1 = 1").Delete(&User{})
// 生成 SQL: UPDATE users SET deleted_at = '...' WHERE 1 = 1 AND deleted_at IS NULL

// 或者在配置中全局允许（不推荐）：
// gorm.Config{AllowGlobalUpdate: true}
```

---

## 9. 原生 SQL

```go
// ===================== Raw 查询 =====================

// 执行原生 SELECT，结果映射到结构体
var users []User
db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)
// 直接执行: SELECT * FROM users WHERE age > 18
// 结果扫描到 []User

// 映射到自定义结构体
var result struct {
    Total int64
    Avg   float64
}
db.Raw("SELECT COUNT(*) as total, AVG(age) as avg FROM users").Scan(&result)
// result.Total = 用户总数, result.Avg = 平均年龄

// ===================== Exec 执行（非查询语句）=====================

db.Exec("UPDATE users SET name = ? WHERE id = ?", "new_name", 1)
// 执行原生 UPDATE

db.Exec("DROP TABLE IF EXISTS temp_users")
// 执行原生 DDL

// ===================== SQL Builder =====================

// 不写完整 SQL，用 GORM 的链式方法拼接
var results []map[string]interface{}
db.Table("users").
    Select("name", "age").
    Where("id > ?", 5).
    Scan(&results)
// 生成 SQL: SELECT name, age FROM users WHERE id > 5
// Table("users") → 指定表名（不依赖 Model）

// ===================== Row / Rows =====================

// 查询单行
row := db.Table("users").Where("name = ?", "Alice").Select("name", "age").Row()
// 生成 SQL: SELECT name, age FROM users WHERE name = 'Alice'

var name string
var age int
row.Scan(&name, &age) // 手动扫描每个列到变量

// 查询多行
rows, _ := db.Model(&User{}).Where("age > ?", 18).Rows()
// 生成 SQL: SELECT * FROM users WHERE age > 18
defer rows.Close() // ⚠️ 必须关闭，否则连接泄漏

for rows.Next() {
    var user User
    db.ScanRows(rows, &user) // 将当前行扫描到 User 结构体（比手动 Scan 方便）
    fmt.Println(user)
}

// ===================== DryRun 模式 =====================

// 只生成 SQL，不真正执行，用于调试
stmt := db.Session(&gorm.Session{DryRun: true}).
    Where("name = ?", "Alice").
    Find(&User{}).Statement

fmt.Println(stmt.SQL.String())   // "SELECT * FROM users WHERE name = ? AND deleted_at IS NULL"
fmt.Println(stmt.Vars)           // ["Alice"]
// 可以看到生成的 SQL 和绑定参数，不会实际查数据库
```

---

## 10. 关联关系

```go
// =============================================
//              Belongs To（属于）一对一/一对多
// =============================================

// 一个 Order 属于一个 User
// 外键放在 Order 表中（"子"表持有外键）
type Order struct {
    gorm.Model
    Amount float64
    UserID uint    // 外键，默认命名规则：关联结构体名 + "ID" → UserID
    User   User    // 关联字段：通过 UserID 找到对应的 User
    // 建表后 orders 表有 user_id 列，指向 users 表的 id
}

// 自定义外键名
type Order2 struct {
    gorm.Model
    Amount    float64
    CreatorID uint                              // 外键字段叫 CreatorID 而非 UserID
    Creator   User `gorm:"foreignKey:CreatorID"` // 告诉 GORM：用 CreatorID 关联 User
    // 建表后 orders 表有 creator_id 列
}

// =============================================
//              Has One（拥有一个）1对1
// =============================================

// 一个 User 拥有一个 Profile
// 外键在 Profile 表中（"子"表持有外键）
type User struct {
    gorm.Model
    Name    string
    Profile Profile // has one：User 拥有一个 Profile
}
type Profile struct {
    gorm.Model
    UserID uint   // 外键：指向 users 表的 id
    Bio    string
    Avatar string
    // 建表后 profiles 表有 user_id 列
}
// Has One vs Belongs To 的区别：
//   Has One:    外键在"对方"表（Profile 表有 user_id）
//   Belongs To: 外键在"自己"表（Order 表有 user_id）
//   本质上是从哪个方向看的区别

// =============================================
//            Has Many（拥有多个）一对多 
// =============================================

// 一个 User 拥有多个 Order
type User struct {
    gorm.Model
    Name   string
    Orders []Order // has many：用切片表示"多个"
}
// 外键在 Order 表中（Order.UserID → User.ID）
// GORM 自动根据 User → Order 的关系找到 Order.UserID 作为外键

// =============================================
//          Many To Many（多对多）
// =============================================

// 一个 User 可以学多门语言，一门语言可以被多个 User 学习
type User struct {
    gorm.Model
    Name      string
    Languages []Language `gorm:"many2many:user_languages;"`
    // many2many:user_languages → 指定中间表名为 user_languages
    // GORM 会自动创建中间表 user_languages，包含 user_id 和 language_id 两列
}
type Language struct {
    gorm.Model
    Name string
}
// 中间表结构（GORM 自动创建）：
//   user_languages:
//     user_id     → 外键指向 users.id
//     language_id → 外键指向 languages.id

// 自定义中间表（需要额外字段时）
type UserLanguage struct {
    UserID     uint `gorm:"primaryKey"`
    LanguageID uint `gorm:"primaryKey"`
    Skill      int  // 额外字段：熟练度
}
// 需要用 SetupJoinTable 注册自定义中间表
db.SetupJoinTable(&User{}, "Languages", &UserLanguage{})

// =============================================
//             多态关联
// =============================================

// Comment 可以属于 Article 或 Video（不同类型的"父"表）
type Comment struct {
    gorm.Model
    Content         string
    CommentableID   uint   // 父记录的 ID（如 Article 的 ID 或 Video 的 ID）
    CommentableType string // 父记录的类型（"articles" 或 "videos"）
}
type Article struct {
    gorm.Model
    Title    string
    Comments []Comment `gorm:"polymorphic:Commentable;"`
    // polymorphic:Commentable → 使用 CommentableID + CommentableType 作为多态外键
}
type Video struct {
    gorm.Model
    URL      string
    Comments []Comment `gorm:"polymorphic:Commentable;"`
}
// 当为 Article（ID=1）创建 Comment 时：
//   CommentableID = 1, CommentableType = "articles"
// 当为 Video（ID=5）创建 Comment 时：
//   CommentableID = 5, CommentableType = "videos"

// =============================================
//       Association Mode（关联操作）
// =============================================

// 创建时自动写关联
db.Create(&User{
    Name:   "Alice",
    Orders: []Order{{Amount: 100}, {Amount: 200}},
})
// 生成 SQL:
//   INSERT INTO users (name, ...) VALUES ('Alice', ...)
//   INSERT INTO orders (user_id, amount, ...) VALUES (1, 100, ...), (1, 200, ...)
// 自动将 user_id 填入 orders 表

// 添加关联（给已有 user 追加一个 order）
db.Model(&user).Association("Orders").Append(&Order{Amount: 300})
// 生成 SQL: INSERT INTO orders (user_id, amount, ...) VALUES (1, 300, ...)

// 替换关联（清空原有关联，用新的替换）
db.Model(&user).Association("Languages").Replace(&Language{Name: "Go"}, &Language{Name: "Rust"})
// 先清除 user_languages 中该 user 的所有记录，再插入新的关联

// 删除关联（解除关系，不删除记录本身）
db.Model(&user).Association("Orders").Delete(&order)
// Has Many/Many2Many: 清除外键（user_id 设为 NULL 或删除中间表记录）
// ⚠️ 不会删除 orders 表中的记录本身

// 清空关联（解除该 user 的所有 orders 关联）
db.Model(&user).Association("Orders").Clear()
// 将所有关联 order 的 user_id 设为 NULL（或删除中间表记录）

// 关联计数
count := db.Model(&user).Association("Orders").Count()
// 生成 SQL: SELECT COUNT(*) FROM orders WHERE user_id = 1
// count = 该 user 有多少个 order

// 查找关联
var orders []Order
db.Model(&user).Association("Orders").Find(&orders)
// 生成 SQL: SELECT * FROM orders WHERE user_id = 1
// orders = 该 user 的所有订单

// 带条件查找关联
db.Model(&user).Association("Orders").Find(&orders, "amount > ?", 50)
// 生成 SQL: SELECT * FROM orders WHERE user_id = 1 AND amount > 50
```

---
