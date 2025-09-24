package passgen

import (
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name        string
		options     []Option
		expectError bool
		errorMsg    string
	}{
		{
			name:        "default configuration",
			options:     nil,
			expectError: false,
		},
		{
			name:        "valid length",
			options:     []Option{WithLength(20)},
			expectError: false,
		},
		{
			name:        "zero length",
			options:     []Option{WithLength(0)},
			expectError: true,
			errorMsg:    "password length must be greater than 0",
		},
		{
			name:        "negative length",
			options:     []Option{WithLength(-5)},
			expectError: true,
			errorMsg:    "password length must be greater than 0",
		},
		{
			name:        "too long password",
			options:     []Option{WithLength(10001)},
			expectError: true,
			errorMsg:    "password length must not exceed 10000",
		},
		{
			name:        "all character sets disabled",
			options:     []Option{WithoutUppercase(), WithoutLowercase(), WithoutDigits(), WithoutSymbols()},
			expectError: true,
			errorMsg:    "at least one character set must be enabled",
		},
		{
			name:        "negative min uppercase",
			options:     []Option{WithMinUppercase(-1)},
			expectError: true,
			errorMsg:    "minimum uppercase count cannot be negative",
		},
		{
			name:        "negative min lowercase",
			options:     []Option{WithMinLowercase(-1)},
			expectError: true,
			errorMsg:    "minimum lowercase count cannot be negative",
		},
		{
			name:        "negative min digits",
			options:     []Option{WithMinDigits(-1)},
			expectError: true,
			errorMsg:    "minimum digits count cannot be negative",
		},
		{
			name:        "negative min symbols",
			options:     []Option{WithMinSymbols(-1)},
			expectError: true,
			errorMsg:    "minimum symbols count cannot be negative",
		},
		{
			name:        "min requirements exceed length",
			options:     []Option{WithLength(10), WithMinUppercase(3), WithMinLowercase(3), WithMinDigits(3), WithMinSymbols(3)},
			expectError: true,
			errorMsg:    "sum of minimum requirements (12) cannot exceed password length (10)",
		},
		{
			name:        "min uppercase with disabled uppercase",
			options:     []Option{WithoutUppercase(), WithMinUppercase(2)},
			expectError: true,
			errorMsg:    "uppercase characters are disabled but minimum uppercase requirement is 2",
		},
		{
			name:        "min lowercase with disabled lowercase",
			options:     []Option{WithoutLowercase(), WithMinLowercase(2)},
			expectError: true,
			errorMsg:    "lowercase characters are disabled but minimum lowercase requirement is 2",
		},
		{
			name:        "min digits with disabled digits",
			options:     []Option{WithoutDigits(), WithMinDigits(2)},
			expectError: true,
			errorMsg:    "digits are disabled but minimum digits requirement is 2",
		},
		{
			name:        "min symbols with disabled symbols",
			options:     []Option{WithoutSymbols(), WithMinSymbols(2)},
			expectError: true,
			errorMsg:    "symbols are disabled but minimum symbols requirement is 2",
		},
		{
			name:        "conflicting options - last wins",
			options:     []Option{WithoutUppercase(), WithUpperCase()},
			expectError: false,
		},
		{
			name:        "valid min requirements",
			options:     []Option{WithLength(20), WithMinRequirements(2, 2, 2, 2)},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewGenerator(tt.options...)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if gen == nil {
					t.Error("expected non-nil generator")
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name            string
		options         []Option
		expectedLength  int
		expectUppercase bool
		expectLowercase bool
		expectDigits    bool
		expectSymbols   bool
		minUppercase    int
		minLowercase    int
		minDigits       int
		minSymbols      int
		expectError     bool
		errorMsg        string
	}{
		{
			name:            "default generation",
			options:         nil,
			expectedLength:  16,
			expectUppercase: true,
			expectLowercase: true,
			expectDigits:    true,
			expectSymbols:   true,
		},
		{
			name:            "custom length",
			options:         []Option{WithLength(24)},
			expectedLength:  24,
			expectUppercase: true,
			expectLowercase: true,
			expectDigits:    true,
			expectSymbols:   true,
		},
		{
			name:            "only uppercase",
			options:         []Option{WithLength(12), WithoutLowercase(), WithoutDigits(), WithoutSymbols()},
			expectedLength:  12,
			expectUppercase: true,
			expectLowercase: false,
			expectDigits:    false,
			expectSymbols:   false,
		},
		{
			name:            "only lowercase",
			options:         []Option{WithLength(12), WithoutUppercase(), WithoutDigits(), WithoutSymbols()},
			expectedLength:  12,
			expectUppercase: false,
			expectLowercase: true,
			expectDigits:    false,
			expectSymbols:   false,
		},
		{
			name:            "only digits",
			options:         []Option{WithLength(6), WithoutUppercase(), WithoutLowercase(), WithoutSymbols()},
			expectedLength:  6,
			expectUppercase: false,
			expectLowercase: false,
			expectDigits:    true,
			expectSymbols:   false,
		},
		{
			name:            "only symbols",
			options:         []Option{WithLength(8), WithoutUppercase(), WithoutLowercase(), WithoutDigits()},
			expectedLength:  8,
			expectUppercase: false,
			expectLowercase: false,
			expectDigits:    false,
			expectSymbols:   true,
		},
		{
			name:            "no symbols",
			options:         []Option{WithLength(16), WithoutSymbols()},
			expectedLength:  16,
			expectUppercase: true,
			expectLowercase: true,
			expectDigits:    true,
			expectSymbols:   false,
		},
		{
			name:            "with min requirements",
			options:         []Option{WithLength(20), WithMinUppercase(3), WithMinLowercase(3), WithMinDigits(3), WithMinSymbols(2)},
			expectedLength:  20,
			expectUppercase: true,
			expectLowercase: true,
			expectDigits:    true,
			expectSymbols:   true,
			minUppercase:    3,
			minLowercase:    3,
			minDigits:       3,
			minSymbols:      2,
		},
		{
			name:            "min requirements equal to length",
			options:         []Option{WithLength(10), WithMinRequirements(2, 3, 3, 2)},
			expectedLength:  10,
			expectUppercase: true,
			expectLowercase: true,
			expectDigits:    true,
			expectSymbols:   true,
			minUppercase:    2,
			minLowercase:    3,
			minDigits:       3,
			minSymbols:      2,
		},
		{
			name:        "invalid configuration",
			options:     []Option{WithLength(-1)},
			expectError: true,
			errorMsg:    "password length must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := Generate(tt.options...)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(password) != tt.expectedLength {
				t.Errorf("expected length %d, got %d", tt.expectedLength, len(password))
			}

			hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
			hasDigits := regexp.MustCompile(`[0-9]`).MatchString(password)
			hasSymbols := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`).MatchString(password)

			if tt.expectUppercase && !hasUppercase {
				t.Error("expected uppercase characters but found none")
			}
			if !tt.expectUppercase && hasUppercase {
				t.Error("found uppercase characters but none were expected")
			}
			if tt.expectLowercase && !hasLowercase {
				t.Error("expected lowercase characters but found none")
			}
			if !tt.expectLowercase && hasLowercase {
				t.Error("found lowercase characters but none were expected")
			}
			if tt.expectDigits && !hasDigits {
				t.Error("expected digits but found none")
			}
			if !tt.expectDigits && hasDigits {
				t.Error("found digits but none were expected")
			}
			if tt.expectSymbols && !hasSymbols {
				t.Error("expected symbols but found none")
			}
			if !tt.expectSymbols && hasSymbols {
				t.Error("found symbols but none were expected")
			}

			if tt.minUppercase > 0 {
				count := len(regexp.MustCompile(`[A-Z]`).FindAllString(password, -1))
				if count < tt.minUppercase {
					t.Errorf("expected at least %d uppercase characters, got %d", tt.minUppercase, count)
				}
			}
			if tt.minLowercase > 0 {
				count := len(regexp.MustCompile(`[a-z]`).FindAllString(password, -1))
				if count < tt.minLowercase {
					t.Errorf("expected at least %d lowercase characters, got %d", tt.minLowercase, count)
				}
			}
			if tt.minDigits > 0 {
				count := len(regexp.MustCompile(`[0-9]`).FindAllString(password, -1))
				if count < tt.minDigits {
					t.Errorf("expected at least %d digits, got %d", tt.minDigits, count)
				}
			}
			if tt.minSymbols > 0 {
				count := len(regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`).FindAllString(password, -1))
				if count < tt.minSymbols {
					t.Errorf("expected at least %d symbols, got %d", tt.minSymbols, count)
				}
			}

			validChars := []rune{}
			if tt.expectUppercase {
				validChars = append(validChars, upper...)
			}
			if tt.expectLowercase {
				validChars = append(validChars, lower...)
			}
			if tt.expectDigits {
				validChars = append(validChars, digits...)
			}
			if tt.expectSymbols {
				validChars = append(validChars, symbols...)
			}

			for _, char := range password {
				if !slices.Contains(validChars, char) {
					t.Errorf("password contains unexpected character: %c", char)
				}
			}
		})
	}
}

