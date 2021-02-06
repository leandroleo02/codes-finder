package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	ucd = "UnicodeData.txt"
)

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
	file, err := openUnicodeData();
	if err != nil {
		log.Fatal(err)
	}
	
	listContent(file)
	file.Close()
}
