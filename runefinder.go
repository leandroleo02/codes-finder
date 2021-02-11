package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	// Ucd is the file path for unicode data table
	Ucd = "UnicodeData.txt"
)

// UnicodeData represents the unicode data from the file table 
type UnicodeData struct {
	code rune
	name string
	deprecatedUnicodeName string
}

// NewUnicodeData create new UnicodeData instance
func NewUnicodeData(code int64, name string, deprecatedUnicodeName string) UnicodeData {
	return UnicodeData {
		rune(code),
		name,
		deprecatedUnicodeName,
	}
}

func (u UnicodeData) String() string {
	newName := u.name
	if u.deprecatedUnicodeName != "" {
		newName = fmt.Sprintf("%s (%s)", u.name, u.deprecatedUnicodeName)
	}
	return fmt.Sprintf("U+%04X\t%[1]c\t%s", u.code, newName)
}

func (u UnicodeData) keyWords() []string {
	keyWords := split(u.name)
	if u.deprecatedUnicodeName != "" {
		for _, word := range split(u.deprecatedUnicodeName) {
			if !contains(keyWords, word) {
				keyWords = append(keyWords, word)
			}
		}
	}
	return keyWords
}

func split(words string) []string {
	splitter := func(c rune) bool {
		return c == ' ' || c == '-'
	}
	return strings.FieldsFunc(words, splitter)
}

// PrepareLine analise the line and returns the fields.
// docs: https://www.unicode.org/Public/5.1.0/ucd/UCD.html#UnicodeData.txt
func PrepareLine(line string) (*UnicodeData, error) {
	var unicodeData UnicodeData
	if line == "" {
		return nil, errors.New("Empty Line")
	}

	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)

	unicodeData = NewUnicodeData(code, fields[1], fields[10])
	return &unicodeData, nil
}

func contains(nameWords []string, word string) bool {
	for _, wordName := range nameWords {
		if wordName == strings.ToUpper(word) { // TODO: melhor forma
			return true
		}
	}
	return false
}

func containsAll(nameWords []string, words []string) bool {
	for _, word := range words {
		if !contains(nameWords, word) {
			return false
		}
	}
	return true
}

// FindRunes search in the file for the words in the description
func FindRunes(r io.Reader, keyWords ...string) []string {
	scanner := bufio.NewScanner(r)
	var runes []string
	for scanner.Scan() {
		line := scanner.Text()
		unicodeData, err := PrepareLine(line)
		if err != nil {
			log.Println(err)
			continue
		}
		if containsAll(unicodeData.keyWords(), keyWords) {
			runes = append(runes, unicodeData.String())
		}
	}
	return runes
}

func openUnicodeData(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func main() {
	file, err := openUnicodeData(Ucd)
	if err != nil {
		log.Fatal(err)
	}

	runes := FindRunes(file, os.Args[1:]...)
	for _, r := range runes {
		fmt.Println(r)
	}
	file.Close()
}