func TestGeneratorReuse(t *testing.T) {
	tests := []struct {
		name           string
		options        []Option
		generateCount  int
		expectedLength int
	}{
		{
			name:           "multiple generations with same config",
			options:        []Option{WithLength(12), WithoutSymbols()},
			generateCount:  10,
			expectedLength: 12,
		},
		{
			name:           "default config multiple times",
			options:        nil,
			generateCount:  5,
			expectedLength: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewGenerator(tt.options...)
			if err != nil {
				t.Fatalf("failed to create generator: %v", err)
			}

			passwords := make([]string, tt.generateCount)
			for i := 0; i < tt.generateCount; i++ {
				password, err := gen.Generate()
				if err != nil {
					t.Fatalf("failed to generate password %d: %v", i, err)
				}
				if len(password) != tt.expectedLength {
					t.Errorf("password %d has wrong length: expected %d, got %d", i, tt.expectedLength, len(password))
				}
				passwords[i] = password
			}

			for i := 0; i < len(passwords); i++ {
				for j := i + 1; j < len(passwords); j++ {
					if passwords[i] == passwords[j] {
						t.Errorf("generated identical passwords at positions %d and %d: %s", i, j, passwords[i])
					}
				}
			}
		})
	}
}

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
			options:         []Option{WithoutUppercase(), WithUpperCase()},
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

