package passgen

type Option func(*config)

func WithLength(n int) Option {
	return func(c *config) {
		c.length = n
	}
}

func WithUppercase() Option {
	return func(c *config) {
		c.useUppercase = true
	}
}

func WithoutUppercase() Option {
	return func(c *config) {
		c.useUppercase = false
	}
}

func WithLowercase() Option {
	return func(c *config) {
		c.useLowercase = true
	}
}

func WithoutLowercase() Option {
	return func(c *config) {
		c.useLowercase = false
	}
}

func WithDigits() Option {
	return func(c *config) {
		c.useDigits = true
	}
}

func WithoutDigits() Option {
	return func(c *config) {
		c.useDigits = false
	}
}

func WithSymbols() Option {
	return func(c *config) {
		c.useSymbols = true
	}
}

func WithoutSymbols() Option {
	return func(c *config) {
		c.useSymbols = false
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
