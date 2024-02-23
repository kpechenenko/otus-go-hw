package main

import (
	"fmt"

	"github.com/lukesiler/stringutil"
)

func main() {
	const greeting = "Hello, OTUS!"
	reversed := stringutil.Reverse(greeting)
	fmt.Print(reversed)
}
