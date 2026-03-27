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
- [5. 复合数据类型](#5-复合数据类型)
- [6. 指针](#6-指针)
- [7. 函数](#7-函数)
- [8. 结构体](#8-结构体)
- [9. 方法](#9-方法)
- [10. 接口](#10-接口)
- [11. 错误处理](#11-错误处理)
- [12. 并发编程](#12-并发编程)
- [13. 泛型（Go 1.18+）](#13-泛型go-118)
- [14. 文件与 I/O 操作](#14-文件与-io-操作)
- [15. 反射（Reflection）](#15-反射reflection)
- [16. 常用标准库概览](#16-常用标准库概览)
- [17. 测试](#17-测试)
- [18. 包、模块与工程化管理](#18-包模块与工程化管理)
- [19. Go 常见易错点总结](#19-go-常见易错点总结)
- [20. 编码规范与最佳实践](#20-编码规范与最佳实践)

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

### 3.5 自定义类型与类型别名

```go
// 自定义类型（全新类型，可以添加方法）
type Celsius float64
type Handler func(string) error

type MyInt int
var a MyInt = 10
// MyInt 是 int 的新类型，不能直接赋值
var b int = a // ❌ 编译错误
var c int = int(a) // ✅ 需要显式转换

// 类型别名（完全等同于原类型）
type Byte = uint8
type Rune = int32
var x Byte = 255
var y uint8 = x // ✅ Byte 是 uint8 的别名，可以直接赋值
var z Rune = 'A' 
z = 300 // ✅ Rune 是 int32 的别名，可以直接赋值
```

### 3.6 运算符

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

### 3.7 输入输出

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

## 5. 复合数据类型

### 5.1 数组

数组长度固定，是**值类型**（赋值/传参会拷贝）。

```go
var arr [3]int                // [0, 0, 0]
arr2 := [3]int{1, 2, 3}
arr3 := [...]int{1, 2, 3}    // 编译器推断长度
arr4 := [5]int{0: 10, 4: 50} // 指定索引初始化

len(arr2) // 3

// 访问和修改数组元素
fmt.Println(arr[0])
arr[1] = 99
fmt.Println(arr)

// 遍历数组
for i := 0; i < len(arr); i++ {
    fmt.Println(arr[i])
}
// for range
for index, value := range arr {
    fmt.Println(index, value)
}

// 多维数组
matrix := [2][3]int{
    {1, 2, 3},
    {4, 5, 6},
}
```

Go 中数组是值类型，赋值会复制整个数组。

```go
a := [3]int{1, 2, 3}
b := a
b[0] = 100

fmt.Println(a) // [1 2 3]
fmt.Println(b) // [100 2 3]
```

这也是为什么在 Go 中更常使用切片而不是数组。

### 5.2 切片 (Slice)

切片是**引用类型**，底层指向一个数组。Go 中实际使用切片远多于数组。

```go
// 声明
var s []int              // nil 切片
s2 := []int{1, 2, 3}     // 字面量
s3 := make([]int, 5)     // len=5, cap=5，长度为 5，容量为 5
s4 := make([]int, 3, 10) // len=3, cap=10，长度为 3，容量为 10

// 从数组/切片截取（左闭右开），包含起始索引，不包含结束索引
arr := [5]int{1, 2, 3, 4, 5}
s5 := arr[1:4]  // [2, 3, 4]
s6 := arr[:3]   // [1, 2, 3]
s7 := arr[2:]   // [3, 4, 5]
s8 := arr[:]    // [1, 2, 3, 4, 5]
/

// 切片也可以从切片截取
s9 := []int{1, 2, 3, 4, 5}
s10 := s9[1:3]       // [2, 3]
fmt.Println(s10)     // [2 3]
fmt.Println(s9[:3])  // [1 2 3]
fmt.Println(s9[2:])  // [3 4 5]

// append：追加元素（底层数组满了会触发扩容）
// 注意：`append` 触发扩容后会返回新的底层数组，因此必须接收返回值。
/*
  原理：切片底层存储是数组，切片的容量即底层数组的长度，数组的长度是不可变的，因此如果切片容量足够，直接在原数组上追加；如果容量不足，`append` 会分配一个新的更大的底层数组，复制原有数据，并追加新元素。
*/
s = append(s, 1, 2, 3)
s = append(s, []int{4, 5}...) // 展开追加

// copy
src := []int{1, 2, 3}
dst := make([]int, len(src))
copy(dst, src)

// 删除元素（无内置方法，用切片技巧）
s = append(s[:i], s[i+1:]...) // 删除索引 i

// len 和 cap，长度和容量
fmt.Println(len(s), cap(s))

// 切片判空用 len，不要用 nil 判断
if len(s) == 0 { ... }
```

**切片底层结构**：`{ pointer *array, len int, cap int }`

### 5.3 映射 (Map)

```go
// 声明与初始化
var m map[string]int              // nil map，不能写入
m1 := map[string]int{"a": 1, "b": 2}
m2 := make(map[string]int)       // 空 map，可以写入

// 增 / 改
m2["key"] = 100

// 查
value := m1["a"]
value, ok := m1["a"] // ok 判断 key 是否存在
if v, ok := m1["c"]; !ok {
    fmt.Println("key not found")
}

// 删
delete(m1, "a")

// 遍历（无序）
for k, v := range m1 {
    fmt.Println(k, v)
}

// map 长度
len(m1)
```

> Map 是**引用类型**，不是并发安全的。并发场景用 `sync.Map` 或加锁。

### 5.4 结构体 (Struct)

结构体后面单独讲，这里简单介绍。

```go

type User struct {
    Name string
    Age  int
}
u := User{Name: "Alice", Age: 30}
// 访问字段
fmt.Println(u.Name)
// 修改字段
u.Age = 31
```

## 6. 指针


Go 有指针，但不像 C/C++ 那样复杂，**没有指针运算**（不能做 `p++`）。

指针保存的是变量的内存地址。

```go
var p *int       // 声明，零值为 nil
x := 42
p = &x           // 取 x 的地址
fmt.Println(*p)  // 解引用，根据地址获取值：42
*p = 100         // 通过指针修改值，x 的值变成 100
fmt.Println(x)   // 100

// new：分配内存，并返回对应类型的指针
p2 := new(int)   // *int，值为 0

// 函数传指针以修改原值，当需要在函数内部修改外部变量时，常使用指针参数
func increment(n *int) {
    *n++
}
val := 10
increment(&val)
fmt.Println(val) // 11
```


> Go 中 slice、map、channel、function、interface 本身就是引用语义，传参时无需取指针。

## 7. 函数

### 7.1 函数定义

```go
func add(a int, b int) int {
    return a + b
}

// 参数类型相同可合并
func add(a, b int) int {
    return a + b
}
```

### 7.2 多返回值

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func calc(a, b int) (int, int) {
	return a + b, a - b
}

result, err := divide(10, 3)
sum, diff := calc(10, 3)
```

### 7.3 命名返回值

```go
func divide(a, b float64) (result float64, err error) {
    if b == 0 {
        err = errors.New("division by zero")
        return // 裸 return，返回命名变量的当前值
    }
    result = a / b
    return
}

func calc(a, b int) (sum int, diff int) {
	sum = a + b
	diff = a - b
	return
}
```

### 7.4 可变参数

```go
func sum(nums ...int) int { // nums 是一个 []int 切片
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3)

// 切片展开传递
nums := []int{1, 2, 3}
sum(nums...)
```

### 7.5 函数作为值 / 参数

```go
// 函数是一等公民，可以赋值给变量，作为参数传递，甚至作为返回值。
func add(a, b int) int {
    return a + b
}

var fn func(int, int) int = add
fmt.Println(fn(1, 2))

// 函数作为参数
func apply(a, b int, op func(int, int) int) int {
    return op(a, b)
}
apply(3, 4, add)
```

### 7.6 匿名函数

```go
// 直接调用
result := func(a, b int) int {
    return a + b
}(3, 4)

// 赋值给变量
double := func(n int) int {
    return n * 2
}
fmt.Println(double(5))
```

### 7.7 闭包

闭包是**函数与其引用的外部变量的组合**，可以访问和修改外部变量。

```go
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

c := counter()
fmt.Println(c()) // 1
fmt.Println(c()) // 2
fmt.Println(c()) // 3
```

### 7.8 init 函数

每个包可以有一个或多个 `init()` 函数，在 `main()` 之前自动执行，常用于初始化。

```go
func init() {
    fmt.Println("初始化...")
}
```

执行顺序：**导入包的 init() → 当前包的 init() → main()**

### 7.9 defer 语句

`defer` 用于注册一个函数，在当前函数返回时执行，常用于资源释放，**后进先出 (LIFO)**。

```go
func main() {
    fmt.Println("start")
    defer fmt.Println("deferred 1")
    defer fmt.Println("deferred 2")
    fmt.Println("end")
}
// 输出：start → end → deferred 2 → deferred 1

// 经典用法：关闭资源
func readFile(path string) {
    f, err := os.Open(path)
    if err != nil {
        return
    }
    defer f.Close() // 函数结束时一定会关闭

    // 读取文件...
}
```

> **注意**：`defer` 的参数在声明时就已求值，而非执行时。

```go
x := 10
defer fmt.Println(x) // 打印 10，不是 20
x = 20
```

---

## 8. 结构体

结构体用于将多个不同类型的数据组合成一个整体，注意：结构体是值类型。

```go
// 定义
type User struct {
    Name string
    Age  int
    Email string
}

// 初始化
u1 := User{Name: "Bob", Age: 30}           // 按字段名（推荐）
u2 := User{"Alice", 25, "alice@go.dev"}    // 按顺序
u3 := new(User)                            // 返回 *User，字段为零值
var u4 User                                // 零值结构体

// 访问与修改
u2.Email = "bob@go.dev"

// 指针访问（自动解引用）
p := &u2
p.Name = "Bob2" // 等同于 (*p).Name = "Bob2"

// 结构体指针

p := &User{
    Name: "Alice",
    Age:  20,
}
fmt.Println(p.Name) // Go 对结构体指针访问字段时可以直接用 `p.Name`，不需要写成 `(*p).Name`，底层会自动解引用成 `(*p).Name`。

// 匿名字段（类型为结构体/组合）
type Admin struct {
    User          // 嵌入 User,相当于 Admin 继承了 User 的字段和方法
    Level int
}
a := Admin{User: User{Name: "Root", Age: 0}, Level: 1}
a.Name // 直接访问 User 的字段。

// 匿名字段（类型为内置类型），不过这种写法在工程实践中不够直观，通常只在特定场景使用。
type User struct {
	string
	int
}

// 结构体嵌套
type Address struct {
    City string
}

type User struct {
    Name    string
    Address Address
}

u := User{
    Name: "Tom",
    Address: Address{
        City: "Shanghai",
    },
}
fmt.Println(u.Address.City)

// 匿名结构体
point := struct {
    X, Y int
}{10, 20}

// 结构体标签（用于 JSON 序列化等）
type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price,omitempty"`
}

```

## 9. 方法

方法本质上是绑定到某种类型上的函数，这个类型也被称为方法的**接收者 (receiver)**。

```go
// 示例1：
type Rectangle struct {
    Width, Height float64
}

// 值接收者：不修改原对象
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// 指针接收者：可修改原对象
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

rect := Rectangle{10, 5}
fmt.Println(rect.Area()) // 50
rect.Scale(2)
fmt.Println(rect.Area()) // 200

// 示例2：
type Person struct {
	Name string
}

// 值接收者，修改不会影响原对象
func (p Person) SetName1(name string) {
    // 这里修改的 p 是副本，不会影响原对象
    p.Name = name
}

func (p *Person) SetName2(name string) {
    // 这里修改的 p 是指针，可以修改原对象
    p.Name = name
}
```

**值接收者**
- 接收的是副本
- 修改不会影响原对象
- 适合小对象、只读场景

**指针接收者**
- 接收的是对象地址
- 修改会影响原对象
- 避免大对象拷贝
- 更适合结构体方法

在实际开发中，结构体方法通常更常使用指针接收者。

## 10 接口

接口是 Go 实现抽象和多态的核心机制，定义了一组方法签名，Go语言通过隐式接口实现，无需显式声明实现接口，只要一个类型包含接口定义的所有方法，该类型就自动实现了该接口，从而实现了灵活的结构化多态（duck typing-鸭子类型）。

### 10.1 定义接口与实现接口
```go
// 定义接口
type Speaker interface {
	Speak()
}


// 定义结构体
type Dog struct{}
// 实现接口
func (d Dog) Speak() {
	fmt.Println("汪汪汪")
}
var s Speaker
s = Dog{}
s.Speak() // 输出：汪汪汪
```
### 10.2 多态

接口类型的变量可以指向任何实现了该接口的类型的实例，从而实现多态。

```go
type Cat struct{}

func (c Cat) Speak() {
	fmt.Println("喵喵喵")
}

func makeSound(s Speaker) {
	s.Speak()
}

makeSound(Dog{}) // 输出：汪汪汪
makeSound(Cat{}) // 输出：喵喵喵
```

### 10.3 空接口

空接口 `interface{}` (Go 1.18+ 可用 `any`)是一个特殊的接口类型，没有任何方法，因此所有类型都实现了空接口，可以用来存储任意类型的值。

在 Go 1.18 之前，空接口常用于“通用类型”场景。现在很多场景也可以考虑泛型。

```go
var i interface{} = "hello"
i = 42
i = true

// 常见用途：通用函数
func PrintAnything(v any) {
    fmt.Println(v)
}
```

### 10.4 类型断言

类型断言用于从接口类型中提取具体类型的值，语法为 `x.(T)`，其中 `x` 是接口类型的变量，`T` 是要断言的具体类型。

```go
var i interface{} = "hello"

// 类型断言
s := i.(string)        // 成功：s = "hello"
// n := i.(int)        // panic!

// 安全断言
s, ok := i.(string)    // ok = true, s = "hello"
n, ok := i.(int)       // ok = false, n = 0

// 类型 switch
switch v := i.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
default:
    fmt.Println("unknown")
}
```
### 10.5 接口组合

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}
```

### 10.6 常用内置接口

```go
// fmt.Stringer —— 类似 Java 的 toString()
type Stringer interface {
    String() string
}

func (u User) String() string {
    return fmt.Sprintf("%s (%d)", u.Name, u.Age)
}

// error 接口
type error interface {
    Error() string
}

// sort.Interface
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

## 11. 错误处理

Go 没有异常机制，**没有 try-catch**，错误通过返回值传递，通常返回一个 `error` 类型的值来表示是否发生错误。

### 11.1 error 接口

```go
// Go 内置 `error` 接口：
type error interface {
    Error() string
}

// 返回错误
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("除数不能为 0")
    }
    return a / b, nil
}

// 调用函数并处理错误
result, err := divide(10, 0)
if err != nil {
    fmt.Println("错误:", err)
    return
}
fmt.Println(result)

// errors.New 创建简单错误，导入包 import "errors"
err := errors.New("这是一个错误")
// divide可改写为：
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("除数不能为 0")
    }
    return a / b, nil
}

// fmt.Errorf 支持格式化错误信息
err := fmt.Errorf("无法打开文件: %s", filename)

// 并且可以使用 %w 包装原错误，便于错误链追踪
err := fmt.Errorf("读取文件失败: %w", originalErr)

// Go 1.13+ 引入了 errors.Is 和 errors.As 来检查错误链
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("文件不存在")
}

var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("路径错误:", pathErr.Path)
}

// 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("field %s: %s", e.Field, e.Message)
}

func validate(name string) error {
    if name == "" {
        return &ValidationError{Field: "name", Message: "cannot be empty"}
    }
    return nil
}
```
### 11.2 panic 和 recover

`panic`（翻译为“恐慌”） 表示程序发生严重错误，通常会中断程序执行，类似 C# 的异常；`recover` 用于在 `defer` 中捕获 panic。

```go
// 触发 panic
func safeDivide(a, b int) {
    defer func() {
        if r := recover(); r != nil { // `recover` 可以捕获 `panic`，防止程序直接崩溃，但必须配合 `defer` 使用。
            fmt.Println("Recovered:", r)
        }
    }()
    fmt.Println(a / b) // b=0 时 panic，被 recover 捕获，程序不会崩溃，类似于 C# 的 try-catch。
}

safeDivide(10, 0) // Recovered: runtime error: integer divide by zero

// 主动 panic
panic("发生严重错误")
```

> **惯例**：正常错误用 `error` 返回值，`panic` 仅用于程序级严重错误（如初始化失败、断言某些绝不应该发生的情况）。

---

## 12. 并发编程

Go 的并发编程模型基于**goroutine**和**channel**。

goroutine 是 Go 的轻量级线程，并非操作系统线程，也称为协程（初始栈仅 ~2KB），使用 `go` 关键字创建。channel 是 goroutine 之间进行通信和同步的管道。

### 12.1 goroutine

```go
func task() {
    fmt.Println("执行任务")
}

func main() {
    go task() // 启动一个新的 goroutine 执行 task 函数，类似于 C# 的 Task.Run(() => task())，但更轻量级。

    go func() { // 启动一个新的 goroutine 执行匿名函数
        fmt.Println("anonymous goroutine")
    }()

    fmt.Println("main")
}
// 注意：main 退出时所有 goroutine 会被终止
// 需要等待机制（WaitGroup / channel）
```

### 12.2 channel

channel 是 goroutine 之间通信的管道。

```go
// 无缓冲 channel（同步：发送和接收会互相等待）
ch := make(chan int)

go func() {
    ch <- 42 // 发送
}()
value := <-ch // 接收
fmt.Println(value) // 42

// 有缓冲 channel（异步：缓冲区满才阻塞）
ch2 := make(chan string, 3)
ch2 <- "a"
ch2 <- "b"
fmt.Println(<-ch2) // "a"

// 关闭 channel，关闭后不能继续发送数据，但仍可以继续接收已经存在的数据。
close(ch2)

// range 遍历 channel（直到 channel 关闭）
go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}()
for v := range ch {
    fmt.Println(v)
}

// 单向 channel（用于函数签名约束）
func producer(out chan<- int) { out <- 1 }  // 只写
func consumer(in <-chan int)  { <-in }      // 只读
```

完整 channel 小示例：

```go
func main() {
    ch := make(chan int)

    // 生产者
    go func() {
        for i := 0; i < 5; i++ {
            ch <- i
            fmt.Println("Produced:", i)
        }
        close(ch) // 生产完成后关闭 channel
    }()

    // 消费者
    for v := range ch { // 从 channel 接收数据，直到 channel 关闭
        fmt.Println("Consumed:", v)
    }
}
```
`select` 用于同时监听多个 channel。

```go
ch1 := make(chan string)
ch2 := make(chan string)

go func() { time.Sleep(1 * time.Second); ch1 <- "one" }()
go func() { time.Sleep(2 * time.Second); ch2 <- "two" }()

// 等待多个 channel 的消息，哪个先到就处理哪个，如果 3 秒内都没有消息到达，则执行 timeout 分支。
select {
case msg := <-ch1:
    fmt.Println(msg)
case msg := <-ch2:
    fmt.Println(msg)
case <-time.After(3 * time.Second):
    fmt.Println("timeout")
}

// 非阻塞操作，如果 ch1 没有消息可接收，则直接执行 default 分支，而不会阻塞等待。
select {
case msg := <-ch1:
    fmt.Println(msg)
default:
    fmt.Println("no message")
}
```

### 12.3 sync 包

Go 的 `sync` 包提供了多种同步原语，如 `WaitGroup`、`Mutex`、`RWMutex` 等，用于更复杂的并发控制。

**sync.WaitGroup**

`WaitGroup` 用于等待一组 goroutine 执行完成。

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Println("worker", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()
    fmt.Println("所有任务完成")
}
```

**sync 包其他常用工具**

```go
// Mutex 互斥锁
var mu sync.Mutex
var count int

func increment() {
    mu.Lock()
    defer mu.Unlock()
    count++
}

// RWMutex 读写锁（多读单写）
var rwmu sync.RWMutex
rwmu.RLock()    // 读锁
rwmu.RUnlock()
rwmu.Lock()     // 写锁
rwmu.Unlock()

// Once 只执行一次
var once sync.Once
once.Do(func() {
    fmt.Println("只执行一次")
})

// sync.Map 并发安全的 map
var sm sync.Map
sm.Store("key", "value")
v, ok := sm.Load("key")
sm.Delete("key")
sm.Range(func(key, value any) bool {
    fmt.Println(key, value)
    return true // 返回 false 停止遍历
})
```

### 12.4 context 包

`context` 包用于在 goroutine 之间传递取消信号和请求范围的数据，常用于控制 API 请求、数据库操作等需要超时或取消的场景。

```go
// 创建一个带有取消功能的 context
ctx, cancel := context.WithCancel(context.Background())
cancel() // 取消 context

// 示例：

// 使用 context 控制 goroutine 生命周期
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // 确保在 main 函数结束时取消 context

go func(ctx context.Context) { // 启动一个 goroutine，监听 context 的取消信号
    for {
        select {
        case <-ctx.Done(): // 监听取消信号
            fmt.Println("cancelled:", ctx.Err())
            fmt.Println("goroutine 退出")
            return
        default:
            fmt.Println("goroutine 工作中...")
            // 模拟工作
            time.Sleep(500 * time.Millisecond)
        }
    }
}(ctx)
```

## 13. 泛型（Go 1.18+）

Go 1.18 引入了泛型（Generics），允许函数和类型在定义时不指定具体类型，而在使用时再指定，从而实现代码的复用和类型安全。

### 13.1 泛型函数

```go
// 类型参数用方括号
func Max[T int | float64 | string](a, b T) T {
    if a > b {
        return a
    }
    return b
}

Max(3, 5)         // 自动推断 T = int
Max[string]("a", "b") // 显式指定
```

### 13.2 类型约束

```go
// 定义约束接口
type Number interface {
    int | int8 | int16 | int32 | int64 |
    float32 | float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

// ~int 表示底层类型为 int 的所有类型
type Integer interface {
    ~int | ~int32 | ~int64
}

type MyInt int // MyInt 底层是 int，满足 ~int

// comparable 内置约束：支持 == 和 != 的类型
func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {
            return true
        }
    }
    return false
}

// any = interface{}
func Print[T any](v T) {
    fmt.Println(v)
}
```

### 13.3 泛型结构体

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// 使用
s := Stack[int]{}
s.Push(1)
s.Push(2)
v, _ := s.Pop() // 2
```
---
## 14. 文件与 I/O 操作

Go 的 I/O 体系围绕两个核心接口构建：`io.Reader` 和 `io.Writer`。理解这一点后，所有 I/O 操作都是一脉相通的。

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

> 文件、网络连接、压缩流、加密流……都实现了这两个接口，可以自由组合。

---

### 14.1 打开与关闭文件

```go
// os.Open —— 只读打开
file, err := os.Open("data.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// os.Create —— 创建/截断文件（可写）
file, err := os.Create("output.txt")

// os.OpenFile —— 完全控制模式和权限
file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
```

**常用标志位组合**：

| 场景                 | 标志位                                      |
| -------------------- | ------------------------------------------- |
| 只读                 | `os.O_RDONLY`                               |
| 只写（覆盖）         | `os.O_WRONLY \| os.O_CREATE \| os.O_TRUNC`  |
| 追加写入             | `os.O_APPEND \| os.O_CREATE \| os.O_WRONLY` |
| 读写                 | `os.O_RDWR \| os.O_CREATE`                  |
| 文件必须不存在才创建 | `os.O_CREATE \| os.O_EXCL \| os.O_WRONLY`   |

> ⚠️ **必须用 `defer file.Close()`**，否则文件描述符泄漏。在循环中打开文件时尤其注意，不要在循环里只 defer（要封装成函数或手动 close）。

---

### 14.2 读取文件

#### 一次性读取整个文件（小文件）

```go
// 最简方式：os.ReadFile（Go 1.16+）
data, err := os.ReadFile("config.json")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(data))

// 等价的手动方式
file, _ := os.Open("config.json")
defer file.Close()
data, _ := io.ReadAll(file)
```

> 适合小文件（< 几十 MB）。整个内容加载到内存，大文件会 OOM。

#### 按行读取（最常用）

```go
func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}
```

> ⚠️ `bufio.Scanner` 默认最大行长为 **64KB**。超长行需要调大缓冲：

```go
scanner := bufio.NewScanner(file)
scanner.Buffer(make([]byte, 0), 10*1024*1024) // 最大 10MB 一行
```

#### 按固定大小块读取（大文件）

```go
func readInChunks(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    buf := make([]byte, 32*1024) // 32KB 缓冲区
    for {
        n, err := file.Read(buf)
        if n > 0 {
            // 处理 buf[:n]
            fmt.Printf("read %d bytes\n", n)
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
    }
    return nil
}
```

#### 使用 bufio.Reader 灵活读取

```go
file, _ := os.Open("data.txt")
defer file.Close()
reader := bufio.NewReader(file)

// 按分隔符读取
for {
    // ReadString 读到分隔符为止（包含分隔符）
    line, err := reader.ReadString('\n')
    if len(line) > 0 {
        line = strings.TrimRight(line, "\r\n")
        fmt.Println(line)
    }
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
}

// 按字节读取
b, err := reader.ReadByte()

// 窥探但不消费
peeked, err := reader.Peek(5) // 查看前 5 个字节
```

#### 从指定位置读取（随机读取）

```go
file, _ := os.Open("data.bin")
defer file.Close()

// Seek 到指定位置
file.Seek(100, io.SeekStart)   // 从头偏移 100 字节
file.Seek(-50, io.SeekEnd)     // 从末尾倒退 50 字节
file.Seek(10, io.SeekCurrent)  // 从当前位置前进 10 字节

buf := make([]byte, 64)
n, _ := file.Read(buf)

// ReadAt 不改变文件偏移量（并发安全）
n, err := file.ReadAt(buf, 1024) // 从偏移 1024 处读取
```

#### io.SectionReader —— 读取文件的一个片段

```go
file, _ := os.Open("large.bin")
defer file.Close()

// 从偏移 1000 处读取 500 字节
section := io.NewSectionReader(file, 1000, 500)
data, _ := io.ReadAll(section)
```

---

### 14.3 写入文件

#### 一次性写入（小内容）

```go
// os.WriteFile（Go 1.16+）
err := os.WriteFile("output.txt", []byte("Hello, Go!\n"), 0644)
if err != nil {
    log.Fatal(err)
}
```

#### 使用 bufio.Writer 高效写入

裸 `file.Write()` 每次都触发系统调用。`bufio.Writer` 在内存中累积数据，批量写入。

```go
func writeLines(path string, lines []string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for _, line := range lines {
        // WriteString 比 Write([]byte(...)) 更高效
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            return err
        }
    }
    return writer.Flush() // 关键！将缓冲区剩余数据刷到磁盘
}
```

> ⚠️ **忘记 `Flush()` 是常见 Bug**。数据还在缓冲区里，文件可能为空或不完整。

#### 使用 fmt.Fprintf 格式化写入

```go
file, _ := os.Create("report.txt")
defer file.Close()

writer := bufio.NewWriter(file)
defer writer.Flush()

fmt.Fprintf(writer, "User: %s\n", "Alice")
fmt.Fprintf(writer, "Score: %d\n", 95)
fmt.Fprintf(writer, "Rate: %.2f%%\n", 87.5)
```

#### 追加写入

```go
func appendToFile(path, content string) error {
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(content)
    return err
}

appendToFile("app.log", "[2024-03-15 10:30:00] server started\n")
```

#### 在指定位置写入

```go
file, _ := os.OpenFile("data.bin", os.O_RDWR, 0644)
defer file.Close()

file.Seek(100, io.SeekStart)
file.Write([]byte{0xFF, 0xFE})

// WriteAt 不改变偏移量
file.WriteAt([]byte("hello"), 200)
```

---

### 14.4 大文件处理实战

#### 示例1：大文件复制（流式，内存恒定）

```go
func copyFile(src, dst string) (int64, error) {
    sourceFile, err := os.Open(src)
    if err != nil {
        return 0, err
    }
    defer sourceFile.Close()

    destFile, err := os.Create(dst)
    if err != nil {
        return 0, err
    }
    defer destFile.Close()

    // io.Copy 内部使用 32KB 缓冲区
    // 无论文件多大，内存占用恒定
    nBytes, err := io.Copy(destFile, sourceFile)
    if err != nil {
        return 0, err
    }

    // 确保数据落盘
    err = destFile.Sync()
    return nBytes, err
}

// 指定缓冲区大小
func copyFileBuffered(src, dst string, bufSize int) (int64, error) {
    sourceFile, _ := os.Open(src)
    defer sourceFile.Close()
    destFile, _ := os.Create(dst)
    defer destFile.Close()

    buf := make([]byte, bufSize)
    return io.CopyBuffer(destFile, sourceFile, buf)
}

// 使用：1MB 缓冲区复制大文件
copyFileBuffered("10GB.dat", "backup.dat", 1024*1024)
```

#### 示例2：大文件逐行处理（如日志分析）

处理一个 10GB 的 Nginx 日志文件，统计每个 IP 的请求次数：

```go
func analyzeLog(path string) (map[string]int, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    ipCount := make(map[string]int)
    scanner := bufio.NewScanner(file)
    // Nginx 日志一行通常不超过 4KB，默认 64KB 足够
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := scanner.Text()

        // 提取第一个字段（IP 地址）
        // 日志格式：192.168.1.1 - - [15/Mar/2024:10:30:00 +0800] "GET / HTTP/1.1" 200 612
        if idx := strings.IndexByte(line, ' '); idx > 0 {
            ip := line[:idx]
            ipCount[ip]++
        }

        // 每百万行打印一次进度
        if lineNum%1_000_000 == 0 {
            fmt.Printf("processed %d million lines...\n", lineNum/1_000_000)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("scan error at line %d: %w", lineNum, err)
    }

    return ipCount, nil
}
```

#### 示例3：大文件分块写入（生成测试数据）

```go
func generateLargeFile(path string, sizeGB int) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriterSize(file, 4*1024*1024) // 4MB 缓冲
    defer writer.Flush()

    targetBytes := int64(sizeGB) * 1024 * 1024 * 1024
    var written int64

    for written < targetBytes {
        line := fmt.Sprintf("%s [INFO] request_id=%d user_agent=Mozilla/5.0 path=/api/v1/data\n",
            time.Now().Format(time.RFC3339), rand.Intn(1_000_000))
        n, err := writer.WriteString(line)
        if err != nil {
            return err
        }
        written += int64(n)
    }

    fmt.Printf("wrote %d bytes (%.2f GB)\n", written, float64(written)/(1024*1024*1024))
    return nil
}
```

#### 示例4：带进度条的大文件复制

```go
// 自定义 Reader 包装器，跟踪读取进度
type ProgressReader struct {
    reader     io.Reader
    total      int64
    current    int64
    lastReport float64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
    n, err := pr.reader.Read(p)
    pr.current += int64(n)

    if pr.total > 0 {
        percent := float64(pr.current) / float64(pr.total) * 100
        // 每增长 1% 才打印，避免刷屏
        if percent-pr.lastReport >= 1 {
            fmt.Printf("\r  progress: %.1f%% (%s / %s)",
                percent, formatBytes(pr.current), formatBytes(pr.total))
            pr.lastReport = percent
        }
    }
    return n, err
}

func formatBytes(b int64) string {
    const unit = 1024
    if b < unit {
        return fmt.Sprintf("%d B", b)
    }
    div, exp := int64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func copyWithProgress(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    info, _ := sourceFile.Stat()

    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()

    pr := &ProgressReader{
        reader: sourceFile,
        total:  info.Size(),
    }

    _, err = io.Copy(destFile, pr)
    fmt.Println() // 换行
    return err
}
```

---

### 14.5 临时文件与临时目录

```go
// 创建临时文件
tmpFile, err := os.CreateTemp("", "myapp-*.txt")
// 第一个参数为空表示使用系统默认临时目录
// 第二个参数中 * 会被替换为随机字符串
if err != nil {
    log.Fatal(err)
}
defer os.Remove(tmpFile.Name()) // 用完删除
defer tmpFile.Close()

fmt.Println(tmpFile.Name()) // 如 /tmp/myapp-384729481.txt
tmpFile.WriteString("temporary data")

// 创建临时目录
tmpDir, err := os.MkdirTemp("", "myapp-*")
if err != nil {
    log.Fatal(err)
}
defer os.RemoveAll(tmpDir) // 递归删除

fmt.Println(tmpDir) // 如 /tmp/myapp-291837465
```

---

### 14.6 文件信息与权限

```go
// 获取文件信息
info, err := os.Stat("data.txt")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("file does not exist")
    }
    log.Fatal(err)
}

fmt.Println(info.Name())    // "data.txt"
fmt.Println(info.Size())    // 字节数
fmt.Println(info.Mode())    // "-rw-r--r--"
fmt.Println(info.ModTime()) // 修改时间
fmt.Println(info.IsDir())   // 是否是目录

// Lstat 不跟踪符号链接（Stat 会跟踪）
info, err = os.Lstat("symlink.txt")

// 修改权限
os.Chmod("script.sh", 0755)

// 修改所有者（需要 root）
os.Chown("data.txt", uid, gid)

// 修改时间戳
os.Chtimes("data.txt", time.Now(), time.Now())

// 判断文件类型
if info.Mode().IsRegular() { /* 普通文件 */ }
if info.Mode().IsDir()     { /* 目录 */ }
if info.Mode()&os.ModeSymlink != 0 { /* 符号链接 */ }
```

---

### 14.7 目录操作

```go
// 创建目录
os.Mkdir("logs", 0755)          // 单层
os.MkdirAll("a/b/c/d", 0755)   // 递归创建

// 读取目录内容
entries, err := os.ReadDir(".")
if err != nil {
    log.Fatal(err)
}
for _, entry := range entries {
    info, _ := entry.Info()
    fmt.Printf("%-30s %10d %s\n", entry.Name(), info.Size(), info.Mode())
}

// 递归遍历目录树
err := filepath.Walk("/var/log", func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    if info.IsDir() && info.Name() == ".git" {
        return filepath.SkipDir // 跳过整个目录
    }
    fmt.Println(path)
    return nil
})

// Go 1.16+ filepath.WalkDir 更高效（不调用 Stat）
err := filepath.WalkDir("/var/log", func(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }
    fmt.Println(path, d.IsDir())
    return nil
})

// 实际场景：查找所有 .go 文件
func findGoFiles(root string) ([]string, error) {
    var files []string
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() && filepath.Ext(path) == ".go" {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}

// Glob 模式匹配
matches, _ := filepath.Glob("/var/log/*.log")
for _, m := range matches {
    fmt.Println(m)
}
```

---

### 14.8 filepath 包详解

```go
import "path/filepath"

// 拼接路径（自动处理分隔符）
p := filepath.Join("home", "user", "docs", "file.txt")
// Linux: "home/user/docs/file.txt"
// Windows: "home\\user\\docs\\file.txt"

// 提取各部分
filepath.Base("/a/b/c.txt")   // "c.txt"
filepath.Dir("/a/b/c.txt")    // "/a/b"
filepath.Ext("/a/b/c.tar.gz") // ".gz"

// 去掉扩展名
name := strings.TrimSuffix("report.csv", filepath.Ext("report.csv")) // "report"

// 绝对路径
abs, _ := filepath.Abs("./data.txt") // "/home/user/project/data.txt"

// 相对路径
rel, _ := filepath.Rel("/home/user", "/home/user/project/main.go")
// "project/main.go"

// 清理路径
filepath.Clean("/a/b/../c/./d") // "/a/c/d"

// 分割路径和文件名
dir, file := filepath.Split("/home/user/data.txt")
// dir = "/home/user/", file = "data.txt"

// 匹配模式
matched, _ := filepath.Match("*.txt", "readme.txt") // true
matched, _ = filepath.Match("log_202?-*", "log_2024-03") // true
```

---

### 14.9 io 包核心工具

```go
import "io"

// io.Copy：流式复制（已在 12.4 演示）
io.Copy(dst, src)             // 使用默认 32KB 缓冲
io.CopyN(dst, src, 1024)     // 只复制 1024 字节
io.CopyBuffer(dst, src, buf) // 使用自定义缓冲

// io.ReadAll：读取全部内容
data, err := io.ReadAll(reader)

// io.ReadFull：确保读满 buf
buf := make([]byte, 100)
_, err := io.ReadFull(reader, buf) // 不够 100 字节会报错

// io.LimitReader：限制最大读取量
limited := io.LimitReader(file, 1024*1024) // 最多读 1MB
data, _ := io.ReadAll(limited)

// io.MultiReader：串联多个 Reader
header := strings.NewReader("HEADER\n")
body, _ := os.Open("body.txt")
footer := strings.NewReader("\nFOOTER")
combined := io.MultiReader(header, body, footer)
io.Copy(os.Stdout, combined) // 依次输出三部分

// io.MultiWriter：同时写入多个目标
logFile, _ := os.Create("app.log")
multi := io.MultiWriter(os.Stdout, logFile)
fmt.Fprintln(multi, "this goes to both console and file")

// io.TeeReader：读取的同时复制一份
var buf bytes.Buffer
tee := io.TeeReader(resp.Body, &buf)
io.ReadAll(tee)
// 此时 buf 中也有一份 resp.Body 的完整内容

// io.Pipe：同步的内存管道
pr, pw := io.Pipe()
go func() {
    defer pw.Close()
    json.NewEncoder(pw).Encode(data) // 写入端
}()
json.NewDecoder(pr).Decode(&result)  // 读取端

// io.NopCloser：给 Reader 加一个空的 Close 方法
rc := io.NopCloser(strings.NewReader("hello"))
// rc 实现了 io.ReadCloser
```

---

### 15.10 bufio 高级用法

```go
import "bufio"

// 自定义缓冲区大小
reader := bufio.NewReaderSize(file, 1*1024*1024) // 1MB 读缓冲
writer := bufio.NewWriterSize(file, 1*1024*1024) // 1MB 写缓冲

// bufio.Scanner 自定义分割函数
scanner := bufio.NewScanner(file)

// 内置分割函数
scanner.Split(bufio.ScanLines)  // 默认：按行
scanner.Split(bufio.ScanWords)  // 按单词
scanner.Split(bufio.ScanBytes)  // 按字节
scanner.Split(bufio.ScanRunes)  // 按 UTF-8 字符

// 自定义分割：按 \0 分割（处理 null-terminated 数据）
scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
    if atEOF && len(data) == 0 {
        return 0, nil, nil
    }
    if i := bytes.IndexByte(data, 0); i >= 0 {
        return i + 1, data[:i], nil
    }
    if atEOF {
        return len(data), data, nil
    }
    return 0, nil, nil // 需要更多数据
})

