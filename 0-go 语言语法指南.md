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

// 类型别名（完全等同于原类型）
type Byte = uint8
type Rune = int32
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

```go
// 定义
type User struct {
    Name string
    Age  int
    Email string
}

// 初始化
u1 := User{"Alice", 25, "alice@go.dev"}    // 按顺序
u2 := User{Name: "Bob", Age: 30}           // 按字段名（推荐）
u3 := new(User)                            // 返回 *User，字段为零值
var u4 User                                // 零值结构体

// 访问与修改
u2.Email = "bob@go.dev"

// 指针访问（自动解引用）
p := &u2
p.Name = "Bob2" // 等同于 (*p).Name = "Bob2"

// 匿名字段（嵌入/组合）
type Admin struct {
    User          // 嵌入 User
    Level int
}
a := Admin{User: User{Name: "Root", Age: 0}, Level: 1}
a.Name // 直接访问 User 的字段（提升）

// 匿名结构体
point := struct {
    X, Y int
}{10, 20}

// 结构体标签（用于 JSON 序列化等）
type Product struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price float64 `json:"price,omitempty"`
}
```

### 5.5 自定义类型与类型别名

```go
// 自定义类型（全新类型，可以添加方法）
type Celsius float64
type Handler func(string) error

// 类型别名（完全等同于原类型）
type Byte = uint8
type Rune = int32
```

---

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
