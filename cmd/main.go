package main

import (
	"fmt"

	"github.com/haadi-coder/passgen"
)

func main() {
	pinGenerator, _ := NewPINGenerator()

	fmt.Println(pinGenerator.gen4.Generate())
	fmt.Println(pinGenerator.gen6.Generate())
}

// Генератор PIN-кодов разной длины
type PINGenerator struct {
	gen4 *passgen.Generator
	gen6 *passgen.Generator
}

func NewPINGenerator() (*PINGenerator, error) {
	gen4, err := passgen.NewGenerator(
		passgen.WithLength(4),
		passgen.WithoutUppercase(),
		passgen.WithoutLowercase(),
		passgen.WithoutSymbols(),
	)
	if err != nil {
		return nil, err
	}

	gen6, err := passgen.NewGenerator(
		passgen.WithLength(6),
		passgen.WithoutUppercase(),
		passgen.WithoutLowercase(),
		passgen.WithoutSymbols(),
	)
	if err != nil {
		return nil, err
	}

	return &PINGenerator{gen4: gen4, gen6: gen6}, nil
}
