# Nature Dopes CLI - Progress Tracker

**Last Updated**: 2025-11-10
**Current Phase**: Phase 2 - Configuration Management (IN PROGRESS)
**Next Phase**: Phase 2 - Complete Save() function and cmd/config.go

---

## ✅ Phase 1: Foundation - COMPLETED

### What You Built:
1. ✅ Project structure created
2. ✅ Go module initialized (`go.mod`)
3. ✅ Cobra CLI framework installed
4. ✅ `main.go` - Entry point created
5. ✅ `cmd/root.go` - Root command with flags created
6. ✅ Tested and verified flags work

### Files Created:
```
naturedopes-cli/
├── main.go                 ✅
├── go.mod                  ✅
├── go.sum                  ✅
├── BUILD_GUIDE.md          ✅
├── PROGRESS.md             ✅ (this file)
└── cmd/
    └── root.go             ✅
```

### What You Learned:
- ✅ Go project structure and packages
- ✅ `package` declaration and imports
- ✅ Cobra command structure (`cobra.Command`)
- ✅ The `init()` function and when it runs
- ✅ Global variables in Go
- ✅ CLI flags (persistent flags)
- ✅ Pointers (`&variable`)
- ✅ Error handling pattern (`if err != nil`)
- ✅ Public vs private functions (capital vs lowercase)

### Code You Wrote:

**main.go**:
```go
package main

import (
	"github.com/wattsmainsanglais/naturedopes-cli/cmd"
)

func main() {
	cmd.Execute()
}
```

**cmd/root.go**:
```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiURL string
	apiKey string
)

var rootCmd = &cobra.Command{
	Use:   "naturedopes-cli",
	Short: "CLI tool for Nature Dopes API",
	Long: `A command-line interface for interacting with the Nature Dopes API.

Manage images, search for flora species, and work with API keys.

Example usage:
  naturedopes-cli images list
  naturedopes-cli keys create --name "My Key"`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "http://localhost:8080", "API base URL")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for authentication")
}
```

### Testing Results:
```bash
# Help text displays correctly
go run main.go --help
✅ Shows CLI description

# Flags are recognized
go run main.go --api-url https://test.com
✅ No error

# Invalid flags are caught
go run main.go --invalid-flag
✅ Shows error
```

---

## 🚧 Phase 2: Configuration Management - IN PROGRESS

### Progress So Far:
1. ✅ Created `pkg/config/config.go` file
2. ✅ Defined `Config` struct with JSON tags
3. ✅ Implemented `getConfigFilePath()` helper function
4. ✅ Implemented `Load()` function (needs 1 typo fix)
5. ⏸️ Need to implement `Save()` function
6. ⏸️ Need to implement `Set()` function
7. ⏸️ Need to implement `Get()` function
8. ⏸️ Need to implement `Clear()` function
9. ⏸️ Need to create `cmd/config.go`

### What You've Learned So Far:

#### 1. **Hidden Files and Directories**
- Files/directories starting with `.` are hidden in Unix/Linux/Mac
- View with `ls -a` (not just `ls`)
- Common examples: `.gitconfig`, `.ssh/`, `.bashrc`
- Used for config files to keep home directory clean

#### 2. **Error Wrapping in Go**
```go
// GOOD ✅
return "", fmt.Errorf("failed to get home directory: %w", err)
//                    ^lowercase  ^colon+space  ^%w wraps error

// BAD ❌
return "", fmt.Errorf("Error %w", err)  // Not descriptive, capitalized
```
- Always use `%w` to wrap errors (preserves error chain)
- Use lowercase messages (errors appear mid-sentence when chained)
- Add descriptive context with `: %w` pattern

