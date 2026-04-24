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
// 安装
// go get -u gorm.io/gorm
// go get -u gorm.io/driver/mysql
// go get -u grom.io/driver/postgres
// go get -u gorm.io/driver/sqlite

package main

import (
    "time"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)