// ReadWriter：组合 Reader 和 Writer
rw := bufio.NewReadWriter(
    bufio.NewReader(conn),
    bufio.NewWriter(conn),
)
```

---

### 14.11 实战：CSV 文件处理

```go
import "encoding/csv"

// 读取 CSV
func readCSV(path string) ([][]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ','        // 分隔符，默认逗号
    reader.Comment = '#'      // 注释行前缀
    reader.TrimLeadingSpace = true

    // 一次读取全部
    records, err := reader.ReadAll()
    return records, err
}

// 逐行读取大 CSV（省内存）
func processLargeCSV(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)

    // 读取表头
    header, err := reader.Read()
    if err != nil {
        return err
    }
    fmt.Println("Columns:", header)

    lineNum := 1
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("line %d: %w", lineNum, err)
        }
        lineNum++

        // 处理每一行
        // record[0], record[1], ...
        _ = record
    }
    return nil
}

// 写入 CSV
func writeCSV(path string, data [][]string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    // 写入 UTF-8 BOM（让 Excel 正确识别中文）
    file.Write([]byte{0xEF, 0xBB, 0xBF})

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // 写表头
    writer.Write([]string{"姓名", "年龄", "邮箱"})

    // 写数据
    for _, row := range data {
        if err := writer.Write(row); err != nil {
            return err
        }
    }

    return writer.Error()
}
```

---

### 14.12 实战：JSON 文件流式处理

当 JSON 文件非常大（如包含百万条记录的数组），不适合 `json.Unmarshal` 全部加载：

```go
// 流式解析大 JSON 数组
// 文件内容：[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"},...]
func processLargeJSON(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)

    // 读取 '[' 开始符号
    token, err := decoder.Token()
    if err != nil {
        return err
    }
    if delim, ok := token.(json.Delim); !ok || delim != '[' {
        return fmt.Errorf("expected '[', got %v", token)
    }

    type User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }

    count := 0
    for decoder.More() {
        var user User
        if err := decoder.Decode(&user); err != nil {
            return fmt.Errorf("record %d: %w", count, err)
        }
        count++
        // 处理每条记录...
    }

    fmt.Printf("processed %d records\n", count)
    return nil
}

