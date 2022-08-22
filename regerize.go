package main

import (
	"fmt"
	"os"

	"github.com/lawrencemq/regerize/parser"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Println("No file given.")
		return
	}

	filename := os.Args[1]
	regex, err := parser.ParseFile(filename)
	if err != nil {
		fmt.Println("Unable to parse file: ", err)
		return
	}
	fmt.Println(regex)

}
