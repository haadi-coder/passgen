package passgen

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
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

var (
	upper = []rune{
		'A', 'B', 'C', 'D',
		'E', 'F', 'G', 'H',
		'I', 'J', 'K', 'L',
		'M', 'N', 'O', 'P',
		'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X',
		'Y', 'Z'}

	lower = []rune{
		'a', 'b', 'c', 'd',
		'e', 'f', 'g', 'h',
		'i', 'j', 'k', 'l',
		'm', 'n', 'o', 'p',
		'q', 'r', 's', 't',
		'u', 'v', 'w', 'x',
		'y', 'z'}

	digits = []rune{
		'0', '1',
		'2', '3',
		'4', '5',
		'6', '7',
		'8', '9',
	}

	symbols = []rune{
		'!', '@', '#', '$',
		'%', '^', '&', '*',
		'(', ')', '_', '+',
		'-', '=', '[', ']',
		'{', '}', '|', ';',
		':', ',', '.', '<',
		'>', '?',
	}
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
	charSets := make([]rune, 0, len(upper)+len(lower)+len(digits)+len(symbols))

	if g.uppercase {
		charSets = append(charSets, upper...)

		if g.minUppercase > 0 {
			entry, err := generatePassEntry(upper, g.minUppercase)
			if err != nil {
				return "", fmt.Errorf("failed to generate password entry: %w", err)
			}
			rawPass.WriteString(entry)
		}
	}

	if g.lowercase {
		charSets = append(charSets, lower...)

		if g.minLowercase > 0 {
			entry, err := generatePassEntry(lower, g.minLowercase)
			if err != nil {
				return "", fmt.Errorf("failed to generate password entry: %w", err)
			}
			rawPass.WriteString(entry)
		}
	}

	if g.digits {
		charSets = append(charSets, digits...)

		if g.minDigits > 0 {
			entry, err := generatePassEntry(digits, g.minDigits)
			if err != nil {
				return "", fmt.Errorf("failed to generate password entry: %w", err)
			}
			rawPass.WriteString(entry)
		}
	}

	if g.symbols {
		charSets = append(charSets, symbols...)

		if g.minSymbols > 0 {
			entry, err := generatePassEntry(symbols, g.minSymbols)
			if err != nil {
				return "", fmt.Errorf("failed to generate password entry: %w", err)
			}
			rawPass.WriteString(entry)
		}
	}

	remaining := g.length - rawPass.Len()

	if remaining > 0 {
		entry, err := generatePassEntry(charSets, remaining)
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

	alphabetLen := uint64(len(alphabet))
	buf := make([]byte, 8)

	for range count {
		if _, err := rand.Read(buf); err != nil {
			return "", fmt.Errorf("failed to read random bytes: %w", err)
		}
		idx := binary.LittleEndian.Uint64(buf) % alphabetLen
		sb.WriteRune(alphabet[idx])
	}

	return sb.String(), nil
}

func shuffleString(s string) (string, error) {
	runes := []rune(s)
	buf := make([]byte, 8)

	for i := len(runes) - 1; i > 0; i-- {
		if _, err := rand.Read(buf); err != nil {
			return "", fmt.Errorf("failed to read random bytes: %w", err)
		}

		j := int(binary.LittleEndian.Uint64(buf) % uint64(i+1))
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}
