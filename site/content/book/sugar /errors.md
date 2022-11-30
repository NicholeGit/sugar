---
title: "Errors"
---

# Errors

这个库可以作为 [github.com/pkg/errors](github.com/pkg/errors) 和 Go 标准错误包的替代。

在 `Golang` 中，我们能够自由的处理程序产生的 `error`。由于 `error` 的处理对于程序来说至关重要，我们需要根据应用程序的相关性质通过调整或者添加更多保存在 `Error` 对象中的数据来提高 `error` 的可追踪性。

## `Error` 结构

一个用于来包裹需要传递的错误的自定义结构。

```golang
type Operation string
type Kind string

type Error struct {
	// 错误发生时所处的函数名。
	op Operation

	// 包裹的错误的类型。例如，`未找到` 的类型为：`NotFound`。
	kind Kind

	// 对包裹的错误的解释。
	message string

	// 包裹的错误对象。
	cause error

	// 发生错误的具体位置。例如：产生错误的文件名和对应具体代码的行号。
	location *Location
}
```
- `Location` 不应该被传递。默认情况下，将通过 `runtime.Caller` 获取产生错误的文件名和对应具体代码的行号。
- `Op` 如果未被指定，将通过 `runtime.Caller` 获取。

## 构造一个 `Error`
为了方便 `Error` 的构造，包提供了一个便于输入的函数：`E`

	func E(args ...interface{}) error

## 包装一个 `Error`
提供 `Wrap` 函数包装 error， 这样可以记录 error 的错误链，在最上层打印 log，就可以把整个错误链输出，无须在每个错误的地方都打印 log。

    func Wrap(err error, args ...interface{}) error {

## 如何使用
```golang
// app/constants.go
const (
	KindInternal        Kind = "INTERNAL"         // Internal error or inconsistency.
	KindNotExist        Kind = "NOT_EXIST"        // Item does not exist.
	KindInvalidArgument Kind = "INVALID_ARGUMENT" // Invalid argument for this type of item.
)

// app/account/account.go
package account

func getUser(userID string) (*User, err) {
    const op errors.Op = "account.getUser"
    err := loginService.Validate(userID)
    if err != nil {
        return nil, errors.E(op, err)
    }
    ...
}

// app/login/login.go
package login

func Validate(userID string) err {
    const op errors.Op = "login.Validate"
    err := db.LookUpUser(userID) // errors.New("not found")
    if err != nil {
        return nil, errors.E(op, err, KindNotExist)
    }
}

// app/main,go
package main

func main() {
    err := account.getUser("9527")

    if errors.Is(KindNotExist, err) {
        log.Errorf("not exist userid = %s, reason = %s", "9527", err)
        // Our Stack Output:
        //   not exist userid = 9527, reason = [account.GetUser] (account.go#5): [login.Validate] <NOT_EXIST> (login.go#5) not found
    }
}
```


