package main

import (
	"fmt"
	"os"
)

func readFile(filename string) string {
	data, err := os.ReadFile("/tmp/dat")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func main() {
	fmt.Println("Hello, World!")
}
