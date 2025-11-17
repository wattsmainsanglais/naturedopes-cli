# Nature Dopes CLI - Progress Tracker

**Last Updated**: 2025-11-17
**Current Phase**: Phase 2 - Configuration Management (COMPLETED)
**Next Phase**: Phase 3 - API Client Foundation

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

## ✅ Phase 2: Configuration Management - COMPLETED

### What You Built:
1. ✅ Created `pkg/config/config.go` file
2. ✅ Defined `Config` struct with JSON tags
3. ✅ Implemented `getConfigFilePath()` helper function
4. ✅ Implemented `Load()` function
5. ✅ Implemented `Save()` function (method on *Config)
6. ✅ Implemented `Set()` function (using switch statement)
7. ✅ Implemented `Get()` function (returns value and error)
8. ✅ Created `cmd/config.go` - configCmd, setCmd, getCmd, listCmd
9. ✅ Used `reflect` package to iterate through struct fields
10. ✅ Extracted JSON tags for clean output formatting
11. ✅ Added `init()` function to wire commands together
12. ✅ Tested all commands successfully

### What You've Learned So Far:

#### 1. **Hidden Files and Directories**
- Files/directories starting with `.` are hidden in Unix/Linux/Mac
- View with `ls -a` (not just `ls`)
- Common examples: `.gitconfig`, `.ssh/`, `.bashrc`
- Used for config files to keep home directory clean

#### 2. **Error Wrapping in Go**
```go
// GOOD ✅ - Wrapping an error from another function
return "", fmt.Errorf("failed to get home directory: %w", err)
//                    ^lowercase  ^colon+space  ^%w wraps error

// GOOD ✅ - Creating a new validation error
return fmt.Errorf("invalid config key: %s", key)  // No %w, new error

// BAD ❌
return "", fmt.Errorf("Error %w", err)  // Not descriptive, capitalized
```
- Use `%w` to wrap errors from other functions (preserves error chain)
- Don't use `%w` for new validation errors you create
- Use lowercase messages (errors appear mid-sentence when chained)
- Add descriptive context with `: %w` pattern

#### 3. **File Paths in Go**
```go
filepath.Join(homeDir, ".naturedopes-cli", "config.json")
filepath.Dir(fullPath)  // Extract directory from full path
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

#### 7. **Methods vs Functions**
```go
// Method (has a receiver)
func (config *Config) Save() error {
    // Called like: config.Save()
}

// Regular function
func Load() (*Config, error) {
    // Called like: Load()
}
```
- Methods are bound to types, have receivers
- `(config *Config)` is the receiver (like `this` or `self`)
- Use pointer receivers (`*Config`) for efficiency and to modify data

#### 8. **Checking File Existence**
```go
if _, err := os.Stat(path); os.IsNotExist(err) {
    // File doesn't exist - return default config (not an error!)
}
```
- Missing config file on first run is NORMAL
- Return default values, don't treat as error

#### 9. **JSON Marshaling and Unmarshaling**
```go
// Unmarshal: JSON → struct
var config Config
err := json.Unmarshal(fileContent, &config)

// Marshal: struct → JSON
data, err := json.MarshalIndent(config, "", "  ")  // Pretty print with 2-space indent
```
- Unmarshal converts JSON bytes to Go struct
- Marshal converts Go struct to JSON bytes
- Pass pointer to struct (`&config`) for Unmarshal so it can be modified

#### 10. **Switch Statements**
```go
switch key {
case "api-url":
    currentConfig.ApiURL = value
case "api-key":
    currentConfig.ApiKey = value
default:
    return fmt.Errorf("invalid key: %s", key)
}
```
- Idiomatic Go way to handle multiple string comparisons
- Cleaner than multiple `if/else if` statements
- No need for `break` (unlike C/Java)

#### 11. **Multiple Return Values**
```go
func Get(key string) (string, error) {
    // Success: return value and nil error
    return currentConfig.ApiURL, nil

    // Error: return zero value and error
    return "", fmt.Errorf("invalid key: %s", key)
}
```
- Functions can return multiple values
- Common pattern: return result and error
- On error, return zero value (empty string, 0, nil, etc.) and error

#### 12. **Creating Directories**
```go
os.MkdirAll(configDir, 0755)  // Creates directory and all parents
```
- Like `mkdir -p` in bash
- Creates parent directories if needed
- `0755` = permissions (rwxr-xr-x)

#### 13. **Cobra Command Arguments**
```go
Run: func(cmd *cobra.Command, args []string) {
    key := args[0]    // First argument
    value := args[1]  // Second argument
}
Args: cobra.ExactArgs(2),  // Require exactly 2 arguments
```
- `args` slice contains command-line arguments
- `cobra.ExactArgs(n)` validates argument count

#### 14. **Printing Errors**
```go
// WRONG ❌
fmt.Errorf("error: %w", err)  // Creates error but doesn't print it!

