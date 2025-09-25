package passgen

import "fmt"

type config struct {
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

func defaultConfig() *config {
	return &config{
		length:    16,
		uppercase: true,
		lowercase: true,
		digits:    true,
		symbols:   true,
	}
}

func (c *config) validate() error {
	if c.length <= 0 {
		return fmt.Errorf("password length must be greater than 0, got %d", c.length)
	}
	if c.length > 10000 {
		return fmt.Errorf("password length must not exceed 10000, got %d", c.length)
	}

	if c.minUppercase < 0 {
		return fmt.Errorf("minimum uppercase count cannot be negative, got %d", c.minUppercase)
	}
	if c.minLowercase < 0 {
		return fmt.Errorf("minimum lowercase count cannot be negative, got %d", c.minLowercase)
	}
	if c.minDigits < 0 {
		return fmt.Errorf("minimum digits count cannot be negative, got %d", c.minDigits)
	}
	if c.minSymbols < 0 {
		return fmt.Errorf("minimum symbols count cannot be negative, got %d", c.minSymbols)
	}

	if !c.uppercase && !c.lowercase && !c.digits && !c.symbols {
		return fmt.Errorf("at least one character set must be enabled")
	}

	if !c.uppercase && c.minUppercase > 0 {
		return fmt.Errorf("uppercase characters are disabled but minimum uppercase requirement is %d", c.minUppercase)
	}
	if !c.lowercase && c.minLowercase > 0 {
		return fmt.Errorf("lowercase characters are disabled but minimum lowercase requirement is %d", c.minLowercase)
	}
	if !c.digits && c.minDigits > 0 {
		return fmt.Errorf("digits are disabled but minimum digits requirement is %d", c.minDigits)
	}
	if !c.symbols && c.minSymbols > 0 {
		return fmt.Errorf("symbols are disabled but minimum symbols requirement is %d", c.minSymbols)
	}

	totalMin := c.minUppercase + c.minLowercase + c.minDigits + c.minSymbols
	if totalMin > c.length {
		return fmt.Errorf("sum of minimum requirements (%d) cannot exceed password length (%d)", totalMin, c.length)
	}

	return nil
}
