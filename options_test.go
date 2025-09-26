package passgen

import (
	"regexp"
	"testing"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		name            string
		options         []Option
		expectUppercase bool
		expectLowercase bool
		expectDigits    bool
		expectSymbols   bool
		expectedLength  int
		minUppercase    int
		minLowercase    int
		minDigits       int
		minSymbols      int
	}{
		{
			name:            "WithUpperCase",
			options:         []Option{WithoutUppercase(), WithUppercase()},
			expectUppercase: true,
		},
		{
			name:            "WithoutUppercase",
			options:         []Option{WithoutUppercase()},
			expectUppercase: false,
		},
		{
			name:            "WithLowercase",
			options:         []Option{WithoutLowercase(), WithLowercase()},
			expectLowercase: true,
		},
		{
			name:            "WithoutLowercase",
			options:         []Option{WithoutLowercase()},
			expectLowercase: false,
		},
		{
			name:         "WithDigits",
			options:      []Option{WithoutDigits(), WithDigits()},
			expectDigits: true,
		},
		{
			name:         "WithoutDigits",
			options:      []Option{WithoutDigits()},
			expectDigits: false,
		},
		{
			name:          "WithSymbols",
			options:       []Option{WithoutSymbols(), WithSymbols()},
			expectSymbols: true,
		},
		{
			name:          "WithoutSymbols",
			options:       []Option{WithoutSymbols()},
			expectSymbols: false,
		},
		{
			name:           "WithLength",
			options:        []Option{WithLength(25)},
			expectedLength: 25,
		},
		{
			name:         "WithMinUppercase",
			options:      []Option{WithMinUppercase(5)},
			minUppercase: 5,
		},
		{
			name:         "WithMinLowercase",
			options:      []Option{WithMinLowercase(4)},
			minLowercase: 4,
		},
		{
			name:      "WithMinDigits",
			options:   []Option{WithMinDigits(3)},
			minDigits: 3,
		},
		{
			name:       "WithMinSymbols",
			options:    []Option{WithMinSymbols(2)},
			minSymbols: 2,
		},
		{
			name:         "WithMinRequirements",
			options:      []Option{WithMinRequirements(2, 3, 4, 1)},
			minUppercase: 2,
			minLowercase: 3,
			minDigits:    4,
			minSymbols:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allOptions := append([]Option{WithLength(20)}, tt.options...)
			password, err := Generate(allOptions...)
			if err != nil {
				t.Fatalf("failed to generate password: %v", err)
			}

			if tt.expectedLength > 0 && len(password) != tt.expectedLength {
				t.Errorf("expected length %d, got %d", tt.expectedLength, len(password))
			}

			if tt.minUppercase > 0 {
				count := len(regexp.MustCompile(`[A-Z]`).FindAllString(password, -1))
				if count < tt.minUppercase {
					t.Errorf("expected at least %d uppercase characters, got %d in password: %s", tt.minUppercase, count, password)
				}
			}

			if tt.minLowercase > 0 {
				count := len(regexp.MustCompile(`[a-z]`).FindAllString(password, -1))
				if count < tt.minLowercase {
					t.Errorf("expected at least %d lowercase characters, got %d in password: %s", tt.minLowercase, count, password)
				}
			}

			if tt.minDigits > 0 {
				count := len(regexp.MustCompile(`[0-9]`).FindAllString(password, -1))
				if count < tt.minDigits {
					t.Errorf("expected at least %d digits, got %d in password: %s", tt.minDigits, count, password)
				}
			}

			if tt.minSymbols > 0 {
				count := len(regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`).FindAllString(password, -1))
				if count < tt.minSymbols {
					t.Errorf("expected at least %d symbols, got %d in password: %s", tt.minSymbols, count, password)
				}
			}
		})
	}
}
