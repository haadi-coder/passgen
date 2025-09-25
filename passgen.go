package passgen

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
)

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

	alphabetPower = len(upper) + len(lower) + len(digits) + len(symbols)
)

type Generator struct {
	config         *config
	availableChars []rune
}

func NewGenerator(opts ...Option) (*Generator, error) {
	gen := Generator{
		config:         defaultConfig(),
		availableChars: make([]rune, 0, alphabetPower),
	}
	for _, opt := range opts {
		opt(gen.config)
	}

	if err := gen.config.validate(); err != nil {
		return nil, fmt.Errorf("failed to create new generator: %w", err)
	}

	if gen.config.uppercase {
		gen.availableChars = append(gen.availableChars, upper...)
	}
	if gen.config.lowercase {
		gen.availableChars = append(gen.availableChars, lower...)
	}
	if gen.config.digits {
		gen.availableChars = append(gen.availableChars, digits...)
	}
	if gen.config.symbols {
		gen.availableChars = append(gen.availableChars, symbols...)
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
	cfg := g.config

	if cfg.minUppercase > 0 {
		entry, err := generatePassEntry(upper, cfg.minUppercase)
		if err != nil {
			return "", fmt.Errorf("failed to generate password entry: %w", err)
		}
		rawPass.WriteString(entry)
	}

	if cfg.minLowercase > 0 {
		entry, err := generatePassEntry(lower, cfg.minLowercase)
		if err != nil {
			return "", fmt.Errorf("failed to generate password entry: %w", err)
		}
		rawPass.WriteString(entry)
	}

	if cfg.minDigits > 0 {
		entry, err := generatePassEntry(digits, cfg.minDigits)
		if err != nil {
			return "", fmt.Errorf("failed to generate password entry: %w", err)
		}
		rawPass.WriteString(entry)
	}

	if cfg.minSymbols > 0 {
		entry, err := generatePassEntry(symbols, cfg.minSymbols)
		if err != nil {
			return "", fmt.Errorf("failed to generate password entry: %w", err)
		}
		rawPass.WriteString(entry)
	}

	remaining := cfg.length - rawPass.Len()

	if remaining > 0 {
		entry, err := generatePassEntry(g.availableChars, remaining)
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

	// Fisher–Yates shuffle
	for i := len(runes) - 1; i > 0; i-- {
		if _, err := rand.Read(buf); err != nil {
			return "", fmt.Errorf("failed to read random bytes: %w", err)
		}

		j := int(binary.LittleEndian.Uint64(buf) % uint64(i+1))
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}