func TestConcurrentGeneration(t *testing.T) {
	tests := []struct {
		name        string
		options     []Option
		goroutines  int
		generations int
	}{
		{
			name:        "concurrent_default_generation",
			options:     nil,
			goroutines:  10,
			generations: 100,
		},
		{
			name:        "concurrent_with_options",
			options:     []Option{WithLength(20), WithMinRequirements(2, 2, 2, 2)},
			goroutines:  20,
			generations: 50,
		},
		{
			name:        "high_concurrency",
			options:     []Option{WithLength(12), WithoutSymbols()},
			goroutines:  100,
			generations: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			var mu sync.Mutex
			passwords := make(map[string]int)
			errors := make([]error, 0)

			for i := 0; i < tt.goroutines; i++ {
				wg.Add(1)
				go func(goroutineID int) {
					defer wg.Done()

					for j := 0; j < tt.generations; j++ {
						password, err := Generate(tt.options...)

						mu.Lock()
						if err != nil {
							errors = append(errors, err)
						} else {
							passwords[password]++
						}
						mu.Unlock()
					}
				}(i)
			}

			wg.Wait()

			if len(errors) > 0 {
				t.Fatalf("got %d errors during concurrent generation, first: %v", len(errors), errors[0])
			}

			expectedPasswords := tt.goroutines * tt.generations
			if len(passwords) != expectedPasswords {

				duplicates := expectedPasswords - len(passwords)
				duplicatePercentage := float64(duplicates) / float64(expectedPasswords) * 100

				if duplicatePercentage > 1.0 {
					t.Errorf("too many duplicate passwords: %d duplicates out of %d (%.2f%%)",
						duplicates, expectedPasswords, duplicatePercentage)
				}
			}

			for password := range passwords {
				if len(password) == 0 {
					t.Error("generated empty password")
				}
			}
		})
	}
}

