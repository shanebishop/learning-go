// Modify dup2 to print the names of all files in which each duplicated
// line occurs.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	files := os.Args[1:]

	if len(files) == 0 {
		checkForDup(os.Stdin, "stdin")
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup: %v\n", err)
				continue
			}
			checkForDup(f, arg)
			f.Close()
		}
	}
}

func checkForDup(f *os.File, filename string) {
	input := bufio.NewScanner(f)
	distinctLines := make(map[string]bool)

	for input.Scan() {
		line := input.Text()

		if found := distinctLines[line]; found {
			// We found a duplicate, so print the file
			// as a file containing duplicates
			fmt.Println(filename)
			return
		}

		distinctLines[line] = true
	}
	// NOTE: ignoring potential errors frmo input.Err()
}
