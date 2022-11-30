package errors

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
)

type Operation string

type Kind string

// Separator is the string used to separate nested errors
const Separator = ":"

// Location is the where the error happened.
type Location struct {
	filename   string
	lineNumber int
}

func (l *Location) String() string {
	// (e.g.) `main.go#3`
	return fmt.Sprintf("%s#%d", l.filename, l.lineNumber)
}

// Error is the type that implements the error interface.
// Error defines a standard application error.
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

// E builds an error value from its arguments.
// There must be at least one argument.

// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.

// If Op is not specified, we set it to use function name by `runtime.Caller`.
// Note: `E` must be return a non-nil error object.
func E(args ...interface{}) error {
	return eWithSkip(2, args...)
}

// Similar to errors.Wrap, if `error` passing nil, it returns nil and does not create a new `Error` struct.
// Use cases for this function:
// - Reduce code redundancy like `if err != nil {}`, You can just write code like `return errors.Wrap(err, ...)`
func Wrap(err error, args ...interface{}) error {
	if err == nil {
		return nil
	}
	var a = append(args, err)
	return eWithSkip(2, a...)
}

// New is simlar to errors.New in the standard package.
func New(text string) error {
	return eWithSkip(2, text)
}

func eWithSkip(skip int, args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}

	e := &Error{}

	var hasOp bool = false

	for _, arg := range args {
		if arg == nil {
			continue
		}
		switch arg := arg.(type) {
		case Operation:
			e.op = arg
			hasOp = true
		case Kind:
			e.kind = arg
		case string:
			e.message = arg
		case *Error:
			// Make a copy
			c := *arg
			e.cause = &c
		case error:
			e.cause = arg
		default:
			_, file, line, _ := runtime.Caller(skip)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			return fmt.Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}

	// Create Location by `runtime.Caller`
	pc, filename, line, _ := runtime.Caller(skip)
	e.location = &Location{
		filename:   lastName(filename),
		lineNumber: line,
	}

	// If `hasOp` is false, it indicates that use function name by `runtime.Caller`.
	if !hasOp {
		funcName := runtime.FuncForPC(pc).Name()
		e.op = genOp(funcName)
	}

	// deduplication
	if cause, ok := e.cause.(*Error); ok && cause.op == e.op {
		return cause
	}

	return e
}

// git.shiyou.kingsoft.com/errors.E => errors.E
func lastName(fullname string) string {
	parts := strings.Split(fullname, "/")
	return parts[len(parts)-1]
}

// (e.g.) `module-name/pkg.(*T).Test` => `T.Test`
func genOp(funcName string) Operation {
	re := regexp.MustCompile(`.*[./]\(?\*?(\w+)\)?(\.\w+)$`)
	return Operation(re.ReplaceAllString(funcName, "$1$2"))
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *Error) Unwrap() error { return w.cause }

// Error returns the string representation of the error message.
// Note: Refer to the unit tests for more detailed output.
func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	fmt.Fprintf(&buf, "[%s]", e.op)

	if e.kind != "" {
		fmt.Fprintf(&buf, " <%s>", e.kind)
	}

	if e.location != nil {
		fmt.Fprintf(&buf, " (%s)", e.location.String())
	}

	if e.message != "" {
		fmt.Fprintf(&buf, " %s", e.message)
	}

	if e.cause != nil {
		// Separate nested errors by the variable `Separator`
		if _, ok := e.cause.(*Error); ok {
			buf.WriteString(Separator)
			fmt.Fprintf(&buf, " %s", e.cause.Error())
		} else {
			fmt.Fprintf(&buf, " %s", e.cause.Error())
		}
	}

	return buf.String()
}

// Ops returns the "stack" of operations
// for each generated error.
func GetOps(e *Error) []Operation {
	res := []Operation{e.op}

	subErr, ok := e.cause.(*Error)
	if !ok {
		return res
	}

	res = append(res, GetOps(subErr)...)

	return res
}

func GetKind(err error) Kind {
	e, ok := err.(*Error)
	if !ok {
		return ""
	}

	if e.kind != "" {
		return e.kind
	}

	return GetKind(e.cause)
}

// GetMessage returns the first human-readable message of the error, if available.
func GetMessage(err error) string {
	e, ok := err.(*Error)
	if !ok {
		return ""
	}
	if e.message != "" {
		return e.message
	}
	if e.cause != nil {
		return GetMessage(e.cause)
	}
	return ""
}

// Notes: `Is` be replaced by `Match` since v1.0.3

// `Match` reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
// If kind of err is empty then Is returns false.
func Match(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.kind != "" {
		return e.kind == kind
	}
	if e.cause != nil {
		return Match(kind, e.cause)
	}
	return false
}
