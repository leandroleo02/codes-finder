package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func unicodeDataFixture() io.Reader {
	var buf bytes.Buffer

	buf.WriteString("1DA19;SIGNWRITING EYES HALF OPEN;Mn;0;NSM;;;;;N;;;;;\n")
	buf.WriteString("1DA1A;SIGNWRITING EYES WIDE OPEN;Mn;0;NSM;;;;;N;;;;;\n")
	buf.WriteString("1DA1B;SIGNWRITING EYES HALF CLOSED;Mn;0;NSM;;;;;N;;;;;\n")
	buf.WriteString("1DA1C;SIGNWRITING EYES WIDENING MOVEMENT;Mn;0;NSM;;;;;N;;;;;\n")
	buf.WriteString("1F43F;CHIPMUNK;So;0;ON;;;;;N;;;;;\n")
	buf.WriteString("1F601;GRINNING FACE WITH SMILING EYES;So;0;ON;;;;;N;;;;;\n")
	buf.WriteString("1F604;SMILING FACE WITH OPEN MOUTH AND SMILING EYES;So;0;ON;;;;;N;;;;;\n")
	buf.WriteString("20D7;COMBINING RIGHT ARROW ABOVE;Mn;230;NSM;;;;;N;NON-SPACING RIGHT ARROW ABOVE;;;;\n")

	return strings.NewReader(buf.String())
}

func TestReadFile(t *testing.T) {
	file, err := openUnicodeData(Ucd)

	assert.NoError(t, err)
	assert.NotNil(t, file)
}

func TestErrorReadFile(t *testing.T) {
	_, err := openUnicodeData("UnicodeDataDoesNotExist.txt")

	assert.Error(t, err)
}

func TestFindRunesEmptyKeyWord(t *testing.T) {
	r := unicodeDataFixture()

	runes := FindRunes(r, "")
	assert.Len(t, runes, 8)
}

func TestFindRunesWithOneKeyWord(t *testing.T) {
	r := unicodeDataFixture()

	runes := FindRunes(r, "CHIPMUNK")
	assert.Len(t, runes, 1)
}

func TestFindRunesWithMoreThanOneWord(t *testing.T) {
	r := unicodeDataFixture()

	runes := FindRunes(r, "FACE EYES")
	assert.Len(t, runes, 2)
}

func TestFindRunesUsingSecondNameColumn(t *testing.T) {
	r := unicodeDataFixture()

	runes := FindRunes(r, "NON SPACING RIGHT ARROW ABOVE")
	assert.Len(t, runes, 1)
}

func TestFindRunesUsingKeyWordsArray(t *testing.T) {
	r := unicodeDataFixture()

	runes := FindRunesNew(r, "NON", "SPACING", "RIGHT", "ARROW", "ABOVE")
	assert.Len(t, runes, 1)
}

func TestPrepareLineIgnoreEmptiness(t *testing.T) {
	code, name, words, err := PrepareLine("")

	assert.Equal(t, code, int32(-1))
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