// 流式写入大 JSON 数组
func writeLargeJSON(path string, users <-chan User) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    encoder := json.NewEncoder(writer)
    encoder.SetIndent("", "  ")

    writer.WriteString("[\n")
    first := true
    for user := range users {
        if !first {
            writer.WriteString(",\n")
        }
        first = false
        if err := encoder.Encode(user); err != nil {
            return err
        }
    }
    writer.WriteString("]\n")
    return nil
}
```

---

### 14.13 实战：安全写入（原子写入）

直接写入文件，如果中途断电或程序崩溃，可能导致文件损坏。安全做法是**先写临时文件，再原子重命名**：

```go
func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
    dir := filepath.Dir(path)

    // 在同一目录创建临时文件（确保在同一文件系统，rename 才是原子的）
    tmpFile, err := os.CreateTemp(dir, ".tmp-*")
    if err != nil {
        return err
    }
    tmpPath := tmpFile.Name()

    // 出错时清理临时文件
    defer func() {
        if err != nil {
            os.Remove(tmpPath)
        }
    }()

    // 写入数据
    if _, err = tmpFile.Write(data); err != nil {
        tmpFile.Close()
        return err
    }

    // 确保数据落盘
    if err = tmpFile.Sync(); err != nil {
        tmpFile.Close()
        return err
    }

    if err = tmpFile.Close(); err != nil {
        return err
    }

    // 设置权限
    if err = os.Chmod(tmpPath, perm); err != nil {
        return err
    }

    // 原子重命名
    return os.Rename(tmpPath, path)
}

