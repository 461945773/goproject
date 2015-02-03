// wordfrequency.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("usage: %s <file> [<file2>[... <fileN>]]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	frequencyForWord := map[string]int{}
	for _, filename := range commandLineFiles(os.Args[1:]) {
		updateFrequencies(filename, frequencyForWord)
	}

	// reportByWords(frequencyForWord)
	// wordsForFrequency := invertStringIntMap(frequencyForWord)
	// reportByFrequency(wordsForFrequency)
}

func commandLineFiles(files []string) []string {
	if runtime.GOOS == "windows" {
		args := make([]string, 0, len(files))
		for _, name := range files {
			if matches, err := filepath.Glob(name); err != nil {
				args = append(args, name)
			} else if matches != nil {
				args = append(args, matches...)
			}
		}
		return args
	}
	return files
}

func updateFrequencies(filename string, frequencyForWord map[string]int) {
	var file *os.File
	var err error
	if file, err = os.Open(filename); err != nil {
		log.Println("failed to open the file : ", err)
		return
	}
	defer file.Close()
	readAndUpdateFrequencies(bufio.NewReader(file), frequencyForWord)
}

func readAndUpdateFrequencies(reader *bufio.Reader, frequencyForWord map[string]int) {
	for {
		line, err := reader.ReadString('\n')
		fmt.Println(line)
		for _, word := range SplitOnNonLetters(strings.TrimSpace(line)) {
			if len(word) > utf8.UTFMax || utf8.RuneCountInString(word) > 1 {
				frequencyForWord[strings.ToLower(word)] += 1
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Println("failed to finishd to finish reading the file: ", err)
			}
			break
		}
	}
}

func SplitOnNonLetters(s string) []string {
	notAletters := func(char rune) bool { return !unicode.IsLetter(char) }
	return strings.FieldsFunc(s, notAletters)
}
