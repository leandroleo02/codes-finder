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

func split(words string) []string {
	splitter := func(c rune) bool {
		return c == ' ' || c == '-'
	}
	return strings.FieldsFunc(words, splitter)
}

// PrepareLine analise the line and returns the fields.
// docs: https://www.unicode.org/Public/5.1.0/ucd/UCD.html#UnicodeData.txt
func PrepareLine(line string) (rune, string, []string, error) {
	if line == "" {
		return -1, "", nil, errors.New("Empty Line")
	}

	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	name := fields[1]
	nameWords := split(fields[1])
	oldUnicodeName := fields[10]

	if oldUnicodeName != "" {
		name += fmt.Sprintf(" (%s)", oldUnicodeName)
		for _, word := range split(oldUnicodeName) {
			if !contains(nameWords, word) {
				nameWords = append(nameWords, word)
			}
		}
	}

	return rune(code), name, nameWords, nil
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
func FindRunes(r io.Reader, keyWords string) []string {
	words := split(keyWords)
	return FindRunesNew(r, words...)
}

// FindRunesNew search in the file for the words in the description
func FindRunesNew(r io.Reader, keyWords ...string) []string {
	scanner := bufio.NewScanner(r)
	var runes []string
	for scanner.Scan() {
		line := scanner.Text()
		code, name, nameWords, err := PrepareLine(line)
		if err != nil {
			log.Println(err)
			continue
		}
		if containsAll(nameWords, keyWords) {
			lineFormatted := fmt.Sprintf("U+%04X\t%[1]c\t%s", code, name)
			runes = append(runes, lineFormatted)
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

	keyWords := strings.Join(os.Args[1:], " ")
	runes := FindRunes(file, strings.ToUpper(keyWords))
	for _, r := range runes {
		fmt.Println(r)
	}
	file.Close()
}
