# dotenv

A lightweight, zero-dependency Go package for loading environment variables from `.env` and `.env.local` files, with built-in support for mapping variables directly into Go structs.

## Features

* **File Priority**: Automatically loads `.env` and overrides with `.env.local` if present.
* **Struct Mapping**: Unmarshal environment variables directly into strict-typed Go structs.
* **Env Generation**: Marshal structs back into `.env` formatted strings.
* **Shell Support**: Recognizes the `export` keyword.
* **Comment Handling**: Ignores lines starting with `#` and strips inline comments.
* **Smart Quoting**: Automatically handles values wrapped in single (`'`) or double (`"`) quotes.
* **Zero Dependencies**: Uses only the Go standard library.

## Installation

```bash
go get [github.com/rickferrdev/dotenv](https://github.com/rickferrdev/dotenv)
```

## Usage

### 1. Basic Loading (`Collect`)

Create your environment files:

**`.env`**
```bash
PORT=8080
DB_URL="postgres://user:password@localhost:5432/db"
DEBUG=true
```

Call `dotenv.Collect()` as early as possible in your `main` function (or `init`) to populate `os.Environ`.

```go
package main

import (
    "fmt"
    "os"
    "[github.com/rickferrdev/dotenv](https://github.com/rickferrdev/dotenv)"
)

func main() {
    // 1. Load variables from .env files into system environment
    dotenv.Collect()

    // 2. Access them normally
    port := os.Getenv("PORT")
    fmt.Printf("Server starting on port: %s\n", port)
}
```

### 2. Struct Mapping (`Unmarshal`)

Instead of parsing strings manually, map environment variables directly to a struct using the `env` tag.

```go
package main

import (
    "fmt"
    "log"
    "[github.com/rickferrdev/dotenv](https://github.com/rickferrdev/dotenv)"
)

type Config struct {
    Port      int     `env:"PORT"`
    DbURL     string  `env:"DB_URL"`
    Debug     bool    `env:"DEBUG"`
    RateLimit float64 `env:"RATE_LIMIT"`
}

func main() {
    dotenv.Collect() // Load files first

    var cfg Config

    // Parse environment variables into the struct
    if err := dotenv.Unmarshal(&cfg); err != nil {
        log.Fatalf("Could not parse config: %v", err)
    }

    fmt.Printf("Config loaded: Port %d, Debug %v\n", cfg.Port, cfg.Debug)
}
```

### 3. Generating .env Content (`Marshal`)

You can also convert a struct back into a `.env` formatted string.

```go
cfg := Config{
    Port:  9090,
    Debug: false,
}

data, err := dotenv.Marshal(&cfg)
if err != nil {
    panic(err)
}

fmt.Println(string(data))
// Output:
// PORT=9090
// DEBUG=false
// ...
```

## Configuration

You can override the files the package looks for by modifying the `FilenameVariables` slice before calling `Collect`.

```go
// Look for production files instead of default .env
dotenv.FilenameVariables = []string{".env.production", ".env"}
dotenv.Collect()
```

## How it Works

* **`Collect()`**: Iterates through `FilenameVariables`. It parses each line, strips `export` prefixes, handles quotes, cleans comments, and sets values using `os.Setenv`.
* **`Unmarshal()`**: Uses Go reflection to inspect struct tags (`env:"KEY"`) and automatically converts string environment values into the appropriate Go types (`int`, `bool`, `float`, `string`).
* **`Marshal()`**: Reads the struct values and tags to generate a key-value string suitable for `.env` files.