// 使用
config := []byte(`{"port": 8080, "debug": false}`)
atomicWriteFile("config.json", config, 0644)
```

---

### 14.14 实战：监控文件变化

简单的轮询方案（生产环境推荐使用 `fsnotify` 库）：

```go
func watchFile(path string, interval time.Duration, callback func()) error {
    initialStat, err := os.Stat(path)
    if err != nil {
        return err
    }

    for {
        time.Sleep(interval)

        stat, err := os.Stat(path)
        if err != nil {
            return err
        }

        if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
            callback()
            initialStat = stat
        }
    }
}

// 使用
go watchFile("config.yaml", 2*time.Second, func() {
    fmt.Println("config file changed! reloading...")
    // 重新加载配置...
})
```

---

### 14.15 内存中的 I/O

有时候不需要真正的文件，可以在内存中使用相同的 Reader/Writer 接口：

```go
import (
    "bytes"
    "strings"
)

// strings.Reader —— 从字符串读取
reader := strings.NewReader("Hello, Go!")
data, _ := io.ReadAll(reader) // []byte("Hello, Go!")

// bytes.Buffer —— 可读可写的内存缓冲区
var buf bytes.Buffer
buf.WriteString("hello ")
buf.WriteString("world")
buf.WriteByte('!')
fmt.Println(buf.String()) // "hello world!"
fmt.Println(buf.Len())    // 12

// 从 Buffer 读取
line, _ := buf.ReadString(' ') // "hello "

// bytes.Reader —— 只读，支持 Seek
r := bytes.NewReader([]byte("Hello"))
r.Seek(2, io.SeekStart)
b, _ := r.ReadByte() // 'l'

// 常见场景：构造 HTTP 请求体
body := bytes.NewBufferString(`{"name":"Alice"}`)
req, _ := http.NewRequest("POST", url, body)
```

---

### 14.16 总结：选择合适的读写方式

| 场景                     | 推荐方式               | 内存占用     |
| ------------------------ | ---------------------- | ------------ |
| 小文件（< 10MB）整体读取 | `os.ReadFile`          | 文件大小     |
| 小文件整体写入           | `os.WriteFile`         | 文件大小     |
| 按行处理                 | `bufio.Scanner`        | 单行大小     |
| 大文件流式复制           | `io.Copy`              | 32KB（固定） |
| 大文件块处理             | `file.Read` + 固定 buf | buf 大小     |
| 高频写入（日志等）       | `bufio.Writer`         | 缓冲大小     |
| 大 JSON 流式解析         | `json.NewDecoder`      | 单条记录     |
| 大 CSV 逐行处理          | `csv.Reader.Read`      | 单行大小     |
| 构造内存数据             | `bytes.Buffer`         | 数据大小     |
| 需要原子写入             | 临时文件 + `os.Rename` | 文件大小     |


## 15. 反射（Reflection）

反射是程序在运行时检查自身类型信息和操纵任意值的能力。Go 通过 `reflect` 包提供反射支持。反射强大但应谨慎使用——它绕过了编译期类型检查，代码更难理解、性能更低。JSON 序列化、ORM 框架、依赖注入、配置映射等场景是反射的典型用武之地。

### 15.1 反射的两大核心：Type 与 Value

Go 中每个 `interface{}` 内部都包含两部分信息：**具体类型（Type）** 和 **具体值（Value）**。反射就是把这两部分拆开检查和操作。

```go
import "reflect"

var x float64 = 3.14

t := reflect.TypeOf(x)   // reflect.Type —— 描述类型
v := reflect.ValueOf(x)  // reflect.Value —— 包含实际值

fmt.Println(t)            // "float64"
fmt.Println(t.Kind())     // "float64"（底层种类）
fmt.Println(v)            // "3.14"
fmt.Println(v.Type())     // "float64"
fmt.Println(v.Float())    // 3.14
fmt.Println(v.Interface()) // 3.14（转回 interface{}）
```

**Type 与 Kind 的区别**：`Type` 是完整的类型名（如 `main.User`），`Kind` 是底层种类（如 `struct`）。自定义类型和底层类型的 Kind 相同，但 Type 不同。

```go
type UserID int64

var id UserID = 42
t := reflect.TypeOf(id)

fmt.Println(t.Name())   // "UserID"   —— 自定义类型名
fmt.Println(t.Kind())   // "int64"    —— 底层种类
fmt.Println(t.String()) // "main.UserID"
```

**所有 Kind 常量**：

```go
reflect.Invalid       // 零值 Value 的 Kind
reflect.Bool
reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64
reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr
reflect.Float32, reflect.Float64
reflect.Complex64, reflect.Complex128
reflect.Array
reflect.Chan
reflect.Func
reflect.Interface
reflect.Map
reflect.Pointer       // 也可写作 reflect.Ptr
reflect.Slice
reflect.String
reflect.Struct
reflect.UnsafePointer
```

---

### 15.2 反射三大定律

Rob Pike 提出了反射的三条定律，理解它们是掌握反射的关键。

**定律一：反射可以从接口值得到反射对象**

```go
var num int = 42
var x interface{} = num

v := reflect.ValueOf(x)  // 从 interface{} -> reflect.Value
t := reflect.TypeOf(x)   // 从 interface{} -> reflect.Type
```

**定律二：反射可以从反射对象还原为接口值**

```go
v := reflect.ValueOf(42)
i := v.Interface()       // reflect.Value -> interface{}
n := i.(int)             // interface{} -> int
fmt.Println(n)           // 42
```

**定律三：要通过反射修改值，反射对象必须是可设置的（settable）**

```go
var x float64 = 3.14
v := reflect.ValueOf(x)
fmt.Println(v.CanSet()) // false —— v 持有的是 x 的副本

// 传入指针，再 Elem() 取得指向的元素
v = reflect.ValueOf(&x).Elem()
fmt.Println(v.CanSet()) // true
v.SetFloat(2.718)
fmt.Println(x)          // 2.718 —— 原始变量被修改了
```

> `Elem()` 的含义：对指针类型，返回指针指向的值；对接口类型，返回接口持有的具体值。

---

### 15.3 检查与操作结构体

结构体反射是最常用的场景，几乎所有序列化/ORM 库都依赖它。

#### 遍历结构体字段

```go
type User struct {
    Name    string `json:"name" validate:"required"`
    Age     int    `json:"age" validate:"min=0,max=150"`
    Email   string `json:"email,omitempty" validate:"email"`
    private string // 未导出字段
}

