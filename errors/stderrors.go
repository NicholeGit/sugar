package errors

import stderrors "errors"

// Is is a proxy for the Is function in Go's standard `errors` library
// (pkg.go.dev/errors).
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// As is a proxy for the As function in Go's standard `errors` library
// (pkg.go.dev/errors).
func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

// Unwrap is a proxy for the Unwrap function in Go's standard `errors` library
// (pkg.go.dev/errors).
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}
