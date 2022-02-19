package xrequire

import (
	"github.com/go-chai/chai/internal/tests/xassert"
	"github.com/stretchr/testify/require"
)

// JSONEq asserts that two JSON strings are equivalent.
//
//  assert.JSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func JSONEq(t require.TestingT, expected string, actual string, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if xassert.JSONEq(t, expected, actual, msgAndArgs...) {
		return
	}
	t.FailNow()
}

type tHelper interface {
	Helper()
}
