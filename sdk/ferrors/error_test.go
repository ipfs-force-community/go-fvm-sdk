// Package ferrors fvm errors
package ferrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSysCallError(t *testing.T) {
	err := NewSysCallError(6, "this is error:")
	assert.True(t, errors.Is(err, NotFound))
}
