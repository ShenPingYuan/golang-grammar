# Gin 框架进阶系列（五）：Gin + GORM 连接数据库实现 CRUD

---

## GORM 简介

GORM 是 Go 生态中最流行的 ORM 库，功能覆盖自动迁移、关联关系、事务、钩子、预加载等。Gin 本身不绑定任何 ORM，但 Gin + GORM 是社区中最主流的组合。

> 本篇文章只是用于演示 Gin 与 GORM 的简单搭配使用，如需用于生产环境很多细节还需完善，比如应该返回 UserDto 而不是直接返回数据库模型 User 。

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql    # MySQL
go get -u gorm.io/driver/postgres # PostgreSQL
go get -u gorm.io/driver/sqlite   # SQLite（本地开发方便）
```

本篇以 MySQL 为主，其他数据库只需换驱动，API 完全一致。

---

## 项目结构

```
├── main.go
├── config/
│   └── database.go    // 数据库初始化
├── model/
│   └── user.go        // 模型定义
├── handler/
│   └── user.go        // 路由处理函数
└── go.mod
```

微小型项目按这个分层足够。核心思路：**model 不依赖 Gin，handler 不写 SQL**。

---

## 数据库初始化

```go
// config/database.go
package config

import (
    "fmt"
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        "root",      // 用户名
        "123456",    // 密码
        "127.0.0.1", // 主机
        3306,        // 端口
        "gin_demo",  // 数据库名
    )

    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info), // 开发阶段打印 SQL
    })
    if err != nil {
        log.Fatalf("数据库连接失败: %v", err)
    }

    // 连接池配置
    sqlDB, _ := DB.DB()
    sqlDB.SetMaxIdleConns(10) // 设置最大空闲连接数
    sqlDB.SetMaxOpenConns(100) // 设置最大打开连接数


    log.Println("数据库连接成功")
}
```

生产环境中 DSN 应从环境变量或配置文件读取，不要硬编码。`logger.Info` 会打印每条 SQL 及耗时，上线时切换为 `logger.Warn` 或 `logger.Silent`。

---

## 模型定义

```go
// model/user.go
package model

import "gorm.io/gorm"

type User struct {
    gorm.Model         // 内嵌 ID、CreatedAt、UpdatedAt、DeletedAt
    Name   string `gorm:"type:varchar(50);not null"  json:"name"`
    Email  string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
    Age    int    `gorm:"default:0"                  json:"age"`
}
```

`gorm.Model` 展开后是四个字段：

```go
type Model struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除
}
```

有了 `DeletedAt`，调用 `Delete` 时不会真正删除记录，而是填充删除时间，查询时自动过滤。

---

## 自动迁移

```go
// main.go 中初始化完数据库后
config.InitDB()
config.DB.AutoMigrate(&model.User{})
```

`AutoMigrate` 会根据 struct 创建表、添加缺失的列和索引。它**不会删除列或修改列类型**，生产环境建议用专业的迁移工具（如 [golang-migrate](https://github.com/golang-migrate/migrate)）管理 schema 变更。

---

## 请求与响应结构体分离

不要直接用 model 接收前端参数，原因有三：model 包含 `ID`、`DeletedAt` 等不该由前端控制的字段；校验规则和数据库约束是两码事；响应时可能需要隐藏某些字段。

```go
// handler/user.go
package handler

// 创建用户请求
type CreateUserReq struct {
    Name  string `json:"name"  binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age"   binding:"gte=0,lte=150"`
}

// 更新用户请求
type UpdateUserReq struct {
    Name  string `json:"name"  binding:"omitempty"`
    Email string `json:"email" binding:"omitempty,email"`
    Age   *int   `json:"age"   binding:"omitempty,gte=0,lte=150"`
}

