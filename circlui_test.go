package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getBranchTest(t *testing.T) {
  a := 1
  b := 2
  assert.Equal(t, 3, a + b)
}