func inspectStruct(v interface{}) {
    val := reflect.ValueOf(v)
    typ := val.Type()

    // 如果传入指针，解引用
    if typ.Kind() == reflect.Pointer {
        val = val.Elem()
        typ = typ.Elem()
    }

    if typ.Kind() != reflect.Struct {
        fmt.Println("not a struct")
        return
    }

    fmt.Printf("Struct: %s (%d fields)\n", typ.Name(), typ.NumField())
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)     // reflect.StructField —— 字段元信息
        value := val.Field(i)     // reflect.Value —— 字段的值

        fmt.Printf("  [%d] %-10s %-10s = %-10v exported=%-5v tag=%s\n",
            i,
            field.Name,
            field.Type,
            value.Interface(), // 注意：未导出字段调用 Interface() 会 panic
            field.IsExported(),
            field.Tag,
        )
    }
}

// 安全版本：处理未导出字段
for i := 0; i < typ.NumField(); i++ {
    field := typ.Field(i)
    value := val.Field(i)

    if field.IsExported() {
        fmt.Printf("  %s = %v\n", field.Name, value.Interface())
    } else {
        fmt.Printf("  %s = (unexported)\n", field.Name)
    }
}
```

#### 读取 Struct Tag

```go
type Config struct {
    Host    string `env:"APP_HOST" default:"localhost"`
    Port    int    `env:"APP_PORT" default:"8080"`
    Debug   bool   `env:"APP_DEBUG" default:"false"`
}

func readTag(v interface{}) {
    typ := reflect.TypeOf(v)
    if typ.Kind() == reflect.Pointer {
        typ = typ.Elem()
    }

    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)

        // Get 返回指定 key 的值
        envTag := field.Tag.Get("env")
        defaultTag := field.Tag.Get("default")
        fmt.Printf("%s: env=%q default=%q\n", field.Name, envTag, defaultTag)

        // Lookup 可以区分"没有这个 tag"和"tag 值为空"
        if val, ok := field.Tag.Lookup("validate"); ok {
            fmt.Printf("  validate: %s\n", val)
        }
    }
}
```

#### 动态设置结构体字段

```go
func setField(obj interface{}, fieldName string, value interface{}) error {
    val := reflect.ValueOf(obj)

    // 必须传入指针
    if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Struct {
        return fmt.Errorf("expected pointer to struct, got %s", val.Type())
    }

    val = val.Elem()
    field := val.FieldByName(fieldName)

    if !field.IsValid() {
        return fmt.Errorf("no field named %q", fieldName)
    }
    if !field.CanSet() {
        return fmt.Errorf("cannot set field %q (unexported?)", fieldName)
    }

    newVal := reflect.ValueOf(value)
    if !newVal.Type().AssignableTo(field.Type()) {
        return fmt.Errorf("type mismatch: cannot assign %s to %s",
            newVal.Type(), field.Type())
    }

    field.Set(newVal)
    return nil
}

// 使用
user := User{Name: "Alice", Age: 25}
setField(&user, "Name", "Bob")
setField(&user, "Age", 30)
fmt.Println(user) // {Bob 30 }
```

#### 按名称获取字段

```go
val := reflect.ValueOf(user)

// FieldByName
nameField := val.FieldByName("Name")
if nameField.IsValid() {
    fmt.Println(nameField.String())
}

// FieldByIndex —— 支持嵌套
// 比如嵌套结构体的第 0 个字段的第 1 个字段
field := val.FieldByIndex([]int{0, 1})
```

---

### 15.4 检查与调用方法和函数

#### 遍历方法

```go
type Calculator struct{}

func (c Calculator) Add(a, b int) int    { return a + b }
func (c Calculator) Mul(a, b int) int    { return a * b }
func (c *Calculator) Reset()             { /* ... */ }

func listMethods(v interface{}) {
    typ := reflect.TypeOf(v)
    fmt.Printf("Type %s has %d methods:\n", typ, typ.NumMethod())

    for i := 0; i < typ.NumMethod(); i++ {
        method := typ.Method(i)
        fmt.Printf("  %s %s\n", method.Name, method.Type)
    }
}

listMethods(Calculator{})
// Type main.Calculator has 2 methods:
//   Add func(main.Calculator, int, int) int
//   Mul func(main.Calculator, int, int) int

listMethods(&Calculator{})
// Type *main.Calculator has 3 methods:
//   Add   func(*main.Calculator, int, int) int
//   Mul   func(*main.Calculator, int, int) int
//   Reset func(*main.Calculator)
```

> 注意：值类型只能看到值接收者方法，指针类型能看到所有方法。

#### 动态调用方法

```go
func callMethod(obj interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
    val := reflect.ValueOf(obj)
    method := val.MethodByName(methodName)

    if !method.IsValid() {
        return nil, fmt.Errorf("method %q not found", methodName)
    }

    // 构造参数
    in := make([]reflect.Value, len(args))
    for i, arg := range args {
        in[i] = reflect.ValueOf(arg)
    }

    // 调用
    results := method.Call(in)

    // 提取返回值
    out := make([]interface{}, len(results))
    for i, result := range results {
        out[i] = result.Interface()
    }
    return out, nil
}

// 使用
calc := Calculator{}
result, _ := callMethod(calc, "Add", 10, 20)
fmt.Println(result[0]) // 30
```

#### 动态调用函数

```go
func dynamicCall(fn interface{}, args ...interface{}) []interface{} {
    fnVal := reflect.ValueOf(fn)
    if fnVal.Kind() != reflect.Func {
        panic("not a function")
    }

    in := make([]reflect.Value, len(args))
    for i, arg := range args {
        in[i] = reflect.ValueOf(arg)
    }

    results := fnVal.Call(in)

    out := make([]interface{}, len(results))
    for i, r := range results {
        out[i] = r.Interface()
    }
    return out
}

// 使用
add := func(a, b int) int { return a + b }
result := dynamicCall(add, 3, 4)
fmt.Println(result[0]) // 7

// 检查函数签名
fnType := reflect.TypeOf(add)
fmt.Println(fnType.NumIn())       // 2（参数个数）
fmt.Println(fnType.In(0))         // int（第一个参数类型）
fmt.Println(fnType.NumOut())      // 1（返回值个数）
fmt.Println(fnType.Out(0))        // int（第一个返回值类型）
fmt.Println(fnType.IsVariadic())  // false
```

---

### 15.5 操作 Slice、Map、Channel

#### Slice

```go
// 创建 slice
sliceType := reflect.SliceOf(reflect.TypeOf(0)) // []int
slice := reflect.MakeSlice(sliceType, 0, 10)    // make([]int, 0, 10)

// 追加元素
slice = reflect.Append(slice, reflect.ValueOf(1), reflect.ValueOf(2), reflect.ValueOf(3))

fmt.Println(slice.Len()) // 3
fmt.Println(slice.Index(0).Int()) // 1

// 设置元素
slice.Index(1).SetInt(20)

// 转回 Go 值
goSlice := slice.Interface().([]int) // [1, 20, 3]

// 遍历
for i := 0; i < slice.Len(); i++ {
    fmt.Println(slice.Index(i).Interface())
}
```

#### Map

```go
// 创建 map
mapType := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(0)) // map[string]int
m := reflect.MakeMap(mapType)

// 设置键值
m.SetMapIndex(reflect.ValueOf("alice"), reflect.ValueOf(90))
m.SetMapIndex(reflect.ValueOf("bob"), reflect.ValueOf(85))

// 获取值
val := m.MapIndex(reflect.ValueOf("alice"))
if val.IsValid() {
    fmt.Println(val.Int()) // 90
}

// 遍历
iter := m.MapRange()
for iter.Next() {
    fmt.Printf("%s: %d\n", iter.Key().String(), iter.Value().Int())
}

// 也可以用 MapKeys
for _, key := range m.MapKeys() {
    fmt.Printf("%s: %v\n", key, m.MapIndex(key))
}

// 删除键
m.SetMapIndex(reflect.ValueOf("bob"), reflect.Value{}) // 设置零值 Value 即删除

fmt.Println(m.Len()) // 1
```

#### Channel

```go
// 创建 channel
chanType := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(0)) // chan int
ch := reflect.MakeChan(chanType, 5) // 缓冲 5

// 发送
ch.Send(reflect.ValueOf(42))

// 接收
val, ok := ch.TryRecv() // 非阻塞接收
if ok {
    fmt.Println(val.Int()) // 42
}

// reflect.Select —— 动态 select
cases := []reflect.SelectCase{
    {Dir: reflect.SelectRecv, Chan: ch},
    {Dir: reflect.SelectDefault},
}
chosen, value, recvOK := reflect.Select(cases)
fmt.Println(chosen, value, recvOK)
```

---

### 15.6 动态创建类型和值

```go
// 创建指针
intVal := reflect.ValueOf(42)
ptr := reflect.New(intVal.Type()) // *int
ptr.Elem().Set(intVal)
fmt.Println(ptr.Elem().Int())    // 42

// 创建结构体实例
type Order struct {
    ID     int
    Amount float64
}
orderType := reflect.TypeOf(Order{})
newOrder := reflect.New(orderType).Elem() // 创建 Order 零值
newOrder.FieldByName("ID").SetInt(1001)
newOrder.FieldByName("Amount").SetFloat(99.9)
order := newOrder.Interface().(Order)
fmt.Println(order) // {1001 99.9}

// 通过 reflect.StructOf 动态构建结构体类型（Go 1.7+）
dynamicType := reflect.StructOf([]reflect.StructField{
    {
        Name: "Name",
        Type: reflect.TypeOf(""),
        Tag:  `json:"name"`,
    },
    {
        Name: "Value",
        Type: reflect.TypeOf(0),
        Tag:  `json:"value"`,
    },
})

instance := reflect.New(dynamicType).Elem()
instance.FieldByName("Name").SetString("temperature")
instance.FieldByName("Value").SetInt(25)
fmt.Println(instance.Interface())
```

---

### 15.7 类型判断与转换

```go
// 类型断言 vs 反射
func describe(i interface{}) string {
    // switch 类型断言：编译期已知的类型列表
    switch v := i.(type) {
    case int:
        return fmt.Sprintf("int: %d", v)
    case string:
        return fmt.Sprintf("string: %q", v)
    default:
        // 反射：处理运行时才知道的类型
        val := reflect.ValueOf(i)
        return fmt.Sprintf("%s: %v", val.Type(), val)
    }
}

// 类型比较
t1 := reflect.TypeOf(0)
t2 := reflect.TypeOf(int(0))
fmt.Println(t1 == t2) // true

// 是否实现某接口
var writerType = reflect.TypeOf((*io.Writer)(nil)).Elem()
fileType := reflect.TypeOf((*os.File)(nil))
fmt.Println(fileType.Implements(writerType)) // true

// 是否可赋值 / 可转换
fmt.Println(t1.AssignableTo(t2))  // true
fmt.Println(t1.ConvertibleTo(reflect.TypeOf(float64(0)))) // true

// 类型转换
intVal := reflect.ValueOf(42)
floatVal := intVal.Convert(reflect.TypeOf(float64(0)))
fmt.Println(floatVal.Float()) // 42.0
```

---

### 15.8 实战：通用的 struct-to-map 转换器

```go
// StructToMap 将结构体转为 map[string]interface{}
// 支持 json tag 作为 key，支持嵌套结构体，支持忽略零值
func StructToMap(obj interface{}, useJSONTag bool) map[string]interface{} {
    result := make(map[string]interface{})
    val := reflect.ValueOf(obj)
    typ := val.Type()

    if typ.Kind() == reflect.Pointer {
        val = val.Elem()
        typ = typ.Elem()
    }
    if typ.Kind() != reflect.Struct {
        return result
    }

    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        value := val.Field(i)

        if !field.IsExported() {
            continue
        }

        // 确定 key 名称
        key := field.Name
        if useJSONTag {
            jsonTag := field.Tag.Get("json")
            if jsonTag == "-" {
                continue // 跳过
            }
            parts := strings.Split(jsonTag, ",")
            if parts[0] != "" {
                key = parts[0]
            }
            // 处理 omitempty
            if len(parts) > 1 && parts[1] == "omitempty" && value.IsZero() {
                continue
            }
        }

        // 递归处理嵌套结构体
        if value.Kind() == reflect.Struct && field.Type.Name() != "Time" {
            result[key] = StructToMap(value.Interface(), useJSONTag)
        } else {
            result[key] = value.Interface()
        }
    }
    return result
}

