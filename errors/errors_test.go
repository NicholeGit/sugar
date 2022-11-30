package errors

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	KindInternal        Kind = "INTERNAL"         // Internal error or inconsistency.
	KindNotExist        Kind = "NOT_EXIST"        // Item does not exist.
	KindInvalidArgument Kind = "INVALID_ARGUMENT" // Invalid argument for this type of item.
)

// Line number is #20.
func funcForAutoGenerateOpTest() error {
	return E(KindInvalidArgument, "invalid arguments")
}

var (
	ErrTest = errors.New("a test error")
)

func TestError_Error(t *testing.T) {
	tests := map[string]struct {
		givenError error
		wantStr    string
	}{
		"simple": {
			E(Operation("DelUser"), "user joe not found"),
			"[DelUser] (errors_test.go#33) user joe not found",
		},
		"simple with kind with Msg": {
			E(Operation("DelUser"), KindNotExist, "user joe not found"),
			"[DelUser] <NOT_EXIST> (errors_test.go#37) user joe not found",
		},
		"simple with kind without Msg": {
			E(Operation("DelUser"), KindNotExist),
			"[DelUser] <NOT_EXIST> (errors_test.go#41)",
		},
		"simple with auto Op and location without Msg": {
			funcForAutoGenerateOpTest(),
			"[errors.funcForAutoGenerateOpTest] <INVALID_ARGUMENT> (errors_test.go#20) invalid arguments",
		},
		"wrap external error with print Msg": {
			E(Operation("DelUser"), KindNotExist, "user joe not found", ErrTest),
			"[DelUser] <NOT_EXIST> (errors_test.go#49) user joe not found a test error",
		},
		"wrap external error": {
			E(Operation("DelUser"), KindNotExist, ErrTest),
			"[DelUser] <NOT_EXIST> (errors_test.go#53) a test error",
		},
		"wrap *Error": {
			E(
				Operation("HandleDelUser"),
				E(Operation("DelUser"), KindNotExist, "user joe not found"),
			),
			"[HandleDelUser] (errors_test.go#57): [DelUser] <NOT_EXIST> (errors_test.go#59) user joe not found",
		},
		"wrap *Error and external error": {
			E(
				Operation("HandleDelUser"), "user joe not found",
				E(Operation("DelUser"), KindNotExist, ErrTest),
			),
			"[HandleDelUser] (errors_test.go#64) user joe not found: [DelUser] <NOT_EXIST> (errors_test.go#66) a test error",
		},
		"wrap errors.New": {
			E(Operation("DelUser"), KindNotExist, errors.New("user joe not found")),
			"[DelUser] <NOT_EXIST> (errors_test.go#71) user joe not found",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantStr, tc.givenError.Error())
		})
	}
}

func TestE(t *testing.T) {
	tests := map[string]struct {
		givenKind  Kind
		givenMsg   string
		givenOp    Operation
		givenCause error
		wantError  error
	}{
		"simple": {
			KindInternal, "data inconsistent", Operation("SetUser"), nil,
			&Error{
				kind:    KindInternal,
				message: "data inconsistent",
				op:      Operation("SetUser"),
				cause:   nil,
				location: &Location{
					filename:   "errors_test.go",
					lineNumber: 138,
				},
			},
		},
		"nested": {
			KindNotExist, "user joe not found", Operation("GetUser"), ErrTest,
			&Error{
				kind:    KindNotExist,
				message: "user joe not found",
				op:      Operation("GetUser"),
				cause:   ErrTest,
				location: &Location{
					filename:   "errors_test.go",
					lineNumber: 138,
				},
			},
		},
		"msg only": {
			Kind(""), "user joe not found", Operation("GetUser"), nil,
			&Error{
				message: "user joe not found",
				op:      Operation("GetUser"),
				location: &Location{
					filename:   "errors_test.go",
					lineNumber: 138,
				},
			},
		},
		"dedup": {
			KindInternal, "data inconsistent", Operation("SetUser"), &Error{op: Operation("SetUser")},
			&Error{
				op: Operation("SetUser"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantError, E(tc.givenOp, tc.givenKind, tc.givenMsg, tc.givenCause))
		})
	}

	t.Run("panic with no valid arguments", func(t *testing.T) {
		assert.Panics(t, func() { E() })
	})

	t.Run("invalid type", func(t *testing.T) {
		assert.True(t, strings.Contains(E(bytes.NewBuffer(nil)).Error(), "unknown type"))
	})

	t.Run("auto-generate op", func(t *testing.T) {
		f1 := func() {
			err := funcForAutoGenerateOpTest()
			errV := err.(*Error)
			assert.Equal(t, Operation("errors.funcForAutoGenerateOpTest"), errV.op)
		}
		f1()
	})
}

func TestWrap(t *testing.T) {
	tests := map[string]struct {
		givenKind  Kind
		givenMsg   string
		givenOp    Operation
		givenCause error
		wantError  error
	}{
		"return nil": {
			KindInternal, "data inconsistent", Operation("SetUser"), nil,
			nil,
		},
		"simple": {
			KindNotExist, "user joe not found", Operation("GetUser"), ErrTest,
			&Error{
				kind:    KindNotExist,
				message: "user joe not found",
				op:      Operation("GetUser"),
				cause:   ErrTest,
				location: &Location{
					filename:   "errors_test.go",
					lineNumber: 189,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantError, Wrap(tc.givenCause, tc.givenOp, tc.givenKind, tc.givenMsg))
		})
	}
}

