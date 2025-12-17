# dotenv

A lightweight, zero-dependency Go package for loading environment variables from `.env` and `.env.local` files into your application.

## Features

* **File Priority**: Automatically loads `.env` and overrides with `.env.local` if present.
* **Shell Support**: Recognizes the `export` keyword.
* **Comment Handling**: Ignores lines starting with `#` and strips inline comments.
* **Smart Quoting**: Automatically handles values wrapped in single (`'`) or double (`"`) quotes.
* **Zero Dependencies**: Uses only the Go standard library.

## Installation

```bash
go get github.com/rickferrdev/dotenv

```

## Usage

### 1. Create your environment files

**`.env`**

```bash
PORT=8080
DB_URL="postgres://user:password@localhost:5432/db"
# This is a comment
DEBUG=true

```

**`.env.local`**

```bash
PORT=3000 # Overrides the value in .env

```

### 2. Load variables in Go

Call `dotenv.Collect()` as early as possible in your `main` function (or in an `init` function) to populate `os.Environ`.

```go
package main

import (
    "fmt"
    "os"
    "github.com/rickferrdev/dotenv"
)

func main() {
    // Load variables from .env and .env.local
    dotenv.Collect()

    port := os.Getenv("PORT")
    fmt.Printf("Server starting on port: %s\n", port)
}

```

## How it Works

* **`Collect()`**: Iterates through `FilenameVariables` (defaulting to `.env` and `.env.local`). It parses each line, strips the `export` prefix if it exists, and ignores comments.
* **`quotes()`**: A helper function that ensures values like `KEY="value"` or `KEY='value'` are stored simply as `value`, while also cleaning up trailing inline comments.

## Configuration

You can override the files the package looks for by modifying the `FilenameVariables` slice before calling `Collect`:

```go
dotenv.FilenameVariables = []string{".env.production", ".env"}
dotenv.Collect()

```