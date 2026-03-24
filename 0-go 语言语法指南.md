# Go 语言语法完全指南

## 前言

Go，也称 Golang，Google 推出的一门静态强类型、编译型语言。特点：**语法简介、编译快、并发能力强、工程化支持完善**。应用领域如：后端开发、云原生、微服务、分布式系统、DevOps 工具链等。

本文是我学习 Golang 后，以 **系统性总结** 的方式对其语法知识进行完整梳理，可作为：

- Go 初学者系统入门笔记
- Go 语法速查手册
- 编写 Go 项目时基础知识参考

---

## 目录

- [1. Go 语言简介](#1-go-语言简介)
- [2. Go 程序基本结构](#2-go-程序基本结构)
- [3. 基础语法](#3-基础语法)
- [4. 流程控制](#4-流程控制)
- [5. 函数](#5-函数)



## 1. Go 语言简介

### 1.1 语言特点

- **语法简介**: Go 语法简洁。借鉴了 C 语言的风格，但去掉了许多复杂的特性。
- **编译型语言**：编译速度极快，生成单一二进制文件。
- **天生并发**：内置 goroutine 和 channel 简化并发编程。
- **垃圾回收**：无需手动管理内存。
- **静态类型 + 类型推断**：既有静态类型的安全性，又有类型推断的便利性。
- **标准库丰富**：网络、文件、编码、加密、测试等开箱即用
- **工程化能力强**：内置格式化、测试、模块管理工具。
- **跨平台**：支持多种操作系统和架构，这没啥好说的市面上大部分语言都支持。

### 1.2 应用场景

- Web 后端开发
- API 服务开发
- 云原生与容器生态
- 微服务架构
- 网络编程
- 运维工具和命令行工具
- 分布式系统

### 1.3 环境搭建

```bash
# 下载安装：https://go.dev/dl/

# 验证安装
go version

# 如果安装成功，会输出类似信息：
go version go1.22.0 windows/amd64
```

### 1.4 Hello World

创建文件 `main.go`，内容如下：

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

```bash
go run main.go     # 直接运行
go build main.go   # 编译生成二进制
```

## 2. Go 程序基本结构

Go 程序通常由以下几个部分组成：

- **包声明**
- **导入包**
- **全局声明**
- **函数定义**
- **main 函数**

### 2.1 package

每个 Go 文件都必须以 `package` 开头，用于声明该文件属于哪个包。

```go
package main
```
如果一个程序要作为可执行程序运行，则必须使用 `main` 包。

### 2.2 import

使用 `import` 导入依赖包。

```go
import "fmt"
```

如果导入多个包，可以使用圆括号：

```go
import (
	"fmt"
	"time"
)
```

## 2.3 main 函数

Go 程序的入口是 `main` 包中的 `main()` 函数。

```go
func main() {
	fmt.Println("程序开始执行")
}
```

## 2.4 注释

Go 支持两种注释方式：

### 单行注释

```go
// 这是单行注释
```

### 多行注释

```go
/*
这是多行注释
可以写多行内容
*/
```

## 2.5 标识符命名规则

标识符包括变量名、函数名、类型名等，命名规则如下：

- 由字母、数字、下划线组成
- 不能以数字开头
- 区分大小写
- 不能与关键字重名

例如：

```go
var name string
var userAge int
var _temp int
```

## 2.6 可见性规则

Go 使用**首字母大小写**控制访问权限：

- **首字母大写**：可被包外访问
- **首字母小写**：仅包内可访问

例如：

```go
type User struct {
	Name string
	age  int
}
```

其中：

- `Name` 可导出
- `age` 不可导出

---

## 3. 基础语法

### 3.1 变量

```go
// 显式声明
var name string = "Go"
var age int      // 零值：0

// 类型推断
var score = 100

// 短变量声明（只能在函数内部使用）
lang := "Golang"

// 批量声明
var a, b, c int
var x, y = 10, 20

// 也可以分组声明：
var (
	name string = "Alice"
	age  int = 20
	city string = "Beijing"
    gender bool  // 零值：false
)

// Go 使用 `_` 表示匿名变量，用于忽略某个值，常见于函数多返回值场景。
n, _ := 10, 20
fmt.Println(n)
```

**零值规则**：Go 中变量声明后自动初始化为零值。

| 类型                            | 零值    |
| ------------------------------- | ------- |
| int/float                       | `0`     |
| string                          | `""`    |
| bool                            | `false` |
| 指针/切片/map/channel/接口/函数 | `nil`   |

### 3.2 常量

```go
const Pi = 3.14159
const (
    StatusOK    = 200
    StatusNotFound = 404
)

// iota：常量生成器，从 0 开始自增
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
)

// iota 高级用法
const (
    _  = iota             // 0 丢弃
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20
    GB                    // 1 << 30
)
```

## 3.3 基本数据类型

**示例：**

```go
// 整型，对应的无符号整型：- uint uint8 uint16 uint32 uint64
var a int     // 平台相关：32 或 64 位
var b int8    // -128 ~ 127
var c int16
var d int32   // 别名 rune（表示 Unicode 码点）
var e int64
var f uint8   // 别名 byte（0 ~ 255）

// 浮点
var g float32
var h float64  // 默认浮点字面量类型

// 复数
var i complex64
var j complex128

// 布尔
var k bool // 布尔类型只有两个值：true 和 false，注意：Go 中不能将整数直接当作布尔值使用。

// 字符串（Go 中字符串本质上是**只读字节序列**）
var s string = "hello"

// rune 与 byte，字节与字符
var r rune = '中'  // int32 别名，Unicode 码点
var by byte = 'A'  // uint8 别名

```

**常见转义字符：**

 - `\n` 换行
 - `\t` 制表符
 - `\"` 双引号
 - `\\` 反斜杠
 - `\r` 回车
 - `\b` 退格
 - `\0` 空字符

**原始字符串：**

使用反引号 `` 定义原始字符串，保留原始格式：

```go
str := `这是一个原始字符串
不会处理转义字符 \n
可以换行`
```

### 3.4 类型转换

Go **没有隐式类型转换**，必须显式转换：

```go
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)

var x float64 = 3.14
var y int = int(x) // 注意：浮点转整数会截断小数部分，y 的值为 3


// 字符串 <-> 数字 需要用 strconv
s := strconv.Itoa(42)        // int -> string: "42"
n, _ := strconv.Atoi("42")   // string -> int: 42

// 字符串 <-> 字节切片
bs := []byte("hello")
s2 := string(bs)
```

### 3.5 运算符

没什么好说的，懂的都懂：

```go
// 算术：+  -  *  /  %  ++  --
// 位运算：&  |  ^  <<  >>  &^(位清除)
// 比较：==  !=  <  >  <=  >=
// 逻辑：&&  ||  !
// 取地址/解引用：&  *
// 通道：<-

// 注意：Go 没有三元运算符 (? :)
// 没有 ++ / -- 表达式（i++ 是语句，不能用于赋值）
i := 0
i++ // ✅
// j := i++ // ❌ 编译错误
```

### 3.6 输入输出

```go
// 输出
fmt.Print("no newline")
fmt.Println("with newline")
fmt.Printf("name: %s, age: %d\n", "Go", 14)

// 格式化浮点输出
fmt.Printf("pi: %.2f\n", 3.14159) // pi: 3.14

// 格式化动词
// %v  默认格式      %+v 带字段名（结构体）   %#v Go 语法格式
// %T  类型          %d  十进制整数           %f  浮点
// %s  字符串        %q  带引号字符串         %p  指针
// %b  二进制        %x  十六进制

// Sprintf 返回格式化字符串
s := fmt.Sprintf("score: %d", 100)

// 输入
var name string
var age int
fmt.Scan(&name) // 以空白字符（空格、制表符、换行）为分隔
fmt.Scanln(&name) // 以换行符为分隔
fmt.Scanf("%s", &name)// 按照格式读取输入，类似于 C 语言的 scanf
fmt.Scanf("%s %d", &name, &age)
```
>注意：输入函数需要传入变量地址，因此要使用 `&`。
---

## 4. 流程控制

### 4.1 if / else

```go
// 条件不需要小括号，但大括号必须
if x > 0 {
    fmt.Println("positive")
} else if x == 0 {
    fmt.Println("zero")
} else {
    fmt.Println("negative")
}

// if 可以带初始化语句（变量作用域仅限 if 块）
if err := doSomething(); err != nil {
    fmt.Println(err)
}
```

### 4.2 for

Go **只有 `for`**，没有 `while` 和 `do-while`。

```go
// 标准 for
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// 类似 while
n := 0
for n < 5 {
    n++
}

// 无限循环
for {
    break // 用 break 退出
}

// for range 遍历
nums := []int{10, 20, 30}
for index, value := range nums {
    fmt.Println(index, value)
}

// 只要 index
for i := range nums { ... }

// 只要 value
for _, v := range nums { ... }

// 遍历字符串（按 rune 遍历）
for i, ch := range "你好Go" {
    fmt.Printf("%d: %c\n", i, ch)
}

// 遍历 map
m := map[string]int{"a": 1, "b": 2}
for key, value := range m {
    fmt.Println(key, value)
}

// 遍历 channel
for msg := range ch {
    fmt.Println(msg)
}
```

### 4.3 switch

Go 的 `switch` 默认自带 `break`，不需要手动写。

```go
// 基本 switch（自动 break，不需要手动写）
switch day {
case "Mon":
    fmt.Println("Monday")
case "Tue", "Wed": // 多个值
    fmt.Println("Tue or Wed")
default:
    fmt.Println("other")
}

switch day {
case 1, 2, 3, 4, 5:
	fmt.Println("工作日")
case 6, 7:
	fmt.Println("周末")
}

// 无条件 switch（替代 if-else 链）
switch {
case score >= 90:
    fmt.Println("A")
case score >= 60:
    fmt.Println("B")
default:
    fmt.Println("C")
}

// fallthrough：强制穿透到下一个 case,只会执行下一个 case，不会继续往下穿透
switch n := 3; n {
case 3:
    fmt.Println("three")
    fallthrough
case 4:
    fmt.Println("four") // 也会执行
}

// 类型 switch
switch v := i.(type) {
case int:
    fmt.Println("int:", v)
case string:
    fmt.Println("string:", v)
default:
    fmt.Println("unknown")
}
```

### 4.4 break / continue / goto

```go
// 跳出当前循环。
for i := 0; i < 10; i++ {
    if i == 5 {
        break
    }
    fmt.Println(i)
}

// 跳过本次循环。
for i := 0; i < 5; i++ {
    if i == 2 {
        continue
    }
    fmt.Println(i)
}

// break / continue 配合标签
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                break outer // 跳出外层循环
            }
        }
    }

outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                continue outer // 跳过外层循环的当前迭代，进入下一次迭代
            }
        }
    }

// goto（很少使用）
goto End
fmt.Println("skipped")
End:
    fmt.Println("reached")
```

---

## 5. 函数