// RIGHT ✅
fmt.Printf("Error: %v\n", err)  // Prints the error
return  // Exit the function
```
- `fmt.Errorf()` creates an error, doesn't print it
- `fmt.Printf()` prints to stdout
- Use `return` to exit early after errors

#### 15. **Pointers: `*` and `&` Operators**
```go
// & = "address of" (get memory address)
ptr := &variable

// * in type = "pointer to"
var ptr *string  // ptr is a pointer to a string

// * for dereferencing = "value at address"
value := *ptr  // Get the value that ptr points to
```
- `&variable` gets the memory address of variable
- `*Type` in declaration means "pointer to Type"
- `*pointer` gets the value at that memory address
- Use with reflect: `reflect.ValueOf(*currentConfig)` to dereference

#### 16. **Reflection (`reflect` package)**
```go
values := reflect.ValueOf(*currentConfig)  // Get reflection value
types := values.Type()                     // Get type information

for i := 0; i < values.NumField(); i++ {
    field := types.Field(i)           // Field metadata
    jsonTag := field.Tag.Get("json")  // Extract JSON tag
    value := values.Field(i)          // Field value
}
```
- Reflection allows inspecting struct fields at runtime
- `NumField()` returns number of fields in struct
- `Field(i).Tag.Get("json")` extracts JSON tag value
- More flexible but slower than direct field access

#### 17. **Multiple `init()` Functions**
```go
// In cmd/root.go
func init() {
    rootCmd.PersistentFlags().StringVar(...)
}

// In cmd/config.go (SEPARATE file, SAME package)
func init() {
    rootCmd.AddCommand(configCmd)
}
```
- Multiple `init()` functions are allowed in the same package
- Each file can have its own `init()`
- All `init()` functions run automatically when package loads
- Use for file-specific setup and initialization

### Code You've Written:

**pkg/config/config.go** (complete):
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
			ApiURL: "http://localhost:8080",
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

func (config *Config) Save() error {
	path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("couldn't get home directory: %w", err)
	}

	configDir := filepath.Dir(path)

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("couldn't marshal JSON: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write data to file: %w", err)
	}

	return nil
}

func Set(key, value string) error {
	currentConfig, err := Load()
	if err != nil {
		return fmt.Errorf("could not load config file: %w", err)
	}

	switch key {
	case "api-url":
		currentConfig.ApiURL = value
	case "api-key":
		currentConfig.ApiKey = value
	default:
		return fmt.Errorf("invalid key: %s", key)
	}

	err = currentConfig.Save()
	if err != nil {
		return fmt.Errorf("could not save config file: %w", err)
	}

	return nil
}

func Get(key string) (string, error) {
	currentConfig, err := Load()
	if err != nil {
		return "", fmt.Errorf("could not load config file: %w", err)
	}

	switch key {
	case "api-url":
		return currentConfig.ApiURL, nil
	case "api-key":
		return currentConfig.ApiKey, nil
	default:
		return "", fmt.Errorf("invalid key: %s", key)
	}
}
```

**cmd/config.go** (partial - in progress):
```go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(command *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		err := config.Set(key, value)
		if err != nil {
			fmt.Printf("could not set: %v\n", err)
			return
		}

		fmt.Printf("New %v has been set", key)
	},
}

// TODO: Add getCmd, listCmd (optional), and init() function
```

