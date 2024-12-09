package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	source := "Hello, OTUS!"
	fmt.Print(reverse.String(source))
}
