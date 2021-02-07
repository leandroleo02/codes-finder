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

func TestFindRunes(t *testing.T) {
	file, _ := openUnicodeData()

	runes := FindRunes(file, "CHIPMUNK")
	assert.Len(t, runes, 1)
}

func TestPrepareLineIgnoreEmptiness(t *testing.T) {
	code, name, words, err := PrepareLine("")

	assert.Equal(t, code, int32(0))
	assert.Equal(t, name, "")
	assert.Len(t, words, 0)
	assert.Error(t, err)
}

func TestPrepareLineWithSingleWordDescription(t *testing.T) {
	code, name, words, err := PrepareLine("1F43F;CHIPMUNK;So;0;ON;;;;;N;;;;;")

	assert.NoError(t, err)
	assert.Equal(t, code, rune(128063))
	assert.Equal(t, name, "CHIPMUNK")
	assert.Len(t, words, 1)
}

func TestPrepareLineWithMultipleWordDescriptionSeparetedOnlyBySpace(t *testing.T) {
	code, name, words, err := PrepareLine("1F601;GRINNING FACE WITH SMILING EYES;So;0;ON;;;;;N;;;;;")

	assert.NoError(t, err)
	assert.Equal(t, code, rune(128513))
	assert.Equal(t, name, "GRINNING FACE WITH SMILING EYES")
	assert.Len(t, words, 5)
}
