# passgen

A fast, secure password generator for Go with extensive customization options.

## Features

- **Secure**: Uses `crypto/rand` for cryptographically secure random generation
- **Fast**: Optimized for performance with minimal allocations
- **Flexible**: Comprehensive options for password requirements
- **Thread-safe**: Safe for concurrent use

## Installation

```bash
go get github.com/haadi-coder/passgen
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/haadi-coder/passgen"
)

func main() {
    // Generate with defaults (16 chars, mixed case + digits + symbols)
    password, err := passgen.Generate()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(password)
    
    // Custom password
    password, err = passgen.Generate(
        passgen.WithLength(20),
        passgen.WithoutSymbols(),
        passgen.WithMinRequirements(2, 2, 2, 0), // min: upper, lower, digits, symbols
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(password)
}
```

## Options

| Option | Description |
|--------|-------------|
| `WithLength(n)` | Set password length (1-10000) |
| `WithUppercase()` / `WithoutUppercase()` | Include/exclude uppercase letters |
| `WithLowercase()` / `WithoutLowercase()` | Include/exclude lowercase letters |
| `WithDigits()` / `WithoutDigits()` | Include/exclude digits |
| `WithSymbols()` / `WithoutSymbols()` | Include/exclude symbols |
| `WithMinUppercase(n)` | Minimum uppercase characters |
| `WithMinLowercase(n)` | Minimum lowercase characters |
| `WithMinDigits(n)` | Minimum digits |
| `WithMinSymbols(n)` | Minimum symbols |
| `WithMinRequirements(u,l,d,s)` | Set all minimums at once |

## Reusable Generator

```go
gen, err := passgen.NewGenerator(
    passgen.WithLength(12),
    passgen.WithMinRequirements(2, 2, 2, 0),
)
if err != nil {
    log.Fatal(err)
}

for i := 0; i < 10; i++ {
    password, err := gen.Generate()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(password)
}
```

## Performance

High-performance implementation achieving **~490K passwords/sec** for default configuration with minimal memory allocations.
