# AGENTS.md - cloudlrc Agent Guide

## Build Commands

### Build
```bash
go build                    # Build for current platform
CGO_ENABLED=0 go build      # Build without CGO
```

### Cross-platform build
```bash
# Windows (amd64)
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

# With build info
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-X 'github.com/ludoux/cloudlrc/cmd.GIT_HASH=$(git show -s --format=%H)' -X 'github.com/ludoux/cloudlrc/cmd.BUILD_TIME=$(date -R)' -X 'github.com/ludoux/cloudlrc/cmd.GO_VER=$(go version)'"
```

### Test Commands
```bash
go test ./...               # Run all tests
go test -v ./...            # Run all tests with verbose output
go test -run TestName ./package  # Run single test in specific package
go test ./netease           # Run tests in specific package
```

### Code Quality
```bash
gofmt -s -w .               # Format code
go vet ./...                # Run go vet
golangci-lint run           # Run linter (if configured)
```

---

## Code Style Guidelines

### Naming Conventions

#### Types
- **Structs**: PascalCase with `_s` suffix (e.g., `NeteaseSingleMusic_s`, `CfgFile_s`)
- **Slice/Array types**: PascalCase with `_t` suffix (e.g., `NeteaseSingleMusics_t`, `LyricLines_t`)
- **Interfaces**: PascalCase without suffixes

#### Functions
- **Exported**: PascalCase (e.g., `NewLyric`, `DownloadSingleMusicLrc`)
- **Private constructors**: `new` prefix (e.g., `newLyricLine`, `newNeteaseSingleMusic`)
- **Private functions**: camelCase (e.g., `convertEncoding`, `responseController`)

#### Variables
- **Local variables**: camelCase (e.g., `rt`, `err`, `resp`)
- **Package-level exported**: PascalCase (e.g., `Client`, `VERSION`)
- **Package-level private**: camelCase (e.g., `responseChannel`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `GIT_HASH`, `BUILD_TIME`)

### Import Organization
```go
import (
    "fmt"                    // Standard library first
    "os"
    "strings"

    "github.com/ludoux/cloudlrc/cfg"  // Internal packages
    "github.com/ludoux/cloudlrc/lyric"

    "github.com/spf13/cobra"           // Third-party packages
    "github.com/spf13/cast"
)
```
- Standard library first, then internal packages, then third-party
- Group by type with blank lines between groups
- Use aliases for clarity: `cookiejar "github.com/orirawlings/persistent-cookiejar"`

### Formatting
- Use `gofmt -s` for formatting
- Tabs for indentation (Go standard)
- Use string builder for string concatenation in loops
- Prefer `%02d` format for aligned numbers

### Type Conversions
- Always use `cast` package from spf13/cast for type conversions
- Example: `cast.ToInt64(value)`, `cast.ToString(value)`, `cast.ToBool(value)`

### Error Handling
- **Fatal errors**: Use `log.Fatal()` or `log.Panic()` (e.g., client initialization failures)
- **Command errors**: Use `fmt.Println()` + `os.Exit(1)` for CLI commands
- **Function errors**: Return `error` type, include context in messages
- Example: `return fmt.Errorf("获取服务端歌词错误: %s", err.Error())`

### Comments
- **Package comments**: Not required in this codebase
- **Function comments**: Optional, use C-style `/* */` for multi-line
- **Inline comments**: Minimal, use only when logic is unclear
- Language: Mix of Chinese and English accepted (existing codebase uses both)

### Struct Organization
```go
type Type_s struct {
    field1 int                  // Basic fields first
    field2 string

    nested *OtherType_s         // Pointer fields
    items  []ItemType_t         // Slice fields

    config CfgFile_s            // Config fields last
}
```

### Method Receivers
- Use pointer receivers `*Type` for methods that modify the struct
- Use value receivers `Type` for methods that don't modify (rare in this codebase)

### File Organization
- One type per file preferred for complex types
- Keep related functions in same file as their type
- Command handlers in `cmd/` package
- Models in respective feature packages

### Constants
- Define package-level constants at top of file
- Use for magic numbers and common strings
- Example: `responseChannel = make(chan string, 15)`

### Dependencies
- Use `github.com/imroc/req/v3` for HTTP client
- Use `github.com/spf13/cobra` for CLI framework
- Use `github.com/buger/jsonparser` for JSON parsing (efficient)
- Use `golang.org/x/text` for encoding/decoding

---

## Project Structure

```
cloudlrc/
├── main.go              # Entry point
├── cmd/                 # CLI commands using cobra
├── cfg/                 # Configuration management
├── lyric/               # Lyric parsing and formatting
├── netease/             # Netease Music API client
│   ├── client.go        # HTTP client setup
│   ├── model.go         # Data models
│   ├── netease.go       # Business logic
│   └── crypto.go        # Encryption utilities
└── utils/               # Utility functions
```

---

## Key Patterns

### Async Operations
- Use `github.com/panjf2000/ants/v2` for worker pools
- Use `sync.WaitGroup` for coordination
- Use channels for communication

### Config Loading
- Load from `config.ini` using `gopkg.in/ini.v1`
- Sections: `[file]`, `[lrc]`
- Provide functions with `WithCfg` suffix for testing

### CLI Commands
- Use cobra framework
- Define commands in `cmd/` package
- Validate arguments in command Run functions
- Print errors to stdout, exit with status 1

---

## Testing Notes
- No test files currently in codebase
- Write tests using standard `testing` package
- Test files should be named `*_test.go`
- Use table-driven tests for multiple cases
