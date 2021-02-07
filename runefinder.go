package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	ucd = "UnicodeData.txt"
)

func split(words string) []string {
	splitter := func(c rune) bool {
		return c == ' ' || c == '-'
	}
	return strings.FieldsFunc(words, splitter)
}

// PrepareLine analise the line and returns the fields
func PrepareLine(line string) (rune, string, []string, error) { // TODO: line should be a struct?
	if line == "" {
		return 0, "", nil, errors.New("Empty Line")
	}

	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	name := fields[1]
	nameWords := split(fields[1])

	if fields[10] != "" {
		name += fmt.Sprintf(" (%s)", fields[10])
		for _, word := range split(fields[10]) {
			if !contains(nameWords, word) {
				nameWords = append(nameWords, word)
			}
		}
	}

	return rune(code), name, nameWords, nil
}

func contains(nameWords []string, word string) bool {
	for _, wordName := range nameWords {
		if wordName == word {
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
func FindRunes(f *os.File, criteria string) []string {
	scanner := bufio.NewScanner(f)
	var runes []string
	for scanner.Scan() {
		line := scanner.Text()
		words := split(criteria)
		code, name, nameWords, err := PrepareLine(line)
		if err != nil {
			// TODO: is this ok to check empty line?
			continue
		}
		if containsAll(nameWords, words) {
			lineFormatted := fmt.Sprintf("U+%04X\t%[1]c\t%s", code, name)
			runes = append(runes, lineFormatted)
		}
	}
	return runes
}

func openUnicodeData() (*os.File, error) {
	file, err := os.Open(ucd)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func main() {
	file, err := openUnicodeData()
	if err != nil {
		log.Fatal(err)
	}

	criteria := strings.Join(os.Args[1:], " ")
	runes := FindRunes(file, strings.ToUpper(criteria))
	for _, r := range runes {
		fmt.Println(r)
	}
	file.Close()
}