#### 3. **File Paths in Go**
```go
filepath.Join(homeDir, ".naturedopes-cli", "config.json")
```
- Cross-platform path handling (works on Windows/Linux/Mac)
- Automatically uses correct separator (`/` or `\`)

#### 4. **Config Storage Location**
```go
os.UserHomeDir()  // Gets /home/andrew (or equivalent)
```
- Config stored in user's home directory (not current working directory)
- CLI can be run from anywhere, config is always accessible
- Standard convention for CLI tools

#### 5. **Structs and JSON Tags**
```go
type Config struct {
    ApiURL string `json:"api_url"`  // Struct field ↔ JSON field
    ApiKey string `json:"api_key"`
}
```
- JSON tags map Go field names to JSON keys
- Go field names are PascalCase (exported/public)
- JSON keys are snake_case (convention)

#### 6. **Pointers and Return Types**
```go
func Load() (*Config, error)  // Returns pointer to Config
return nil, err              // Return nil for pointer on error
return &config, nil          // Return pointer to struct on success
```
- `*Config` means "pointer to Config"
- `&config` means "address of config"
- Return `nil` when you can't return a valid pointer

#### 7. **Checking File Existence**
```go
if _, err := os.Stat(path); os.IsNotExist(err) {
    // File doesn't exist - return default config (not an error!)
}
```
- Missing config file on first run is NORMAL
- Return default values, don't treat as error

#### 8. **JSON Unmarshaling**
```go
var config Config
err := json.Unmarshal(fileContent, &config)
// fileContent ([]byte) → config (struct)
```
- Converts JSON bytes to Go struct
- Pass pointer to struct (`&config`) so it can be modified

### Code You've Written:

**pkg/config/config.go** (partial):
```go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ApiURL string `json:"api_url"`
	ApiKey string `json:"api_key"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("user home directory not found: %w", err)
	}

	fullPath := filepath.Join(homeDir, ".naturedopes-cli", "config.json")
	return fullPath, nil
}

func Load() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("couldn't get home directory: %w", err)
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{
			ApiURL: "http://localhost8080",  // ⚠️ FIX: Missing colon after localhost
			ApiKey: "",
		}, nil
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file: %w", err)
	}

	var config Config
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal JSON: %w", err)
	}

	return &config, nil
}
```

### Next Steps When You Return:

#### 1. Fix the typo in Load() (line 37)
Change `"http://localhost8080"` to `"http://localhost:8080"`

#### 2. Implement the `Save()` function
This function will:
- Get the config file path
- Create the `.naturedopes-cli` directory if it doesn't exist
- Convert Config struct to pretty JSON
- Write the JSON to the file

**New concepts you'll need:**
```go
// Create directory (and parents) if needed
os.MkdirAll(dirPath, 0755)

// Get directory from full path
filepath.Dir("/home/user/.naturedopes-cli/config.json")
// Returns: "/home/user/.naturedopes-cli"

// Convert struct to pretty JSON
data, err := json.MarshalIndent(cfg, "", "  ")
//                                    ^^  ^^^^
//                                    |   2-space indent
//                                    no prefix

// Write file with permissions
os.WriteFile(path, data, 0644)
//                       ^^^^
//                       rw-r--r-- permissions
```

#### 3. Implement helper functions
- `Set(key, value string) error` - Update a config value and save
- `Get(key string) (string, error)` - Get a config value
- `Clear() error` - Delete the config file

#### 4. Create `cmd/config.go`
This will create the CLI commands:
- `naturedopes-cli config set api-url <url>`
- `naturedopes-cli config get api-key`
- `naturedopes-cli config list`
- `naturedopes-cli config clear`

---

## 🔜 Next Session: Phase 2 - Configuration Management (Continued)

### What You'll Build Next:
A config system so users don't have to type `--api-url` and `--api-key` every time!

### Commands You'll Create:
```bash
# Set configuration values
naturedopes-cli config set api-url https://naturedopesapi-production.up.railway.app
naturedopes-cli config set api-key 69e9b39d...

# Get configuration values
naturedopes-cli config get api-url
naturedopes-cli config get api-key

# List all config
naturedopes-cli config list

# Clear config
naturedopes-cli config clear
```

### Files You'll Create:
1. `pkg/config/config.go` - Configuration management package
2. `cmd/config.go` - Config command handlers

### What You'll Learn:
- Go structs and JSON tags
- File I/O (`os.ReadFile`, `os.WriteFile`)
- JSON marshaling/unmarshaling
- User home directory handling
- Error wrapping
- Subcommands in Cobra
- Cross-platform path handling

### How It Will Work:
```
User runs: naturedopes-cli config set api-key abc123
    ↓
Save to: ~/.naturedopes-cli/config.json
    ↓
{
  "api_url": "http://localhost:8080",
  "api_key": "abc123"
}
    ↓
Future commands automatically load this config!
```

---

## 📊 Overall Progress

```
Phase 1: Foundation              ████████████████████ 100% ✅
Phase 2: Configuration           ████████░░░░░░░░░░░░  40% 🚧
Phase 3: API Client              ░░░░░░░░░░░░░░░░░░░░   0%
Phase 4: Images Commands         ░░░░░░░░░░░░░░░░░░░░   0%
Phase 5: Search Functionality    ░░░░░░░░░░░░░░░░░░░░   0%
Phase 6: API Keys Commands       ░░░░░░░░░░░░░░░░░░░░   0%
Phase 7: Polish & Error Handling ░░░░░░░░░░░░░░░░░░░░   0%
Phase 8: Testing                 ░░░░░░░░░░░░░░░░░░░░   0%

Total Project: ███░░░░░░░░░░░░░░░ 17.5% Complete
```

---

## 🎯 Quick Start for Next Session

When you're ready to continue:

1. **Review**: Read BUILD_GUIDE.md Phase 2 section
2. **Say**: "I'm ready for Phase 2" or "Let's do configuration"
3. **We'll build**: The config system together

### Key Concepts Preview:
- **Structs**: Data structures (like TypeScript interfaces)
- **JSON tags**: How to map struct fields to JSON
- **File paths**: `~/.naturedopes-cli/config.json`
- **Marshaling**: Convert struct → JSON
- **Unmarshaling**: Convert JSON → struct

---

## 💡 Recap of Key Go Concepts

### 1. Packages
```go
package cmd  // This file belongs to "cmd" package
```

### 2. Imports
```go
import (
	"fmt"                    // Standard library
	"github.com/spf13/cobra" // External package
)
```

### 3. Variables
```go
var apiURL string           // Package-level variable
var rootCmd = &cobra.Command{...}  // Variable with initialization
```

### 4. The `&` Operator (Pointer)
```go
StringVar(&apiURL, ...)     // & = "address of" apiURL
```

### 5. The `init()` Function
```go
func init() {
	// Runs automatically when package loads
	// Use for setup/initialization
}
```

### 6. Error Handling
```go
err := rootCmd.Execute()
if err != nil {
	// Handle error
}
```

### 7. Public vs Private
```go
Execute()   // Capital = Public (exported)
init()      // Lowercase = Private (package only)
```

---

## 📝 Questions Answered This Session

### Q: What is `init()` for?
**A**: Special function that runs automatically when package loads. Used for setup/initialization like configuring flags.

### Q: What are flags?
**A**: Optional command-line arguments that modify behavior. Like `ls -la` where `-l` and `-a` are flags.

### Q: Why aren't flags showing in `--help`?
**A**: Cobra only shows flags section when there are subcommands. They still work! Try using them or wait until we add subcommands.

---

## 🚀 Commands to Remember

```bash
# Run the CLI
go run main.go

# Run with flags
go run main.go --api-url https://test.com

# Build executable
go build -o naturedopes-cli

# Run built executable
./naturedopes-cli --help

# Format your code
go fmt ./...

# Get help on any command
go help <command>
```

---

## 🎉 Achievements Unlocked

- ✅ Created your first Go CLI project
- ✅ Used Cobra framework
- ✅ Implemented command-line flags
- ✅ Understood `init()` function
- ✅ Learned about pointers
- ✅ Built working help system

---

## 📚 Resources Used

- [Cobra Documentation](https://cobra.dev/)
- [Go Tour](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)

---

**Great work today! See you in Phase 2! 🚀**

When you're ready to continue, just say: "Ready for Phase 2" or "Let's code configuration"
