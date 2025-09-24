package main

import (
	"fmt"

	"github.com/haadi-coder/passgen"
)

func main() {
	password, err := passgen.Generate()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("length: ", len(password))

	fmt.Println(password)
}
