package openapi2

import (
	"testing"

	"github.com/go-chai/chai/examples/shared/model"
	"github.com/stretchr/testify/assert"
)

func TestSpecGen(t *testing.T) {
	b := new(model.Bottle)
	_, err := SpecGen(b)
	assert.NoError(t, err)
}

func TestSpecGen2(t *testing.T) {
	b := new(model.Bottle)
	_, err := SpecGen2(b)
	assert.NoError(t, err)
}