### Next Steps When You Return:

#### 1. Complete `cmd/config.go`

Add the `getCmd` command:
```go
var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		value, err := config.Get(key)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("%s: %s\n", key, value)
	},
}
```

Add a `listCmd` command (optional):
```go
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("api-url: %s\n", cfg.ApiURL)
		fmt.Printf("api-key: %s\n", cfg.ApiKey)
	},
}
```

Add the `init()` function to wire everything together:
```go
func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(getCmd)
	configCmd.AddCommand(listCmd)  // Optional
}
```

#### 2. Test your commands!

```bash
# Set values
go run main.go config set api-url https://api.example.com
go run main.go config set api-key abc123

# Get values
go run main.go config get api-url
go run main.go config get api-key

# List all
go run main.go config list

# Check the config file was created
ls -la ~/.naturedopes-cli/
cat ~/.naturedopes-cli/config.json
```

---

## 🔜 Next Session: Phase 3 - API Client Foundation

### What You'll Build Next:
An HTTP client to communicate with the Nature Dopes API!

### Files You'll Create:
1. `pkg/models/types.go` - Data structures for Image and ApiKey
2. `pkg/api/client.go` - HTTP client wrapper

### What You'll Learn:
- Go structs and JSON tags for API responses
- HTTP client usage (`net/http` package)
- Making GET/POST/DELETE requests
- Request/response handling
- Type definitions and pointer semantics
- Authentication headers

### Data Models You'll Define:
```go
type Image struct {
    ID         int     `json:"id"`
    SpeciesName string `json:"species_name"`
    GpsLong    float64 `json:"gps_long"`
    GpsLat     float64 `json:"gps_lat"`
    ImagePath  string  `json:"image_path"`
    UserID     int     `json:"user_id"`
}

type ApiKey struct {
    ID        int     `json:"id"`
    Key       string  `json:"key"`
    Name      string  `json:"name"`
    CreatedAt string  `json:"created_at"`
}
```

---

## 📊 Overall Progress

```
Phase 1: Foundation              ████████████████████ 100% ✅
Phase 2: Configuration           ████████████████████ 100% ✅
Phase 3: API Client              ░░░░░░░░░░░░░░░░░░░░   0%
Phase 4: Images Commands         ░░░░░░░░░░░░░░░░░░░░   0%
Phase 5: Search Functionality    ░░░░░░░░░░░░░░░░░░░░   0%
Phase 6: API Keys Commands       ░░░░░░░░░░░░░░░░░░░░   0%
Phase 7: Polish & Error Handling ░░░░░░░░░░░░░░░░░░░░   0%
Phase 8: Testing                 ░░░░░░░░░░░░░░░░░░░░   0%

Total Project: █████░░░░░░░░░░░░░ 25% Complete
```

---

## 🎯 Quick Start for Next Session

When you're ready to continue:

1. **Review**: Read BUILD_GUIDE.md Phase 3 section
2. **Say**: "I'm ready for Phase 3" or "Let's build the API client"
3. **We'll build**: HTTP client to communicate with the Nature Dopes API

### Key Concepts Preview:
- **HTTP requests**: Using `net/http` package
- **API models**: Structs for Image and ApiKey data
- **JSON unmarshaling**: Convert API responses to Go structs
- **Authentication**: Adding API key headers to requests
- **Error handling**: Handling HTTP errors gracefully

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
- ✅ Built complete configuration system
- ✅ Mastered pointers (`*` and `&`)
- ✅ Used reflection to inspect structs
- ✅ Extracted JSON tags dynamically
- ✅ Implemented subcommands with Cobra
- ✅ Understood multiple `init()` functions
- ✅ Built working help system

---

## 📚 Resources Used

- [Cobra Documentation](https://cobra.dev/)
- [Go Tour](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Go Reflection](https://go.dev/blog/laws-of-reflection)

---

**Great work! Phase 2 Complete! See you in Phase 3! 🚀**

When you're ready to continue, just say: "Ready for Phase 3" or "Let's build the API client"