// 统一响应
type Response struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}
```

`UpdateUserReq` 中 `Age` 用指针 `*int`，这样可以区分"没传"和"传了 0"。

---

## Create：创建用户

```go
func CreateUser(c *gin.Context) {
    var req CreateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, Response{Code: 400, Msg: err.Error()})
        return
    }

    user := model.User{
        Name:  req.Name,
        Email: req.Email,
        Age:   req.Age,
    }

    if err := config.DB.Create(&user).Error; err != nil {
        // 唯一索引冲突判断
        if strings.Contains(err.Error(), "Duplicate") {
            c.JSON(409, Response{Code: 409, Msg: "邮箱已存在"})
            return
        }
        c.JSON(500, Response{Code: 500, Msg: "创建失败"})
        return
    }

    c.JSON(201, Response{Code: 0, Msg: "创建成功", Data: user})
}
```

`Create` 执行后，`user.ID` 会自动被 GORM 回填。

---

## Read：查询用户

```go
// 按 ID 查询
func GetUser(c *gin.Context) {
    var uri struct {
        ID uint `uri:"id" binding:"required,gte=1"`
    }
    if err := c.ShouldBindUri(&uri); err != nil {
        c.JSON(400, Response{Code: 400, Msg: err.Error()})
        return
    }

    var user model.User
    if err := config.DB.First(&user, uri.ID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, Response{Code: 404, Msg: "用户不存在"})
            return
        }
        c.JSON(500, Response{Code: 500, Msg: "查询失败"})
        return
    }

    c.JSON(200, Response{Code: 0, Data: user})
}
```

```go
// 分页列表
func ListUsers(c *gin.Context) {
    var query struct {
        Page int    `form:"page" binding:"gte=1"`
        Size int    `form:"size" binding:"gte=1,lte=100"`
        Name string `form:"name"`
    }
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(400, Response{Code: 400, Msg: err.Error()})
        return
    }

    // 默认值
    if query.Page == 0 {
        query.Page = 1
    }
    if query.Size == 0 {
        query.Size = 10
    }

    var users []model.User
    var total int64

    db := config.DB.Model(&model.User{})

    // 条件查询
    if query.Name != "" {
        db = db.Where("name LIKE ?", "%"+query.Name+"%")
    }

    db.Count(&total)
    db.Offset((query.Page - 1) * query.Size).Limit(query.Size).Find(&users)

    c.JSON(200, Response{
        Code: 0,
        Data: gin.H{
            "list":  users,
            "total": total,
            "page":  query.Page,
            "size":  query.Size,
        },
    })
}
```

注意 `Count` 要在 `Offset/Limit` 之前调用，否则 count 也会受分页影响。实际上 GORM 的 `Count` 会忽略 `Offset` 和 `Limit`，但把 `Count` 放前面语义更清晰，也避免某些边界情况。

---

## Update：更新用户

```go
func UpdateUser(c *gin.Context) {
    var uri struct {
        ID uint `uri:"id" binding:"required,gte=1"`
    }
    if err := c.ShouldBindUri(&uri); err != nil {
        c.JSON(400, Response{Code: 400, Msg: err.Error()})
        return
    }

    var req UpdateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, Response{Code: 400, Msg: err.Error()})
        return
    }

    // 先查用户是否存在
    var user model.User
    if err := config.DB.First(&user, uri.ID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, Response{Code: 404, Msg: "用户不存在"})
            return
        }
        c.JSON(500, Response{Code: 500, Msg: "查询失败"})
        return
    }

    // 用 map 做部分更新，只更新前端传了的字段
    updates := make(map[string]interface{})
    if req.Name != "" {
        updates["name"] = req.Name
    }
    if req.Email != "" {
        updates["email"] = req.Email
    }
    if req.Age != nil {
        updates["age"] = *req.Age
    }

    if len(updates) == 0 {
        c.JSON(400, Response{Code: 400, Msg: "没有需要更新的字段"})
        return
    }

    if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
        c.JSON(500, Response{Code: 500, Msg: "更新失败"})
        return
    }

    c.JSON(200, Response{Code: 0, Msg: "更新成功", Data: user})
}
```

为什么用 `map` 而不是直接传 struct？因为 GORM 的 `Updates` 接收 struct 时会**忽略零值字段**，如果用户想把 `Age` 更新为 0，struct 方式会跳过它。用 map 可以精确控制更新哪些字段。

---

## Delete：删除用户

```go
func DeleteUser(c *gin.Context) {
    var uri struct {
        ID uint `uri:"id" binding:"required,gte=1"`
    }
    if err := c.ShouldBindUri(&uri); err != nil {
        c.JSON(400, Response{Code: 400, Msg: err.Error()})
        return
    }

    result := config.DB.Delete(&model.User{}, uri.ID)
    if result.Error != nil {
        c.JSON(500, Response{Code: 500, Msg: "删除失败"})
        return
    }
    if result.RowsAffected == 0 {
        c.JSON(404, Response{Code: 404, Msg: "用户不存在"})
        return
    }

    c.JSON(200, Response{Code: 0, Msg: "删除成功"})
}
```

因为 model 中有 `gorm.DeletedAt`，这里执行的是**软删除**——`UPDATE users SET deleted_at = NOW() WHERE id = ?`。查询时 GORM 自动加 `WHERE deleted_at IS NULL`，已删除的记录对业务层透明。

如果确实需要物理删除：

```go
config.DB.Unscoped().Delete(&model.User{}, uri.ID)
```

---

## 注册路由

```go
// main.go
package main

import (
    "your-project/config"
    "your-project/handler"
    "your-project/model"

    "github.com/gin-gonic/gin"
)

func main() {
    config.InitDB()
    config.DB.AutoMigrate(&model.User{})

    r := gin.Default()

    v1 := r.Group("/api/v1")
    {
        v1.POST("/users", handler.CreateUser)
        v1.GET("/users", handler.ListUsers)
        v1.GET("/users/:id", handler.GetUser)
        v1.PUT("/users/:id", handler.UpdateUser)
        v1.DELETE("/users/:id", handler.DeleteUser)
    }

    r.Run(":8080")
}
```

---

## 测试一下

```bash
# 创建
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"tom","email":"tom@example.com","age":25}'

# 列表
curl "http://localhost:8080/api/v1/users?page=1&size=10&name=tom"

# 查询
curl http://localhost:8080/api/v1/users/1

# 更新
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"jerry"}'

# 删除
curl -X DELETE http://localhost:8080/api/v1/users/1
```

---

## 常见踩坑

**全局 DB 变量并发安全吗？** 安全。`*gorm.DB` 内部维护了连接池（`database/sql`），每次操作从池中取连接，天然支持并发。全局变量只是持有连接池的引用，不会有竞态问题。

**`First` 和 `Find` 的区别？** `First` 找不到记录返回 `gorm.ErrRecordNotFound`，`Find` 找不到返回空切片且 `Error` 为 nil。查单条用 `First`，查列表用 `Find`。

**更新时 `Updates` 和 `Save` 的区别？** `Save` 会更新所有字段（包括零值），`Updates` 只更新非零值字段（struct）或指定字段（map）。部分更新永远用 `Updates`。

**为什么不在 handler 里直接写复杂 SQL？** 当业务逻辑变复杂，应该再抽一层 service 或 repository。本篇为了演示清晰省略了这一层，后续架构篇会补上。

---

## 小结

这篇完成了从数据库连接到完整 CRUD 的全流程。核心要点：请求结构体与 model 分离、用 map 做部分更新、理解软删除机制、善用 `First` 与 `Find` 的语义差异。

下一篇进入**统一响应封装与错误处理**：自定义错误码体系、全局错误捕获、panic recovery 与业务错误的优雅处理。

---