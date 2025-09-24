package passgen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

type Generator struct {
	length    int
	uppercase bool
	lowercase bool
	digits    bool
	symbols   bool

	minUppercase int
	minLowercase int
	minDigits    int
	minSymbols   int
}

const (
	upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower   = "abcdefghijklmnopqrstuvwxyz"
	digits  = "0123456789"
	symbols = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

func NewGenerator(opts ...Option) (*Generator, error) {
	gen := Generator{
		length:    16,
		uppercase: true,
		lowercase: true,
		digits:    true,
		symbols:   true,
	}
	for _, opt := range opts {
		opt(&gen)
	}

	if err := gen.validate(); err != nil {
		return nil, fmt.Errorf("failed to create new generator: %w", err)
	}

	return &gen, nil
}

func Generate(opts ...Option) (string, error) {
	gen, err := NewGenerator(opts...)
	if err != nil {
		return "", fmt.Errorf("failed to generate: %w", err)
	}

	return gen.Generate()
}

func (g *Generator) Generate() (string, error) {
	var rawPass strings.Builder
	var charSets strings.Builder

	if g.uppercase {
		charSets.WriteString(upper)

		if g.minUppercase > 0 {
			entry, err := generatePassEntry([]rune(upper), g.minUppercase)
			if err != nil {
				return "", fmt.Errorf("failed to generate password entry: %w", err)
			}
			rawPass.WriteString(entry)
		}
	}

	if g.lowercase {
		charSets.WriteString(lower)

		if g.minLowercase > 0 {
			entry, err := generatePassEntry([]rune(lower), g.minLowercase)
			if err != nil {
				return "", fmt.Errorf("failed to generate password entry: %w", err)
			}
			rawPass.WriteString(entry)
		}
	}

	if g.digits {
		charSets.WriteString(digits)

		if g.minDigits > 0 {
			entry, err := generatePassEntry([]rune(digits), g.minDigits)
			if err != nil {
				return "", fmt.Errorf("failed to generate password entry: %w", err)
			}
			rawPass.WriteString(entry)
		}
	}

	if g.symbols {
		charSets.WriteString(symbols)

		if g.minSymbols > 0 {
			entry, err := generatePassEntry([]rune(symbols), g.minSymbols)
			if err != nil {
				return "", fmt.Errorf("failed to generate password entry: %w", err)
			}
			rawPass.WriteString(entry)
		}
	}

	remaining := g.length - rawPass.Len()

	if remaining > 0 {
		entry, err := generatePassEntry([]rune(charSets.String()), remaining)
		if err != nil {
			return "", fmt.Errorf("failed to generate password entry: %w", err)
		}
		rawPass.WriteString(entry)
	}

	pass, err := shuffleString(rawPass.String())
	if err != nil {
		return "", fmt.Errorf("failed to shuffle password: %w", err)
	}

	return pass, nil
}

func generatePassEntry(alphabet []rune, count int) (string, error) {
	var sb strings.Builder

	infig := big.NewInt(int64(len(alphabet)))

	for range int(count) {
		idx, err := rand.Int(rand.Reader, infig)
		if err != nil {
			return "", fmt.Errorf("failed to generate random index: %w", err)
		}
		sb.WriteRune(alphabet[idx.Int64()])
	}

	return sb.String(), nil
}

func shuffleString(s string) (string, error) {
	runes := []rune(s)

	for i := range runes {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", fmt.Errorf("failed to generate random index to shuffle: %w", err)
		}

		runes[i], runes[j.Int64()] = runes[j.Int64()], runes[i]
	}

	return string(runes), nil
}
