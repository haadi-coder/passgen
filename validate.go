package passgen

import "fmt"

func (g *Generator) validate() error {
	if g.length <= 0 {
		return fmt.Errorf("password length must be greater than 0, got %d", g.length)
	}
	if g.length > 10000 {
		return fmt.Errorf("password length must not exceed 10000, got %d", g.length)
	}

	if g.minUppercase < 0 {
		return fmt.Errorf("minimum uppercase count cannot be negative, got %d", g.minUppercase)
	}
	if g.minLowercase < 0 {
		return fmt.Errorf("minimum lowercase count cannot be negative, got %d", g.minLowercase)
	}
	if g.minDigits < 0 {
		return fmt.Errorf("minimum digits count cannot be negative, got %d", g.minDigits)
	}
	if g.minSymbols < 0 {
		return fmt.Errorf("minimum symbols count cannot be negative, got %d", g.minSymbols)
	}

	if !g.uppercase && !g.lowercase && !g.digits && !g.symbols {
		return fmt.Errorf("at least one character set must be enabled")
	}

	if !g.uppercase && g.minUppercase > 0 {
		return fmt.Errorf("uppercase characters are disabled but minimum uppercase requirement is %d", g.minUppercase)
	}
	if !g.lowercase && g.minLowercase > 0 {
		return fmt.Errorf("lowercase characters are disabled but minimum lowercase requirement is %d", g.minLowercase)
	}
	if !g.digits && g.minDigits > 0 {
		return fmt.Errorf("digits are disabled but minimum digits requirement is %d", g.minDigits)
	}
	if !g.symbols && g.minSymbols > 0 {
		return fmt.Errorf("symbols are disabled but minimum symbols requirement is %d", g.minSymbols)
	}

	totalMin := g.minUppercase + g.minLowercase + g.minDigits + g.minSymbols
	if totalMin > g.length {
		return fmt.Errorf("sum of minimum requirements (%d) cannot exceed password length (%d)", totalMin, g.length)
	}

	return nil
}
