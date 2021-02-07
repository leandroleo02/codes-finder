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

func matchLine(line string, word string) bool {
	return false
}

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

// FindRunes search in the file for the words in the description
func FindRunes(f *os.File, criteria string) []string {
	scanner := bufio.NewScanner(f)
	var runes []string
	for scanner.Scan() {
		line := scanner.Text()
		words := splitWords(criteria)
		for _, word := range words {
			if matchLine(line, word) {
				// TODO: format line
				runes = append(runes, line)
			}
		}
	}
	return runes
}

func listContent(f *os.File) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
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

	listContent(file)
	file.Close()
}
