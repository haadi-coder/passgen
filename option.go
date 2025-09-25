package passgen

type Option func(*config)

func WithLength(n int) Option {
	return func(c *config) {
		c.length = n
	}
}

func WithUpperCase() Option {
	return func(c *config) {
		c.uppercase = true
	}
}

func WithoutUppercase() Option {
	return func(c *config) {
		c.uppercase = false
	}
}

func WithLowercase() Option {
	return func(c *config) {
		c.lowercase = true
	}
}

func WithoutLowercase() Option {
	return func(c *config) {
		c.lowercase = false
	}
}

func WithDigits() Option {
	return func(c *config) {
		c.digits = true
	}
}

func WithoutDigits() Option {
	return func(c *config) {
		c.digits = false
	}
}

func WithSymbols() Option {
	return func(c *config) {
		c.symbols = true
	}
}

func WithoutSymbols() Option {
	return func(c *config) {
		c.symbols = false
	}
}

func WithMinUppercase(n int) Option {
	return func(c *config) {
		c.minUppercase = n
	}
}

func WithMinLowercase(n int) Option {
	return func(c *config) {
		c.minLowercase = n
	}
}

func WithMinDigits(n int) Option {
	return func(c *config) {
		c.minDigits = n
	}
}

func WithMinSymbols(n int) Option {
	return func(c *config) {
		c.minSymbols = n
	}
}

func WithMinRequirements(upper, lower, digits, symbols int) Option {
	return func(c *config) {
		c.minDigits = digits
		c.minUppercase = upper
		c.minLowercase = lower
		c.minSymbols = symbols
	}
}
