package passgen

type Option func(*Generator)

func WithLength(n int) Option {
	return func(c *Generator) {
		c.length = n
	}
}

func WithUpperCase() Option {
	return func(c *Generator) {
		c.uppercase = true
	}
}

func WithoutUppercase() Option {
	return func(c *Generator) {
		c.uppercase = false
	}
}

func WithLowercase() Option {
	return func(c *Generator) {
		c.lowercase = true
	}
}

func WithoutLowercase() Option {
	return func(c *Generator) {
		c.lowercase = false
	}
}

func WithDigits() Option {
	return func(c *Generator) {
		c.digits = true
	}
}

func WithoutDigits() Option {
	return func(c *Generator) {
		c.digits = false
	}
}

func WithSymbols() Option {
	return func(c *Generator) {
		c.symbols = true
	}
}

func WithoutSymbols() Option {
	return func(c *Generator) {
		c.symbols = false
	}
}

func WithMinUppercase(n int) Option {
	return func(c *Generator) {
		c.minUppercase = n
	}
}

func WithMinLowercase(n int) Option {
	return func(c *Generator) {
		c.minLowercase = n
	}
}

func WithMinDigits(n int) Option {
	return func(c *Generator) {
		c.minDigits = n
	}
}

func WithMinSymbols(n int) Option {
	return func(c *Generator) {
		c.minSymbols = n
	}
}

func WithMinRequirements(upper, lower, digits, symbols int) Option {
	return func(c *Generator) {
		c.minDigits = digits
		c.minUppercase = upper
		c.minLowercase = lower
		c.minSymbols = symbols
	}
}
