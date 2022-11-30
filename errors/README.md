# Errors

[中文文档](README_zh.md)

This library aims to be used as a drop-in replacement to github.com/pkg/errors and Go's standard errors package.

In Golang, we have the liberty to handle the errors however we want. But it comes with great responsibility since error handling is very critical to the application. Depending on the nature of our application, we can adjust the data we keep in our `Error` object and even add more data in order to improve the traceability.

## Error Struct

We will be using a custom-defined struct in order to wrap the error rather than passing the plain error.

```golang
type Operation string
type Kind string

type Error struct {
	// `Operation` is the operation being performed, usually the name of the method being invoked.
	op Operation

	// `kind` field contains the type of the error. For example, an error can be of type `NotFound`
	kind Kind

	// Human-readable message.
	message string

	// `Error` contain the error object before wrapping it using the struct.
	cause error

	// `Location` contain the error happened information. (e.g. filename, linenumber)
	location *Location
}
```
- `Location` Should not be passed. By default, the `runtime.Caller` will get the line number and file name.
- `Op` If not specified. it will be used `runtime.Caller` get the function name.


## Constructing an Error

To facilitate error construction, the package provides a function named `E`, which is short and easy to type.

    func E(args ...interface{}) error

As the doc comment for the function says, E builds an error value from its arguments. The type of each argument determines its meaning. The idea is to look at the types of the arguments and assign each argument to the field of the corresponding type in the constructed Error struct.

## Usage

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

## References

### Articles
- [practical-go: gophercon-singapore-2019#error_handling](https://dave.cheney.net/practical-go/presentations/gophercon-singapore-2019.html#_error_handling)
- [The Go Blog: Error handling and Go](https://go.dev/blog/error-handling-and-go)
- [The Go Blog: Errors are values](https://go.dev/blog/errors-are-values)
- [Exploring Error Handling Patterns in Go](https://8thlight.com/blog/kyle-krull/2018/08/13/exploring-error-handling-patterns-in-go.html)
- [Error handling in Upspin](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html)
- [pkg errors](https://github.com/pkg/errors)
- [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- [Failure is your Domain - Ben Johnson](https://middlemost.com/failure-is-your-domain/)
- [hashicorp errwrap](https://github.com/hashicorp/errwrap)

### Videos

- [GopherCon 2019: Marwan Sulaiman - Handling Go Errors](https://www.youtube.com/watch?v=4WIhhzTTd0Y)
- [GopherCon 2016: Dave Cheney - Dont Just Check Errors Handle Them Gracefully](https://www.youtube.com/watch?v=lsBF58Q-DnY)

## Contributing

If you would like to contribute, please:

- Create a GitLab issue regarding the contribution. Features and bugs should be discussed beforehand.
- Fork this repository Or can directly create a branch for this repository, e.g. `lile/custom-exit-code`.
- Create a pull request with your solution. This pull request should reference and close the issues (e.g. Fix #2).

All pull requests should:

- Be `go fmt` formatted.
