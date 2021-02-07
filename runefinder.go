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

func splitWords(words string) []string {
	return strings.Split(words, " ")
}

// PrepareLine analise the line and returns the fields
func PrepareLine(line string) (rune, string, []string, error) { // TODO: line should be a struct?
	if line == "" {
		return 0, "", nil, errors.New("Empty Line")
	}

	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	name := fields[1]
	words := splitWords(fields[1])

	return rune(code), name, words, nil
}

func matchLine(wordNames []string, word string) bool {
	for _, wordName := range wordNames {
		if wordName == word {
			return true
		}
	}
	return false
}

// FindRunes search in the file for the words in the description
func FindRunes(f *os.File, criteria string) []string {
	scanner := bufio.NewScanner(f)
	var runes []string
	for scanner.Scan() {
		line := scanner.Text()
		words := splitWords(criteria)
		code, name, wordNames, err := PrepareLine(line)
		if err != nil {
			// TODO: is this ok to check empty line?
			continue
		}
		for _, word := range words {
			if matchLine(wordNames, word) {
				lineFormatted := fmt.Sprintf("U+%04X\t%[1]c\t%s", code, name)
				runes = append(runes, lineFormatted)
			}
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

	runes := FindRunes(file, strings.ToUpper(strings.Join(os.Args[1:], " ")))
	for _, r := range runes {
		fmt.Println(r)
	}
	file.Close()
}