// 使用
type Address struct {
    City    string `json:"city"`
    ZipCode string `json:"zip_code"`
}

type Employee struct {
    Name    string  `json:"name"`
    Age     int     `json:"age,omitempty"`
    Salary  float64 `json:"salary"`
    Address Address `json:"address"`
    Secret  string  `json:"-"`
}

emp := Employee{
    Name:    "Alice",
    Salary:  50000,
    Address: Address{City: "Beijing", ZipCode: "100000"},
    Secret:  "should be hidden",
}

m := StructToMap(emp, true)
// map[name:Alice salary:50000 address:map[city:Beijing zip_code:100000]]
// Age 因 omitempty 被忽略，Secret 因 "-" 被忽略
```

---

### 15.9 实战：基于 struct tag 的验证器

```go
// 简易验证器：支持 required、min、max、email 标签
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("field %q: %s", e.Field, e.Message)
}

func Validate(obj interface{}) []ValidationError {
    var errs []ValidationError
    val := reflect.ValueOf(obj)
    typ := val.Type()

    if typ.Kind() == reflect.Pointer {
        val = val.Elem()
        typ = typ.Elem()
    }

    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        value := val.Field(i)
        tag := field.Tag.Get("validate")

        if tag == "" || !field.IsExported() {
            continue
        }

        rules := strings.Split(tag, ",")
        for _, rule := range rules {
            rule = strings.TrimSpace(rule)

            switch {
            case rule == "required":
                if value.IsZero() {
                    errs = append(errs, ValidationError{
                        Field:   field.Name,
                        Message: "is required",
                    })
                }

            case strings.HasPrefix(rule, "min="):
                minStr := strings.TrimPrefix(rule, "min=")
                minVal, _ := strconv.Atoi(minStr)
                switch value.Kind() {
                case reflect.Int, reflect.Int64:
                    if value.Int() < int64(minVal) {
                        errs = append(errs, ValidationError{
                            Field:   field.Name,
                            Message: fmt.Sprintf("must be >= %d", minVal),
                        })
                    }
                case reflect.String:
                    if len(value.String()) < minVal {
                        errs = append(errs, ValidationError{
                            Field:   field.Name,
                            Message: fmt.Sprintf("length must be >= %d", minVal),
                        })
                    }
                }

            case strings.HasPrefix(rule, "max="):
                maxStr := strings.TrimPrefix(rule, "max=")
                maxVal, _ := strconv.Atoi(maxStr)
                switch value.Kind() {
                case reflect.Int, reflect.Int64:
                    if value.Int() > int64(maxVal) {
                        errs = append(errs, ValidationError{
                            Field:   field.Name,
                            Message: fmt.Sprintf("must be <= %d", maxVal),
                        })
                    }
                }

            case rule == "email":
                email := value.String()
                if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
                    errs = append(errs, ValidationError{
                        Field:   field.Name,
                        Message: "is not a valid email",
                    })
                }
            }
        }
    }
    return errs
}

// 使用
type RegisterForm struct {
    Username string `validate:"required,min=3"`
    Age      int    `validate:"required,min=0,max=150"`
    Email    string `validate:"required,email"`
}

form := RegisterForm{Username: "ab", Age: -1, Email: "invalid"}
errs := Validate(form)
for _, e := range errs {
    fmt.Println(e.Error())
}
// field "Username": length must be >= 3
// field "Age": must be >= 0
// field "Email": is not a valid email
```

---

### 15.10 实战：简易依赖注入容器

```go
type Container struct {
    providers map[reflect.Type]interface{}
}

func NewContainer() *Container {
    return &Container{providers: make(map[reflect.Type]interface{})}
}

// Register 注册一个实例
func (c *Container) Register(value interface{}) {
    t := reflect.TypeOf(value)
    c.providers[t] = value
}

// Resolve 自动填充结构体字段（通过 inject tag）
func (c *Container) Resolve(target interface{}) error {
    val := reflect.ValueOf(target)
    if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Struct {
        return fmt.Errorf("target must be pointer to struct")
    }

    val = val.Elem()
    typ := val.Type()

    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        if field.Tag.Get("inject") != "auto" {
            continue
        }

        fieldType := field.Type
        provider, ok := c.providers[fieldType]
        if !ok {
            return fmt.Errorf("no provider for type %s", fieldType)
        }

        val.Field(i).Set(reflect.ValueOf(provider))
    }
    return nil
}

// 使用
type Logger struct{ Prefix string }
type Database struct{ DSN string }

type App struct {
    Logger   *Logger   `inject:"auto"`
    DB       *Database `inject:"auto"`
    Name     string    // 不注入
}

container := NewContainer()
container.Register(&Logger{Prefix: "[APP]"})
container.Register(&Database{DSN: "postgres://localhost/mydb"})

app := &App{Name: "MyApp"}
container.Resolve(app)

fmt.Println(app.Logger.Prefix)  // "[APP]"
fmt.Println(app.DB.DSN)         // "postgres://localhost/mydb"
fmt.Println(app.Name)           // "MyApp"
```

---

### 15.11 实战：环境变量映射到结构体

```go
// 通过反射 + struct tag 自动从环境变量填充配置
func LoadFromEnv(cfg interface{}) error {
    val := reflect.ValueOf(cfg)
    if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Struct {
        return fmt.Errorf("expected pointer to struct")
    }

    val = val.Elem()
    typ := val.Type()

    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        envKey := field.Tag.Get("env")
        if envKey == "" {
            continue
        }

        envVal, exists := os.LookupEnv(envKey)
        if !exists {
            // 使用 default tag
            if def := field.Tag.Get("default"); def != "" {
                envVal = def
            } else {
                continue
            }
        }

        fv := val.Field(i)
        switch fv.Kind() {
        case reflect.String:
            fv.SetString(envVal)
        case reflect.Int, reflect.Int64:
            if field.Type == reflect.TypeOf(time.Duration(0)) {
                d, err := time.ParseDuration(envVal)
                if err != nil {
                    return fmt.Errorf("field %s: %w", field.Name, err)
                }
                fv.SetInt(int64(d))
            } else {
                n, err := strconv.ParseInt(envVal, 10, 64)
                if err != nil {
                    return fmt.Errorf("field %s: %w", field.Name, err)
                }
                fv.SetInt(n)
            }
        case reflect.Bool:
            b, err := strconv.ParseBool(envVal)
            if err != nil {
                return fmt.Errorf("field %s: %w", field.Name, err)
            }
            fv.SetBool(b)
        case reflect.Float64:
            f, err := strconv.ParseFloat(envVal, 64)
            if err != nil {
                return fmt.Errorf("field %s: %w", field.Name, err)
            }
            fv.SetFloat(f)
        }
    }
    return nil
}

// 使用
type ServerConfig struct {
    Host         string        `env:"SERVER_HOST" default:"0.0.0.0"`
    Port         int           `env:"SERVER_PORT" default:"8080"`
    Debug        bool          `env:"SERVER_DEBUG" default:"false"`
    ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" default:"30s"`
}

os.Setenv("SERVER_PORT", "9090")
os.Setenv("SERVER_DEBUG", "true")

var cfg ServerConfig
LoadFromEnv(&cfg)
fmt.Printf("%+v\n", cfg)
// {Host:0.0.0.0 Port:9090 Debug:true ReadTimeout:30s}
```

---

### 15.12 实战：通用的深度比较与深拷贝

```go
// 深度比较：reflect.DeepEqual
a := []int{1, 2, 3}
b := []int{1, 2, 3}
fmt.Println(reflect.DeepEqual(a, b)) // true（直接 == 会编译错误）
fmt.Println(reflect.DeepEqual(a, []int{1, 2})) // false

// map 比较
m1 := map[string]int{"a": 1, "b": 2}
m2 := map[string]int{"b": 2, "a": 1}
fmt.Println(reflect.DeepEqual(m1, m2)) // true

// 注意：DeepEqual 对 nil 和空 slice 区分
var s1 []int            // nil
s2 := []int{}           // 空但非 nil
fmt.Println(reflect.DeepEqual(s1, s2)) // false

// 通用深拷贝（简化版，处理基本场景）
func DeepCopy(src interface{}) interface{} {
    srcVal := reflect.ValueOf(src)
    dst := reflect.New(srcVal.Type()).Elem()
    deepCopyValue(dst, srcVal)
    return dst.Interface()
}

func deepCopyValue(dst, src reflect.Value) {
    switch src.Kind() {
    case reflect.Pointer:
        if !src.IsNil() {
            dst.Set(reflect.New(src.Elem().Type()))
            deepCopyValue(dst.Elem(), src.Elem())
        }
    case reflect.Struct:
        for i := 0; i < src.NumField(); i++ {
            if dst.Field(i).CanSet() {
                deepCopyValue(dst.Field(i), src.Field(i))
            }
        }
    case reflect.Slice:
        if !src.IsNil() {
            dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
            for i := 0; i < src.Len(); i++ {
                deepCopyValue(dst.Index(i), src.Index(i))
            }
        }
    case reflect.Map:
        if !src.IsNil() {
            dst.Set(reflect.MakeMap(src.Type()))
            for _, key := range src.MapKeys() {
                newVal := reflect.New(src.MapIndex(key).Type()).Elem()
                deepCopyValue(newVal, src.MapIndex(key))
                dst.SetMapIndex(key, newVal)
            }
        }
    default:
        dst.Set(src)
    }
}
```

---

### 15.13 reflect.Value 的零值与有效性检查

```go
var v reflect.Value // 零值
fmt.Println(v.IsValid()) // false —— 零值 Value
fmt.Println(v.Kind())    // reflect.Invalid

// 查找不存在的字段/方法返回零值 Value
val := reflect.ValueOf(User{})
field := val.FieldByName("Nonexistent")
fmt.Println(field.IsValid()) // false

method := val.MethodByName("DoesNotExist")
fmt.Println(method.IsValid()) // false

// IsNil 只能用于 chan/func/interface/map/pointer/slice
var p *int
v = reflect.ValueOf(p)
fmt.Println(v.IsNil())  // true
fmt.Println(v.IsValid()) // true（Value 本身有效，只是持有的指针是 nil）

// IsZero（Go 1.13+）—— 检查是否为类型零值
fmt.Println(reflect.ValueOf(0).IsZero())     // true
fmt.Println(reflect.ValueOf("").IsZero())    // true
fmt.Println(reflect.ValueOf(false).IsZero()) // true
fmt.Println(reflect.ValueOf(User{}).IsZero()) // true（所有字段都是零值）
```

---

### 15.14 性能考量与最佳实践

#### 性能基准

```go
// 直接访问 vs 反射访问，性能差距约 50~100 倍
func BenchmarkDirect(b *testing.B) {
    u := User{Name: "Alice", Age: 30}
    for i := 0; i < b.N; i++ {
        _ = u.Name
    }
}

func BenchmarkReflect(b *testing.B) {
    u := User{Name: "Alice", Age: 30}
    val := reflect.ValueOf(u)
    for i := 0; i < b.N; i++ {
        _ = val.FieldByName("Name").String()
    }
}

// 优化：缓存 reflect.Type 和字段索引
var (
    userType      = reflect.TypeOf(User{})
    nameFieldIdx  int
)

func init() {
    field, _ := userType.FieldByName("Name")
    nameFieldIdx = field.Index[0]
}

func BenchmarkReflectCached(b *testing.B) {
    u := User{Name: "Alice", Age: 30}
    val := reflect.ValueOf(u)
    for i := 0; i < b.N; i++ {
        _ = val.Field(nameFieldIdx).String() // 用索引比用名称快
    }
}
```

#### 最佳实践

```go
// ✅ 1. 能用类型断言就不用反射
func process(v interface{}) {
    switch val := v.(type) {
    case string:
        handleString(val)
    case int:
        handleInt(val)
    default:
        handleViaReflect(v) // 最后手段
    }
}

