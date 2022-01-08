package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func LoadFile(t *testing.T, path string) string {
	b, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return string(b)
}

func JS(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