func BenchmarkGenerate(b *testing.B) {
	benchmarks := []struct {
		name    string
		options []Option
	}{
		{
			name:    "default",
			options: nil,
		},
		{
			name:    "short_password",
			options: []Option{WithLength(8)},
		},
		{
			name:    "long_password",
			options: []Option{WithLength(64)},
		},
		{
			name:    "no_symbols",
			options: []Option{WithLength(16), WithoutSymbols()},
		},
		{
			name:    "digits_only",
			options: []Option{WithLength(16), WithoutUppercase(), WithoutLowercase(), WithoutSymbols()},
		},
		{
			name:    "with_min_requirements",
			options: []Option{WithLength(20), WithMinRequirements(3, 3, 3, 3)},
		},
		{
			name:    "complex_requirements",
			options: []Option{WithLength(32), WithMinUppercase(5), WithMinLowercase(5), WithMinDigits(5), WithMinSymbols(5)},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Generate(bm.options...)
				if err != nil {
					b.Fatalf("failed to generate password: %v", err)
				}
			}
		})
	}
}

// #1
// === RUN   BenchmarkGenerate/default
// BenchmarkGenerate/default
// BenchmarkGenerate/default-12              186519              6163 ns/op            2216 B/op        104 allocs/op
// === RUN   BenchmarkGenerate/short_password
// BenchmarkGenerate/short_password
// BenchmarkGenerate/short_password-12       349222              3333 ns/op            1408 B/op         55 allocs/op
// === RUN   BenchmarkGenerate/long_password
// BenchmarkGenerate/long_password
// BenchmarkGenerate/long_password-12         50288             23632 ns/op            7280 B/op        395 allocs/op
// === RUN   BenchmarkGenerate/no_symbols
// BenchmarkGenerate/no_symbols
// BenchmarkGenerate/no_symbols-12           204147              5624 ns/op            1992 B/op        103 allocs/op
// === RUN   BenchmarkGenerate/digits_only
// BenchmarkGenerate/digits_only
// BenchmarkGenerate/digits_only-12          199051              6055 ns/op            1656 B/op        101 allocs/op
// === RUN   BenchmarkGenerate/with_min_requirements
// BenchmarkGenerate/with_min_requirements
// BenchmarkGenerate/with_min_requirements-12                134268              8073 ns/op            2656 B/op        133 allocs/op
// === RUN   BenchmarkGenerate/complex_requirements
// BenchmarkGenerate/complex_requirements
// BenchmarkGenerate/complex_requirements-12                  94597             12355 ns/op            3848 B/op        206 allocs/op