// ✅ 2. 缓存 reflect.Type，避免重复调用 TypeOf
var cachedType = reflect.TypeOf((*MyStruct)(nil)).Elem()

// ✅ 3. 用 FieldByIndex([]int{i}) 替代 FieldByName("xxx")
//    FieldByName 内部需要遍历查找，FieldByIndex 直接定位

// ✅ 4. 用代码生成替代反射（性能关键路径）
//    如 easyjson、msgp 等工具在编译期生成序列化代码

// ✅ 5. 尽量在初始化阶段做反射，运行时使用缓存结果
type fieldMapping struct {
    index    int
    jsonName string
}

var mappingCache = make(map[reflect.Type][]fieldMapping)

func getMapping(t reflect.Type) []fieldMapping {
    if m, ok := mappingCache[t]; ok {
        return m
    }
    var mapping []fieldMapping
    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)
        if !f.IsExported() {
            continue
        }
        name := f.Tag.Get("json")
        if name == "" {
            name = f.Name
        }
        mapping = append(mapping, fieldMapping{index: i, jsonName: name})
    }
    mappingCache[t] = mapping
    return mapping
}
```

---

### 15.15 常见陷阱

```go
// ❌ 陷阱1：对不可寻址的值调用 Set
v := reflect.ValueOf(42)
// v.SetInt(100) // panic: reflect.Value.SetInt using unaddressable value

// ✅ 正确：通过指针
x := 42
reflect.ValueOf(&x).Elem().SetInt(100)

// ❌ 陷阱2：对未导出字段调用 Interface() 或 Set()
type secret struct {
    hidden int
}
v := reflect.ValueOf(secret{42})
// v.Field(0).Interface() // panic: unexported field

// ❌ 陷阱3：nil interface 与 nil pointer 的区别
var p *int = nil
var i interface{} = p

fmt.Println(i == nil)                       // false（interface 有类型信息）
fmt.Println(reflect.ValueOf(i).IsNil())     // true（指针值是 nil）
fmt.Println(reflect.ValueOf(nil).IsValid()) // false（真正的 nil interface）

// ❌ 陷阱4：忘记 Elem()
var num int = 42
v = reflect.ValueOf(&num)
fmt.Println(v.Kind())        // ptr
fmt.Println(v.Elem().Kind()) // int
// v.SetInt(100)             // panic —— v 是指针，不是 int
v.Elem().SetInt(100)         // 正确

// ❌ 陷阱5：Kind 判断错误
type MyString string
v = reflect.ValueOf(MyString("hello"))
// if v.Kind() == reflect.String  ✅（底层种类是 string）
// if v.Type() == reflect.TypeOf("") ❌（MyString != string）
```

---

### 15.16 总结速查

| 操作 | 方法 |
|------|------|
| 获取类型 | `reflect.TypeOf(v)` |
| 获取值 | `reflect.ValueOf(v)` |
| 底层种类 | `v.Kind()` |
| 还原为 interface{} | `v.Interface()` |
| 解引用指针 | `v.Elem()` |
| 是否可设置 | `v.CanSet()` |
| 字段数量 | `t.NumField()` / `v.NumField()` |
| 按名称取字段 | `v.FieldByName("X")` |
| 读取 tag | `t.Field(i).Tag.Get("json")` |
| 方法数量 | `t.NumMethod()` |
| 调用方法 | `v.MethodByName("M").Call(args)` |
| 创建新实例 | `reflect.New(t)` |
| 创建 slice | `reflect.MakeSlice(t, len, cap)` |
| 创建 map | `reflect.MakeMap(t)` |
| 深度比较 | `reflect.DeepEqual(a, b)` |
| 判断零值 | `v.IsZero()` |
| 判断 nil | `v.IsNil()`（仅指针/map/slice/chan/func/interface） |

---

## 16. 常用标准库概览

### 16.1 strings

```go
import "strings"

strings.Contains("hello", "ell")     // true
strings.HasPrefix("hello", "he")     // true
strings.HasSuffix("hello", "lo")     // true
strings.Index("hello", "ll")         // 2
strings.ToUpper("hello")             // "HELLO"
strings.ToLower("HELLO")             // "hello"
strings.TrimSpace("  hi  ")          // "hi"
strings.Trim("**hi**", "*")          // "hi"
strings.Replace("aaa", "a", "b", 2)  // "bba"
strings.ReplaceAll("aaa", "a", "b")  // "bbb"
strings.Split("a,b,c", ",")          // ["a", "b", "c"]
strings.Join([]string{"a","b"}, "-")  // "a-b"
strings.Repeat("ab", 3)              // "ababab"
strings.Count("hello", "l")          // 2

// strings.Builder（高效字符串拼接）
var sb strings.Builder
sb.WriteString("hello")
sb.WriteString(" world")
s := sb.String() // "hello world"
```

### 16.2 strconv

```go
import "strconv"

// int <-> string
strconv.Itoa(42)           // "42"
strconv.Atoi("42")         // 42, nil

// Parse 系列
strconv.ParseBool("true")      // true, nil
strconv.ParseFloat("3.14", 64) // 3.14, nil
strconv.ParseInt("FF", 16, 64) // 255, nil

// Format 系列
strconv.FormatBool(true)         // "true"
strconv.FormatFloat(3.14, 'f', 2, 64) // "3.14"
strconv.FormatInt(255, 16)       // "ff"
```

### 16.3 sort

```go
import "sort"

// 基本排序
nums := []int{3, 1, 4, 1, 5}
sort.Ints(nums)             // [1, 1, 3, 4, 5]
sort.Strings([]string{"b", "a", "c"})

// 自定义排序
sort.Slice(nums, func(i, j int) bool {
    return nums[i] > nums[j] // 降序
})

// 搜索（需已排序）
idx := sort.SearchInts(nums, 3)

// Go 1.21+ slices 包
import "slices"
slices.Sort(nums)
slices.Contains(nums, 3)
```

### 16.4 time

```go
import "time"

// 获取时间
now := time.Now()
t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// 格式化（Go 用固定的参考时间 "2006-01-02 15:04:05"）
now.Format("2006-01-02 15:04:05")
now.Format(time.RFC3339)

// 解析
t, err := time.Parse("2006-01-02", "2024-03-15")

// 时间运算
d := 2 * time.Hour + 30 * time.Minute
future := now.Add(d)
diff := future.Sub(now)  // Duration

// Sleep
time.Sleep(1 * time.Second)

// Ticker / Timer
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for t := range ticker.C {
    fmt.Println("tick:", t)
}

timer := time.NewTimer(5 * time.Second)
<-timer.C // 5 秒后触发
```

### 16.5 encoding/json

```go
import "encoding/json"

type User struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email,omitempty"`
}

// 序列化
u := User{Name: "Alice", Age: 25}
data, err := json.Marshal(u)          // []byte
pretty, err := json.MarshalIndent(u, "", "  ") // 格式化

// 反序列化
var u2 User
err = json.Unmarshal(data, &u2)

// 处理动态 JSON
var result map[string]interface{}
json.Unmarshal(data, &result)

// 编码到 Writer / 从 Reader 解码
json.NewEncoder(os.Stdout).Encode(u)
json.NewDecoder(resp.Body).Decode(&u2)
```

### 16.6 net/http

```go
import "net/http"

// --- HTTP 服务器 ---
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
})
http.ListenAndServe(":8080", nil)

// 使用 ServeMux
mux := http.NewServeMux()
mux.HandleFunc("/api/users", handleUsers)
http.ListenAndServe(":8080", mux)

