package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	file, err := openUnicodeData()

	assert.NoError(t, err)
	assert.NotNil(t, file)
}
