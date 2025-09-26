package passgen

import (
	"strings"
	"testing"
)

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
				length:       20,
				useUppercase: true,
				useLowercase: true,
				useDigits:    true,
				useSymbols:   false,
			},
			expectError: false,
		},
		{
			name: "valid_with_min_requirements",
			config: &config{
				length:       16,
				useUppercase: true,
				useLowercase: true,
				useDigits:    true,
				useSymbols:   true,
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
				length:       10,
				useUppercase: false,
				useLowercase: false,
				useDigits:    true,
				useSymbols:   false,
			},
			expectError: false,
		},
		{
			name: "valid_max_length",
			config: &config{
				length:       10000,
				useUppercase: true,
				useLowercase: false,
				useDigits:    false,
				useSymbols:   false,
			},
			expectError: false,
		},
		{
			name: "valid_min_length",
			config: &config{
				length:       1,
				useUppercase: true,
				useLowercase: false,
				useDigits:    false,
				useSymbols:   false,
			},
			expectError: false,
		},
		{
			name: "valid_zero_min_requirements",
			config: &config{
				length:       10,
				useUppercase: true,
				useLowercase: true,
				useDigits:    true,
				useSymbols:   true,
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
				useUppercase: true,
				useLowercase: true,
				useDigits:    true,
				useSymbols:   true,
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
				length:       0,
				useUppercase: true,
			},
			expectError: true,
			errorMsg:    "password length must be greater than 0, got 0",
		},
		{
			name: "negative_length",
			config: &config{
				length:       -5,
				useUppercase: true,
			},
			expectError: true,
			errorMsg:    "password length must be greater than 0, got -5",
		},
		{
			name: "too_long_password",
			config: &config{
				length:       10001,
				useUppercase: true,
			},
			expectError: true,
			errorMsg:    "password length must not exceed 10000, got 10001",
		},
		{
			name: "way_too_long_password",
			config: &config{
				length:       50000,
				useUppercase: true,
			},
			expectError: true,
			errorMsg:    "password length must not exceed 10000, got 50000",
		},

		{
			name: "negative_min_uppercase",
			config: &config{
				length:       10,
				useUppercase: true,
				minUppercase: -1,
			},
			expectError: true,
			errorMsg:    "minimum uppercase count cannot be negative, got -1",
		},
		{
			name: "negative_min_lowercase",
			config: &config{
				length:       10,
				useLowercase: true,
				minLowercase: -3,
			},
			expectError: true,
			errorMsg:    "minimum lowercase count cannot be negative, got -3",
		},
		{
			name: "negative_min_digits",
			config: &config{
				length:    10,
				useDigits: true,
				minDigits: -2,
			},
			expectError: true,
			errorMsg:    "minimum digits count cannot be negative, got -2",
		},
		{
			name: "negative_min_symbols",
			config: &config{
				length:     10,
				useSymbols: true,
				minSymbols: -4,
			},
			expectError: true,
			errorMsg:    "minimum symbols count cannot be negative, got -4",
		},
		{
			name: "multiple_negative_minimums",
			config: &config{
				length:       10,
				useUppercase: true,
				useLowercase: true,
				minUppercase: -1,
				minLowercase: -2,
			},
			expectError: true,
			errorMsg:    "minimum uppercase count cannot be negative, got -1",
		},

		{
			name: "all_character_sets_disabled",
			config: &config{
				length:       10,
				useUppercase: false,
				useLowercase: false,
				useDigits:    false,
				useSymbols:   false,
			},
			expectError: true,
			errorMsg:    "at least one character set must be enabled",
		},

		{
			name: "min_uppercase_with_disabled_uppercase",
			config: &config{
				length:       10,
				useUppercase: false,
				useLowercase: true,
				minUppercase: 2,
			},
			expectError: true,
			errorMsg:    "uppercase characters are disabled but minimum uppercase requirement is 2",
		},
		{
			name: "min_lowercase_with_disabled_lowercase",
			config: &config{
				length:       10,
				useUppercase: true,
				useLowercase: false,
				minLowercase: 3,
			},
			expectError: true,
			errorMsg:    "lowercase characters are disabled but minimum lowercase requirement is 3",
		},
		{
			name: "min_digits_with_disabled_digits",
			config: &config{
				length:       10,
				useUppercase: true,
				useDigits:    false,
				minDigits:    1,
			},
			expectError: true,
			errorMsg:    "digits are disabled but minimum digits requirement is 1",
		},
		{
			name: "min_symbols_with_disabled_symbols",
			config: &config{
				length:       10,
				useUppercase: true,
				useSymbols:   false,
				minSymbols:   5,
			},
			expectError: true,
			errorMsg:    "symbols are disabled but minimum symbols requirement is 5",
		},

		{
			name: "min_requirements_exceed_length_simple",
			config: &config{
				length:       5,
				useUppercase: true,
				useLowercase: true,
				useDigits:    true,
				useSymbols:   true,
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
				useUppercase: true,
				useLowercase: true,
				useDigits:    true,
				useSymbols:   true,
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
				useUppercase: true,
				useLowercase: false,
				useDigits:    false,
				useSymbols:   false,
				minUppercase: 6,
			},
			expectError: true,
			errorMsg:    "sum of minimum requirements (6) cannot exceed password length (5)",
		},
		{
			name: "min_requirements_way_exceed_length",
			config: &config{
				length:       2,
				useUppercase: true,
				useLowercase: true,
				useDigits:    true,
				useSymbols:   true,
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
				useUppercase: true,
				useLowercase: true,
				useDigits:    true,
				useSymbols:   true,
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
				useUppercase: true,
				useLowercase: false,
				useDigits:    true,
				useSymbols:   false,
				minUppercase: 10,
				minDigits:    10,
			},
			expectError: false,
		},
		{
			name: "mixed_errors_priority_check",
			config: &config{
				length:       -1,
				useUppercase: false,
				useLowercase: false,
				useDigits:    false,
				useSymbols:   false,
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