// --- HTTP 客户端 ---
resp, err := http.Get("https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)

// 自定义请求
client := &http.Client{Timeout: 10 * time.Second}
req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
req.Header.Set("Content-Type", "application/json")
resp, err := client.Do(req)
```

### 16.7 os / os/exec

```go
import "os"

// 环境变量
os.Getenv("HOME")
os.Setenv("KEY", "value")

// 命令行参数
args := os.Args // os.Args[0] 是程序路径

// 退出
os.Exit(1)

// 执行外部命令
import "os/exec"
out, err := exec.Command("ls", "-la").Output()
fmt.Println(string(out))
```

### 16.8 log

```go
import "log"

log.Println("info message")
log.Printf("user: %s\n", name)
log.Fatal("fatal error")  // 打印后 os.Exit(1)
log.Panic("panic!")        // 打印后 panic

// 自定义 logger
logger := log.New(os.Stdout, "[APP] ", log.Ldate|log.Ltime|log.Lshortfile)
logger.Println("custom log")

// Go 1.21+ slog 结构化日志
import "log/slog"
slog.Info("user login", "name", "Alice", "age", 25)
```

### 16.9 regexp

```go
import "regexp"

re := regexp.MustCompile(`\d+`)
re.MatchString("abc123")          // true
re.FindString("abc123def456")     // "123"
re.FindAllString("abc123def456", -1) // ["123", "456"]
re.ReplaceAllString("abc123", "X")   // "abcX"
```

### 16.10 math/rand 与 crypto/rand

```go
// 伪随机（Go 1.20+ 自动 seed）
import "math/rand"
rand.Intn(100)        // [0, 100)
rand.Float64()        // [0.0, 1.0)

// 密码学安全随机
import "crypto/rand"
import "math/big"
n, _ := rand.Int(rand.Reader, big.NewInt(100))
```

---

## 17. 测试

### 17.1 单元测试

文件命名 `xxx_test.go`，函数命名 `TestXxx(t *testing.T)`。

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}

// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}

// 表驱动测试（推荐）
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.expected)
            }
        })
    }
}
```

```bash
go test ./...          # 运行所有测试
go test -v ./...       # 详细输出
go test -run TestAdd   # 运行匹配的测试
go test -cover         # 覆盖率
```

### 17.2 基准测试

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

```bash
go test -bench=. -benchmem
```

### 17.3 TestMain

```go
func TestMain(m *testing.M) {
    // 测试前的 setup
    fmt.Println("setup")

    code := m.Run() // 运行所有测试

    // 测试后的 teardown
    fmt.Println("teardown")

    os.Exit(code)
}
```

---

## 18. 包、模块与工程化管理

### 18.1 包的概念与组织

Go 语言中的包（package）是代码组织的基本单位，不仅提供了命名空间隔离，还支持代码重用和模块化设计。

每个 Go 文件都必须声明所属的包，因此同一个目录下的 Go 文件，通常应该属于同一个包，包名最好与所在目录名相同，一个目录通常对应一个包。

### 18.2 基本使用

```go
// 声明包
package mypackage

// 导入
import (
    "fmt"
    "math/rand"
    "os"

    // 别名
    myjson "encoding/json"

    // 匿名导入（仅执行 init）
    _ "github.com/lib/pq"

    // 点导入（不推荐）
    . "fmt" // 可以直接用 Println 而非 fmt.Println
)

// Go 中每个包都可以定义 `init()` 函数，在包初始化时自动执行，一个包中可以有多个 `init`，在 `main()` 执行前运行
func init() {
	fmt.Println("init 执行")
}

```

**可见性规则**：**大写开头 = 导出（public），小写开头 = 未导出（private）**。

```go
type User struct {     // 导出
    Name  string       // 导出
    email string       // 未导出
}

func Hello() {}       // 导出
func helper() {}      // 未导出
```


### 18.3 go mod

Go Modules 是 Go 1.11 引入的官方依赖管理工具，解决了 GOPATH 模式下的依赖混乱问题。

```bash
# 初始化模块
go mod init github.com/username/project
# 添加依赖
go get github.com/some/dependency@v1.2.3
go get github.com/gin-gonic/gin@latest

# 更新依赖
go get -u github.com/some/dependency

# 整理依赖，自动添加需要的依赖，删除未使用的依赖，更新 `go.sum`
go mod tidy

# 查看依赖
go list -m all

# 下载依赖
go mod download

# 替换依赖（本地开发调试,使用 fork 版本等）
go mod edit -replace github.com/old/pkg=../local/pkg
go mod edit -replace github.com/old/pkg=github.com/new/pkg@v1.0.0
```

`go.mod` 文件示例：

```
module github.com/user/project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
)
```

### 18.4 项目结构惯例

众所周知，项目结构各有各的风格，这里给出一个常见的结构示例，供参考。

**简易版**

```
myproject/
├── go.mod
├── go.sum
├── main.go
├── cmd/            # 多个可执行入口
│   └── server/
│       └── main.go
├── internal/       # 私有包（仅本模块可导入）
│   └── auth/
├── pkg/            # 公共包（可被外部导入）
│   └── utils/
├── api/            # API 定义
├── configs/        # 配置文件
└── test/           # 额外测试
```

**较完整版**

```
myproject/
├── cmd/                        # 可执行程序入口（每个子目录一个二进制）
│   ├── server/                 # 主 HTTP API 服务
│   │   └── main.go
│   ├── grpcserver/             # gRPC 服务（独立进程）
│   │   └── main.go
│   ├── gateway/                # gRPC-Gateway（HTTP↔gRPC 转换代理）
│   │   └── main.go
│   ├── worker/                 # 后台任务 / 消息消费者
│   │   └── main.go
│   ├── scheduler/              # 定时任务 / Cron
│   │   └── main.go
│   ├── migrate/                # 数据库迁移工具
│   │   └── main.go
│   └── cli/                    # 命令行工具
│       └── main.go
│
├── internal/                   # 内部包（Go 编译器强制：仅本模块可导入）
│   ├── config/                 # 配置加载与结构定义
│   │   ├── config.go           #   Config struct + Load()
│   │   ├── env.go              #   EnvConfig struct + Load()
│   │   └── loader.go           #   配置加载器
│   │
│   │── router/                 # HTTP 路由注册（聚合所有 handler）
│   │   └── router.go           #   func NewRouter() *mux.Router
│   ├── middleware/             # HTTP 中间件
│   │   ├── auth.go             #   JWT / Session 校验
│   │   ├── cors.go             #   跨域处理
│   │   ├── logging.go          #   请求日志
│   │   ├── recovery.go         #   panic 恢复
│   │   └── ratelimit.go        #   限流
│   ├── handler/                # HTTP 处理器（控制器层）
│   │   ├── user.go             #   UserHandler：注册/登录/查询
│   │   ├── order.go            #   OrderHandler
│   │   └── health.go           #   健康检查 /healthz
│   │
│   ├── grpc/                   # gRPC 服务层
│   │   ├── server.go           #   gRPC Server 初始化 + 注册
│   │   ├── interceptor/        #   gRPC 拦截器（等同 HTTP 中间件）
│   │   │   ├── auth.go         #     认证拦截器
│   │   │   ├── logging.go      #     日志拦截器
│   │   │   └── recovery.go     #     panic 恢复
│   │   └── service/            #   gRPC Service 实现
│   │       ├── user.go         #     实现 pb.UserServiceServer
│   │       └── order.go
│   │
│   ├── event/                  # 事件系统
│   │   ├── event.go            #   事件类型定义（Event interface / 各事件 struct）
│   │   ├── bus.go              #   事件总线接口 + 内存实现
│   │   ├── publisher.go        #   发布者（封装 Kafka/RabbitMQ/NATS 等）
│   │   ├── subscriber.go       #   订阅者注册
│   │   └── handler/            #   事件处理器
│   │       ├── user_created.go #     处理 UserCreated 事件
│   │       ├── order_paid.go   #     处理 OrderPaid 事件
│   │       └── notification.go #     发送通知
│   │
│   ├── mq/                     # 消息队列适配层
│   │   ├── producer.go         #   Producer interface
│   │   ├── consumer.go         #   Consumer interface
│   │   ├── kafka.go            #   Kafka 实现
│   │   ├── rabbitmq.go         #   RabbitMQ 实现
│   │   └── nats.go             #   NATS 实现
│   │
│   ├── service/                # 业务逻辑层（核心领域逻辑）
│   │   ├── user.go             #   UserService interface + impl
│   │   └── order.go
│   ├── repository/             # 数据访问层（DB / 缓存操作）
│   │   ├── user.go             #   UserRepository interface + impl
│   │   └── order.go
│   ├── model/                  # 数据模型 / 领域实体
│   │   ├── user.go             #   type User struct
│   │   └── order.go
│   ├── dto/                    # 请求/响应数据传输对象
│   │   ├── request.go          #   CreateUserRequest 等
│   │   └── response.go         #   UserResponse / ErrorResponse
│   │
│   ├── auth/                   # 认证与授权
│   │   ├── jwt.go
│   │   ├── rbac.go
│   │   └── oauth.go            #   第三方 OAuth
│   ├── database/               # 数据库连接与初始化
│   │   ├── mysql.go
│   │   ├── redis.go
│   │   └── mongo.go
│   ├── cache/                  # 缓存抽象层
│   │   ├── cache.go            #   Cache interface
│   │   ├── redis.go            #   Redis 实现
│   │   └── memory.go           #   本地内存实现
│   ├── notify/                 # 通知推送
│   │   ├── notifier.go         #   Notifier interface
│   │   ├── email.go            #   邮件
│   │   ├── sms.go              #   短信
│   │   └── webhook.go          #   Webhook 回调
│   ├── observability/          # 可观测性
│   │   ├── metrics.go          #   Prometheus 指标
│   │   └── tracing.go          #   OpenTelemetry 链路追踪
│   ├── errors/                 # 自定义错误类型与错误码
│   │   └── errors.go
│   ├── cron/                   # 定时任务定义
│   │   ├── scheduler.go        #   任务注册与调度
│   │   └── jobs/
│   │       ├── cleanup.go      #   过期数据清理
│   │       └── report.go       #   定时报表
│   └── util/                   # 内部工具函数
│       ├── hash.go
│       └── validator.go
│
├── pkg/                        # 公共包（可被外部项目导入）
│   ├── logger/                 #   统一日志封装
│   │   └── logger.go
│   ├── httpclient/             #   HTTP 客户端封装
│   │   └── client.go
│   ├── pagination/             #   分页工具
│   └── retry/                  #   重试工具
│       └── retry.go
│
├── api/                        # API 契约定义
│   ├── openapi/                #   OpenAPI / Swagger
│   │   └── spec.yaml
│   └── proto/                  #   Protobuf 源文件
│       ├── user/
│       │   └── user.proto      #   service UserService { ... }
│       ├── order/
│       │   └── order.proto
│       └── common/
│           └── common.proto    #   公共 message（Pagination 等）
│
├── gen/                        # 自动生成代码（勿手动编辑）
│   └── proto/                  #   protoc 生成的 Go 代码
│       ├── user/
│       │   ├── user.pb.go
│       │   └── user_grpc.pb.go
│       ├── order/
│       │   ├── order.pb.go
│       │   └── order_grpc.pb.go
│       └── common/
│           └── common.pb.go
│
├── migrations/                 # 数据库迁移文件
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
│
├── configs/                    # 配置文件（按环境区分）
│   ├── config.yaml             #   默认配置
│   ├── config.dev.yaml         #   开发环境
│   ├── config.prod.yaml        #   生产环境
│   └── config.test.yaml        #   测试环境
│
├── deployments/                # 部署相关文件
│   ├── docker/
│   │   ├── Dockerfile.server
│   │   ├── Dockerfile.grpc
│   │   └── Dockerfile.worker
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   └── docker-compose.yaml
│
├── scripts/                    # 构建与运维脚本
│   ├── build.sh
│   ├── test.sh
│   ├── protogen.sh             #   protoc 代码生成脚本
│   └── seed.go                 #   填充测试数据
│
├── docs/                       # 项目文档
│   ├── architecture.md
│   ├── api.md
│   └── changelog.md
│
├── test/                       # 集成测试 / E2E 测试
│   ├── integration/
│   │   └── user_test.go
│   └── testdata/               #   测试固定数据
│       └── users.json
│
├── web/                        # 前端静态资源（如有）
│   ├── static/
│   └── templates/
│
├── .gitignore
├── .golangci.yml
├── buf.yaml                    # buf 工具配置（protobuf 管理）
├── buf.gen.yaml                # buf 代码生成配置
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## 19. Go 常见易错点总结

### 已声明但未使用

Go 对未使用变量和未使用导入非常严格。

```go
x := 10 // 如果不用，会编译报错
```

这是 Go 强调代码整洁的一种体现。

### map 未初始化不能直接赋值

错误示例：

```go
var m map[string]int
m["a"] = 1
```

正确写法：

```go
m := make(map[string]int)
m["a"] = 1
```

### 切片 append 后可能指向新底层数组

```go
s := []int{1, 2, 3}
s = append(s, 4)
```

如果底层数组容量不足，`append` 会分配新数组，因此必须重新接收返回值。

数组是值类型，切片是引用语义

数组赋值是整份复制，切片赋值共享底层数据，这一点非常容易混淆。

### defer 参数求值时机

```go
x := 1
defer fmt.Println(x)
x = 2
```

输出为：

```bash
1
```

因为 `defer` 在声明时就已经确定参数值。

### range 变量问题

在某些循环场景中，`range` 中的迭代变量是复用的，闭包引用时要特别小心。写并发代码时尤其容易踩坑。

例如应尽量这样写：

```go
for _, v := range nums {
	v := v
	go func() {
		fmt.Println(v)
	}()
}
```

### nil 问题

以下类型可能为 `nil`：

- 指针
- 切片
- map
- channel
- 函数
- 接口

使用前要注意是否已经初始化。

**接口值不等于 nil 的陷阱**

一个接口如果内部持有“带类型的 nil 指针”，接口本身不一定等于 `nil`。这是 Go 中较经典的坑，实际开发中需要谨慎判断。

```go
var i interface{}
fmt.Println(i == nil) // true

i = (*struct{})(nil)
fmt.Println(i == nil) // false
```

---

## 20. 编码规范与最佳实践

### 20.1 代码格式化

```bash
gofmt -w .    # 格式化代码
goimports -w . # 格式化 + 自动管理 import
```

Go 有且仅有一种代码风格，由 `gofmt` 强制统一。

### 20.2 命名规范

```go
// 包名：小写、简短、单个单词
package http

// 变量/函数：驼峰命名，不用下划线
var userName string
func getUserName() string

// 首字母大写 = 导出，小写 = 私有

// 接口命名：单方法接口以 -er 结尾
type Reader interface { Read(p []byte) (n int, err error) }
type Stringer interface { String() string }

// 缩写保持全大写或全小写
var httpURL string
var xmlParser XMLParser
```

### 20.3 错误处理惯例

```go
// 错误变量以 Err 开头
var ErrNotFound = errors.New("not found")
var ErrTimeout  = errors.New("timeout")

// 错误类型以 Error 结尾
type NotFoundError struct { ... }

// 不要忽略错误
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething failed: %w", err)
}

// 尽早返回，减少嵌套
func process() error {
    if err := step1(); err != nil {
        return err
    }
    if err := step2(); err != nil {
        return err
    }
    return nil
}
```

### 20.4 实用的 Go 命令

```bash
go run main.go         # 编译并运行
go build               # 编译
go install             # 编译并安装到 $GOPATH/bin
go fmt ./...           # 格式化
go vet ./...           # 静态分析
go test ./...          # 测试
go mod tidy            # 清理依赖
go doc fmt.Println     # 查看文档
go generate            # 执行 //go:generate 指令
go env                 # 查看环境变量
go tool pprof          # 性能分析
```

### 20.5 常用编译标签与构建

```bash
# 交叉编译
GOOS=linux GOARCH=amd64 go build -o app
GOOS=windows GOARCH=amd64 go build -o app.exe
GOOS=darwin GOARCH=arm64 go build -o app

# 减小体积
go build -ldflags="-s -w" -o app

# 编译时注入变量
go build -ldflags="-X main.version=1.0.0" -o app
```

### 20.6 核心原则

> 不要通过共享内存来通信，而要通过通信来共享内存。

> 少量复制好过引入依赖。

> 让零值有意义（如 `sync.Mutex` 零值就是未锁状态）。

> 参数用接口，返回值用具体类型。

---

*本文基于 Go 1.21+，涵盖了 Go 语言的核心语法知识，大量示例。建议收藏作为速查手册，配合 [官方文档](https://go.dev/doc/) 和 [Go by Example](https://gobyexample.com/) 一起学习。*