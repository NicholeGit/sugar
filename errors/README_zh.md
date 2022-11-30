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

## 参考文献

### 文章
- [practical-go: gophercon-singapore-2019#error_handling](https://dave.cheney.net/practical-go/presentations/gophercon-singapore-2019.html#_error_handling)
- [The Go Blog: Error handling and Go](https://go.dev/blog/error-handling-and-go)
- [The Go Blog: Errors are values](https://go.dev/blog/errors-are-values)
- [Exploring Error Handling Patterns in Go](https://8thlight.com/blog/kyle-krull/2018/08/13/exploring-error-handling-patterns-in-go.html)
- [Error handling in Upspin](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html)
- [pkg errors](https://github.com/pkg/errors)
- [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- [Failure is your Domain - Ben Johnson](https://middlemost.com/failure-is-your-domain/)
- [hashicorp errwrap](https://github.com/hashicorp/errwrap)

### 视频

- [GopherCon 2019: Marwan Sulaiman - Handling Go Errors](https://www.youtube.com/watch?v=4WIhhzTTd0Y)
- [GopherCon 2016: Dave Cheney - Dont Just Check Errors Handle Them Gracefully](https://www.youtube.com/watch?v=lsBF58Q-DnY)

## 贡献

如果你想要贡献代码, 请参考:

- 如果是增加新的功能和 bug 修复，请先创建一个 issue 并做简单描述以及大致的实现方法.
- 不要 Fork 这个仓库，可以直接提交自己的开发分支, 例如: `lile/custom-exit-code`.
- 创建一个 PR 关于这个解决方案，并且引用相关的 issue (例如: `Fix #2`).

所有的 PR 应该:

- 使用 `go fmt` 进行代码格式化.