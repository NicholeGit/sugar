package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppend(t *testing.T) {
	err1 := E("err1")
	err2 := E("err2")
	mErr := Append(err1, err2)
	s := mErr.Error()
	assert.Contains(t, s, "err1")
	assert.Contains(t, s, "err2")
}