func TestGenOp(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		op := genOp("")
		assert.Equal(t, Operation(""), op)
	})

	t.Run("main.Func", func(t *testing.T) {
		op := genOp("main.Func")
		assert.Equal(t, Operation("main.Func"), op)
	})

	t.Run("package.struct.Func", func(t *testing.T) {
		op := genOp("package.struct.Func")
		assert.Equal(t, Operation("struct.Func"), op)
	})

	t.Run("package.(*struct).Func", func(t *testing.T) {
		op := genOp("module-name/package.(*struct).Func")
		assert.Equal(t, Operation("struct.Func"), op)
	})

	t.Run("case sensitive", func(t *testing.T) {
		op := genOp("module-name/package.(*stRUct).Func")
		assert.NotEqual(t, Operation("struct.Func"), op)
		assert.Equal(t, Operation("stRUct.Func"), op)
	})

	t.Run("more than one package", func(t *testing.T) {
		op := genOp("module-name/pkg-0/pkg_1.(*struct).Func")
		assert.Equal(t, Operation("struct.Func"), op)
	})

	t.Run("more symbols", func(t *testing.T) {
		op := genOp("git.shiyou.kingsoft.com/module-name///pkg_1.(*struct).Func")
		assert.Equal(t, Operation("struct.Func"), op)
	})
}

func TestMessage(t *testing.T) {
	tests := map[string]struct {
		givenError error
		wantErrMsg string
	}{
		"nil": {
			nil,
			"",
		},
		"not of type *Error": {
			ErrTest,
			"",
		},
		"simple": {
			E(Operation("DelUser"), KindNotExist, "user joe not found"),
			"user joe not found",
		},
		"simple code only": {
			E(Operation("DelUser"), KindNotExist),
			"",
		},
		"simple msg only": {
			E(Operation("DelUser"), KindNotExist, "user joe not found"),
			"user joe not found",
		},
		"wrapped": {
			E(
				Operation("HandleDelUser"), KindInternal, "userService error",
				E(Operation("DelUser"), KindNotExist, "user joe not found"),
			),
			"userService error",
		},
		"wrapped and first without msg": {
			E(
				Operation("HandleDelUser"), KindInternal,
				E(Operation("DelUser"), KindNotExist, "user joe not found"),
			),
			"user joe not found",
		},
		"without msg, code and cause": {
			E(
				Operation("HandleDelUser"), KindInternal,
				E(Operation("DelUser"), KindNotExist),
			),
			"",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantErrMsg, GetMessage(tc.givenError))
		})
	}
}

func TestMatch(t *testing.T) {
	tests := map[string]struct {
		givenErr  error
		givenKind Kind
		wantIs    bool
	}{
		"nil error": {
			nil, KindInternal, false,
		},
		"simple is": {
			E(Operation("DelUser"), KindInternal, "user joe not found"),
			KindInternal, true,
		},
		"simple is not": {
			E(Operation("DelUser"), KindInternal, "user joe not found"),
			KindInvalidArgument, false,
		},
		"types other than *Error": {
			errors.New("invalid username"), KindInvalidArgument, false,
		},
		"wrapped is": {
			E(
				Operation("HandleDelUser"), "userService error",
				E(Operation("DelUser"), KindNotExist, "user joe not found"),
			),
			KindNotExist, true,
		},
		"wrapped is not": {
			E(
				Operation("HandleDelUser"), "userService error",
				E(Operation("DelUser"), KindNotExist, "user joe not found"),
			),
			KindInvalidArgument, false,
		},
		"first non-empty is not": {
			E(
				Operation("HandleDelUser"), KindInternal, "userService error",
				E(Operation("DelUser"), KindNotExist, "user joe not found"),
			),
			KindNotExist, false,
		},
		"is empty": {
			E(Operation("HandleDelUser"), "userService error"),
			"", false,
		},
		"wrapped is empty": {
			E(
				Operation("HandleDelUser"), "userService error",
				E(Operation("DelUser"), "user joe not found"),
			),
			"", false,
		},
		"is not empty": {
			E(
				Operation("HandleDelUser"), "userService error",
				E(Operation("DelUser"), KindInternal, "user joe not found"),
			),
			"", false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantIs, Match(tc.givenKind, tc.givenErr))
		})
	}
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		givenText string
		wantError error
	}{
		"name": {
			"message",
			&Error{
				message: "message",
				op:      "TestNew.func1",
				location: &Location{
					filename:   "errors_test.go",
					lineNumber: 375,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantError, New(tc.givenText))
		})
	}
}

func TestIs(t *testing.T) {
	tests := map[string]struct {
		givenErr  error
		targetErr error
		wantIs    bool
	}{
		"normal error": {
			E(ErrTest), ErrTest, true,
		},
		"no match error": {
			E(ErrTest), errors.New("unexpected error"), false,
		},
		"wrap error": {
			E("wrap error", E(ErrTest)), errors.New("unexpected error"), false,
		},
		"nil error": {
			E(ErrTest), nil, false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantIs, Is(tc.givenErr, tc.targetErr))
		})
	}
}