// #2 После использования в shuffle буфферезированные значения
// BenchmarkGenerate/default
// BenchmarkGenerate/default-12              258638              4468 ns/op            1456 B/op         57 allocs/op
// === RUN   BenchmarkGenerate/short_password
// BenchmarkGenerate/short_password
// BenchmarkGenerate/short_password-12       459564              2454 ns/op            1040 B/op         32 allocs/op
// === RUN   BenchmarkGenerate/long_password
// BenchmarkGenerate/long_password
// BenchmarkGenerate/long_password-12         74122             16029 ns/op            4216 B/op        204 allocs/op
// === RUN   BenchmarkGenerate/no_symbols
// BenchmarkGenerate/no_symbols
// BenchmarkGenerate/no_symbols-12           293736              3907 ns/op            1232 B/op         56 allocs/op
// === RUN   BenchmarkGenerate/digits_only
// BenchmarkGenerate/digits_only
// BenchmarkGenerate/digits_only-12          269953              4248 ns/op             896 B/op         54 allocs/op
// === RUN   BenchmarkGenerate/with_min_requirements
// BenchmarkGenerate/with_min_requirements
// BenchmarkGenerate/with_min_requirements-12                212169              5559 ns/op            1704 B/op         74 allocs/op
// === RUN   BenchmarkGenerate/complex_requirements
// BenchmarkGenerate/complex_requirements
// BenchmarkGenerate/complex_requirements-12                 136867              8533 ns/op            2320 B/op        111

// #3 После изменения generatePassEntry
// BenchmarkGenerate/default-12              502628              2309 ns/op             688 B/op          9 allocs/op
// === RUN   BenchmarkGenerate/short_password
// BenchmarkGenerate/short_password
// BenchmarkGenerate/short_password-12       804098              1387 ns/op             656 B/op          8 allocs/op
// === RUN   BenchmarkGenerate/long_password
// BenchmarkGenerate/long_password
// BenchmarkGenerate/long_password-12        154750              7787 ns/op            1144 B/op         12 allocs/op
// === RUN   BenchmarkGenerate/no_symbols
// BenchmarkGenerate/no_symbols
// BenchmarkGenerate/no_symbols-12           498453              2174 ns/op             464 B/op          8 allocs/op
// === RUN   BenchmarkGenerate/digits_only
// BenchmarkGenerate/digits_only
// BenchmarkGenerate/digits_only-12          569515              1949 ns/op             128 B/op          6 allocs/op
// === RUN   BenchmarkGenerate/with_min_requirements
// BenchmarkGenerate/with_min_requirements
// BenchmarkGenerate/with_min_requirements-12                374904              2951 ns/op             744 B/op         14 allocs/op
// === RUN   BenchmarkGenerate/complex_requirements
// BenchmarkGenerate/complex_requirements
// BenchmarkGenerate/complex_requirements-12                 265250              4308 ns/op             784 B/op         15 allocs/op

// #4 После не большой оптимизации Generate. Был переход от строк на руны и ввиду избавления от кросс преобразований аллокации снизились
// === RUN   BenchmarkGenerate
// BenchmarkGenerate
// === RUN   BenchmarkGenerate/default
// BenchmarkGenerate/default
// BenchmarkGenerate/default-12              533120              2054 ns/op             464 B/op          6 allocs/op
// === RUN   BenchmarkGenerate/short_password
// BenchmarkGenerate/short_password
// BenchmarkGenerate/short_password-12      1052134              1122 ns/op             432 B/op          5 allocs/op
// === RUN   BenchmarkGenerate/long_password
// BenchmarkGenerate/long_password
// BenchmarkGenerate/long_password-12        161817              7436 ns/op             920 B/op          9 allocs/op
// === RUN   BenchmarkGenerate/no_symbols
// BenchmarkGenerate/no_symbols
// BenchmarkGenerate/no_symbols-12           524682              2031 ns/op             464 B/op          6 allocs/op
// === RUN   BenchmarkGenerate/digits_only
// BenchmarkGenerate/digits_only
// BenchmarkGenerate/digits_only-12          556071              2044 ns/op             464 B/op          6 allocs/op
// === RUN   BenchmarkGenerate/with_min_requirements
// BenchmarkGenerate/with_min_requirements
// BenchmarkGenerate/with_min_requirements-12                404955              2689 ns/op             520 B/op         11 allocs/op
// === RUN   BenchmarkGenerate/complex_requirements
// BenchmarkGenerate/complex_requirements
// BenchmarkGenerate/complex_requirements-12                 285315              4017 ns/op             560 B/op         12 allocs/op
