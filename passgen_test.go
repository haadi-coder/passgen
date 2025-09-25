package passgen

import (
	"math"
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

func TestStatisticalUniformity(t *testing.T) {
	tests := []struct {
		name           string
		options        []Option
		expectedLength int
		charSets       [][]rune
	}{
		{
			name:           "default configuration",
			options:        nil,
			expectedLength: 16,
			charSets:       [][]rune{upper, lower, digits, symbols},
		},
		{
			name:           "only uppercase",
			options:        []Option{WithLength(12), WithoutLowercase(), WithoutDigits(), WithoutSymbols()},
			expectedLength: 12,
			charSets:       [][]rune{upper},
		},
		{
			name:           "only lowercase",
			options:        []Option{WithLength(12), WithoutUppercase(), WithoutDigits(), WithoutSymbols()},
			expectedLength: 12,
			charSets:       [][]rune{lower},
		},
		{
			name:           "only digits",
			options:        []Option{WithLength(6), WithoutUppercase(), WithoutLowercase(), WithoutSymbols()},
			expectedLength: 6,
			charSets:       [][]rune{digits},
		},
		{
			name:           "only symbols",
			options:        []Option{WithLength(8), WithoutUppercase(), WithoutLowercase(), WithoutDigits()},
			expectedLength: 8,
			charSets:       [][]rune{symbols},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewGenerator(tt.options...)
			if err != nil {
				t.Fatalf("failed to create generator: %v", err)
			}

			var allChars []rune
			for _, charSet := range tt.charSets {
				allChars = append(allChars, charSet...)
			}
			if len(allChars) == 0 {
				t.Fatal("no characters available for testing")
			}

			const numPasswords = 10000
			totalChars := numPasswords * tt.expectedLength
			// Ожидаемая частота для каждого символа
			expectedFreq := float64(totalChars) / float64(len(allChars))
			// Допуск: 3 стандартных отклонения для биномиального распределения
			// Стандартное отклонение = sqrt(p * (1-p) * N), где p = 1/len(allChars), N = totalChars
			p := 1.0 / float64(len(allChars))
			stdDev := math.Sqrt(p * (1 - p) * float64(totalChars))
			// Допуск ±4σ (99.9937% доверительный интервал)
			tolerance := 4 * stdDev

			charCount := make(map[rune]int)
			var mu sync.Mutex
			var wg sync.WaitGroup

			// Генерируем пароли конкурентно
			for i := 0; i < numPasswords; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					password, err := gen.Generate()
					if err != nil {
						t.Errorf("failed to generate password: %v", err)
						return
					}
					if len(password) != tt.expectedLength {
						t.Errorf("expected length %d, got %d", tt.expectedLength, len(password))
						return
					}

					mu.Lock()
					for _, char := range password {
						if slices.Contains(allChars, char) {
							charCount[char]++
						} else {
							t.Errorf("unexpected character %c in password", char)
						}
					}
					mu.Unlock()
				}()
			}
			wg.Wait()

			// Проверяем частоту каждого символа
			for _, char := range allChars {
				count := float64(charCount[char])
				if math.Abs(count-expectedFreq) > tolerance {
					t.Errorf("character %c has frequency %v, expected %v ± %v", char, count, expectedFreq, tolerance)
				}
			}

			// Проверяем, что все ожидаемые символы присутствуют
			for _, char := range allChars {
				if _, exists := charCount[char]; !exists {
					t.Errorf("character %c was not generated at all", char)
				}
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid_default_config",
			config:      defaultConfig(),
			expectError: false,
		},
		{
			name: "valid_custom_config",
			config: &config{
				length:    20,
				uppercase: true,
				lowercase: true,
				digits:    true,
				symbols:   false,
			},
			expectError: false,
		},
		{
			name: "valid_with_min_requirements",
			config: &config{
				length:       16,
				uppercase:    true,
				lowercase:    true,
				digits:       true,
				symbols:      true,
				minUppercase: 2,
				minLowercase: 2,
				minDigits:    2,
				minSymbols:   2,
			},
			expectError: false,
		},
		{
			name: "valid_single_character_set",
			config: &config{
				length:    10,
				uppercase: false,
				lowercase: false,
				digits:    true,
				symbols:   false,
			},
			expectError: false,
		},
		{
			name: "valid_max_length",
			config: &config{
				length:    10000,
				uppercase: true,
				lowercase: false,
				digits:    false,
				symbols:   false,
			},
			expectError: false,
		},
		{
			name: "valid_min_length",
			config: &config{
				length:    1,
				uppercase: true,
				lowercase: false,
				digits:    false,
				symbols:   false,
			},
			expectError: false,
		},
		{
			name: "valid_zero_min_requirements",
			config: &config{
				length:       10,
				uppercase:    true,
				lowercase:    true,
				digits:       true,
				symbols:      true,
				minUppercase: 0,
				minLowercase: 0,
				minDigits:    0,
				minSymbols:   0,
			},
			expectError: false,
		},
		{
			name: "valid_min_requirements_equal_length",
			config: &config{
				length:       4,
				uppercase:    true,
				lowercase:    true,
				digits:       true,
				symbols:      true,
				minUppercase: 1,
				minLowercase: 1,
				minDigits:    1,
				minSymbols:   1,
			},
			expectError: false,
		},

		{
			name: "zero_length",
			config: &config{
				length:    0,
				uppercase: true,
			},
			expectError: true,
			errorMsg:    "password length must be greater than 0, got 0",
		},
		{
			name: "negative_length",
			config: &config{
				length:    -5,
				uppercase: true,
			},
			expectError: true,
			errorMsg:    "password length must be greater than 0, got -5",
		},
		{
			name: "too_long_password",
			config: &config{
				length:    10001,
				uppercase: true,
			},
			expectError: true,
			errorMsg:    "password length must not exceed 10000, got 10001",
		},
		{
			name: "way_too_long_password",
			config: &config{
				length:    50000,
				uppercase: true,
			},
			expectError: true,
			errorMsg:    "password length must not exceed 10000, got 50000",
		},

		{
			name: "negative_min_uppercase",
			config: &config{
				length:       10,
				uppercase:    true,
				minUppercase: -1,
			},
			expectError: true,
			errorMsg:    "minimum uppercase count cannot be negative, got -1",
		},
		{
			name: "negative_min_lowercase",
			config: &config{
				length:       10,
				lowercase:    true,
				minLowercase: -3,
			},
			expectError: true,
			errorMsg:    "minimum lowercase count cannot be negative, got -3",
		},
		{
			name: "negative_min_digits",
			config: &config{
				length:    10,
				digits:    true,
				minDigits: -2,
			},
			expectError: true,
			errorMsg:    "minimum digits count cannot be negative, got -2",
		},
		{
			name: "negative_min_symbols",
			config: &config{
				length:     10,
				symbols:    true,
				minSymbols: -4,
			},
			expectError: true,
			errorMsg:    "minimum symbols count cannot be negative, got -4",
		},
		{
			name: "multiple_negative_minimums",
			config: &config{
				length:       10,
				uppercase:    true,
				lowercase:    true,
				minUppercase: -1,
				minLowercase: -2,
			},
			expectError: true,
			errorMsg:    "minimum uppercase count cannot be negative, got -1",
		},

		{
			name: "all_character_sets_disabled",
			config: &config{
				length:    10,
				uppercase: false,
				lowercase: false,
				digits:    false,
				symbols:   false,
			},
			expectError: true,
			errorMsg:    "at least one character set must be enabled",
		},

		{
			name: "min_uppercase_with_disabled_uppercase",
			config: &config{
				length:       10,
				uppercase:    false,
				lowercase:    true,
				minUppercase: 2,
			},
			expectError: true,
			errorMsg:    "uppercase characters are disabled but minimum uppercase requirement is 2",
		},
		{
			name: "min_lowercase_with_disabled_lowercase",
			config: &config{
				length:       10,
				uppercase:    true,
				lowercase:    false,
				minLowercase: 3,
			},
			expectError: true,
			errorMsg:    "lowercase characters are disabled but minimum lowercase requirement is 3",
		},
		{
			name: "min_digits_with_disabled_digits",
			config: &config{
				length:    10,
				uppercase: true,
				digits:    false,
				minDigits: 1,
			},
			expectError: true,
			errorMsg:    "digits are disabled but minimum digits requirement is 1",
		},
		{
			name: "min_symbols_with_disabled_symbols",
			config: &config{
				length:     10,
				uppercase:  true,
				symbols:    false,
				minSymbols: 5,
			},
			expectError: true,
			errorMsg:    "symbols are disabled but minimum symbols requirement is 5",
		},

		{
			name: "min_requirements_exceed_length_simple",
			config: &config{
				length:       5,
				uppercase:    true,
				lowercase:    true,
				digits:       true,
				symbols:      true,
				minUppercase: 2,
				minLowercase: 2,
				minDigits:    2,
				minSymbols:   2,
			},
			expectError: true,
			errorMsg:    "sum of minimum requirements (8) cannot exceed password length (5)",
		},
		{
			name: "min_requirements_exceed_length_by_one",
			config: &config{
				length:       10,
				uppercase:    true,
				lowercase:    true,
				digits:       true,
				symbols:      true,
				minUppercase: 3,
				minLowercase: 3,
				minDigits:    3,
				minSymbols:   3,
			},
			expectError: true,
			errorMsg:    "sum of minimum requirements (12) cannot exceed password length (10)",
		},
		{
			name: "single_requirement_exceeds_length",
			config: &config{
				length:       5,
				uppercase:    true,
				lowercase:    false,
				digits:       false,
				symbols:      false,
				minUppercase: 6,
			},
			expectError: true,
			errorMsg:    "sum of minimum requirements (6) cannot exceed password length (5)",
		},
		{
			name: "min_requirements_way_exceed_length",
			config: &config{
				length:       2,
				uppercase:    true,
				lowercase:    true,
				digits:       true,
				symbols:      true,
				minUppercase: 10,
				minLowercase: 10,
				minDigits:    10,
				minSymbols:   10,
			},
			expectError: true,
			errorMsg:    "sum of minimum requirements (40) cannot exceed password length (2)",
		},

		{
			name: "large_valid_minimums",
			config: &config{
				length:       1000,
				uppercase:    true,
				lowercase:    true,
				digits:       true,
				symbols:      true,
				minUppercase: 250,
				minLowercase: 250,
				minDigits:    250,
				minSymbols:   250,
			},
			expectError: false,
		},
		{
			name: "mixed_enabled_disabled_sets",
			config: &config{
				length:       20,
				uppercase:    true,
				lowercase:    false,
				digits:       true,
				symbols:      false,
				minUppercase: 10,
				minDigits:    10,
			},
			expectError: false,
		},
		{
			name: "mixed_errors_priority_check",
			config: &config{
				length:       -1,
				uppercase:    false,
				lowercase:    false,
				digits:       false,
				symbols:      false,
				minUppercase: -1,
			},
			expectError: true,
			errorMsg:    "password length must be greater than 0, got -1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()

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
				}
			}
		})
	}
}
