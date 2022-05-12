// Modify the echo program introduced in section 1.3 of the book
// to also print os.Args[0], the name of the command that invoked
// it.

// echo program to modify:
// func main() {
//	fmt.Println(strings.Join(os.Args[1:], " "))
// }

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args, " "))
}
