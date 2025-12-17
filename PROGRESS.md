# Nature Dopes CLI - Progress Tracker

**Last Updated**: 2025-12-17
**Current Phase**: Phase 7 - Polish & Error Handling (IN PROGRESS)
**Next Phase**: Phase 8 - Testing

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

## ✅ Phase 3: API Client Foundation - COMPLETED

### What You Built:
1. ✅ Created `pkg/models/types.go` - Data structures for API responses
2. ✅ Created `pkg/api/client.go` - HTTP client wrapper
3. ✅ Implemented `doRequest()` method for making HTTP requests
4. ✅ Added API key authentication header support

### Files Created:
```
naturedopes-cli/
└── pkg/
    ├── models/
    │   └── types.go        ✅  (Image and ApiKey structs)
    └── api/
        └── client.go       ✅  (Client struct and doRequest method)
```

### What You Learned:

#### 1. **Data Models with JSON Tags**
```go
type Image struct {
    ID          int     `json:"id"`
    SpeciesName string  `json:"species_name"`
    GpsLong     float64 `json:"gps_long"`
    GpsLat      float64 `json:"gps_lat"`
    ImagePath   string  `json:"image_path"`
    UserID      int     `json:"user_id"`
}
```
- JSON tags map Go fields to JSON keys from API
- Use backticks for struct tags

#### 2. **Pointers for Nullable Fields**
```go
LastUsed *string `json:"last_used"`  // Can be nil (like null in DB)
```
- `*string` = pointer to string (optional field)
- `nil` in Go = `null` in JSON/databases
- Same concept as `String?` in Prisma

#### 3. **Constructor Pattern**
```go
func NewClient(baseURL, apiKey string) *Client {
    return &Client{
        BaseUrl:    baseURL,
        APIKey:     apiKey,
        HTTPClient: &http.Client{},
    }
}
```
- Functions that create and return instances
- Returns pointer so struct can be modified
- Common Go pattern for initialization

#### 4. **HTTP Client Usage (`net/http`)**
```go
req, err := http.NewRequest(method, url, nil)
resp, err := c.HTTPClient.Do(req)
```
- Create requests with `http.NewRequest()`
- Send requests with `HTTPClient.Do()`
- Third parameter (`nil`) is request body (for POST/PUT)

#### 5. **Setting HTTP Headers**
```go
req.Header.Set("X-API-Key", c.APIKey)
```
- Add authentication headers to requests
- Common pattern for API authentication

#### 6. **`defer` for Cleanup**
```go
defer resp.Body.Close()
```
- `defer` runs at end of function
- Always close response bodies to prevent memory leaks
- Cleanup pattern in Go

#### 7. **Reading Response Bodies**
```go
body, err := io.ReadAll(resp.Body)
```
- Reads all bytes from HTTP response
- Returns `[]byte` (byte slice)
- Can be converted to string or unmarshaled to struct

#### 8. **HTTP Status Code Handling**
```go
if resp.StatusCode >= 400 {
    return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
}
```
- Check status codes: 200-299 = success, 400+ = error
- Return descriptive errors with status and message

#### 9. **Format Specifiers in Error Messages**
```go
%w  // Wrap an error (type: error)
%d  // Integer/number (type: int)
%s  // String (type: string)
%v  // Any value (generic)
```
- Use correct specifier for each type
- `%w` only for wrapping existing errors

### Code You Wrote:

**pkg/models/types.go**:
```go
package models

type Image struct {
	ID          int     `json:"id"`
	SpeciesName string  `json:"species_name"`
	GpsLong     float64 `json:"gps_long"`
	GpsLat      float64 `json:"gps_lat"`
	ImagePath   string  `json:"image_path"`
	UserID      int     `json:"user_id"`
}

type ApiKey struct {
	ID        int     `json:"id"`
	Key       string  `json:"key"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"created_at"`
	ExpiresAt string  `json:"expires_at"`
	LastUsed  *string `json:"last_used"`
	Revoked   bool    `json:"revoked"`
}
```

**pkg/api/client.go**:
```go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseUrl    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(BaseUrl string, APIKey string) *Client {
	return &Client{
		BaseUrl:    BaseUrl,
		APIKey:     APIKey,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) doRequest(method string, path string) ([]byte, error) {
	url := c.BaseUrl + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create http request err: %w", err)
	}

	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not send request err: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response, err: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Returned status code: %d , message: %s ", resp.StatusCode, string(body))
	}

	return body, nil
}
```

### Key Achievements:
- ✅ Created reusable HTTP client structure
- ✅ Understood pointer semantics for nullable values
- ✅ Learned HTTP request/response handling
- ✅ Implemented authentication with headers
- ✅ Added proper error handling and cleanup
- ✅ Made connection between Go pointers and database nulls (Prisma concept)

---

## ✅ Phase 4: Images Commands - COMPLETED

### What You Built:
1. ✅ Implemented `ListImages()` method in `pkg/api/images.go`
2. ✅ Implemented `GetImage(id int)` method in `pkg/api/images.go`
3. ✅ Created `cmd/images.go` with `imagesCmd`, `listImagesCmd`, and `getImageCmd`
4. ✅ Used positional arguments for simple ID parameter
5. ✅ Tested commands successfully with real API

### Files Modified/Created:
```
naturedopes-cli/
├── pkg/api/
│   └── images.go       ✅  (ListImages, GetImage methods)
└── cmd/
    └── images.go       ✅  (images, list, get commands)
```

### What You Learned:

#### 1. **Positional Arguments vs Flags**
```go
// Positional (what you used - cleaner for single required values)
Args: cobra.ExactArgs(1),
id := args[0]
// Usage: naturedopes-cli images get 5

// vs Flags (better for multiple optional parameters)
getImageCmd.Flags().IntVar(&imageID, "id", 0, "Image ID")
// Usage: naturedopes-cli images get --id 5
```
- Positional arguments are simpler for single required values
- Flags are better for optional or multiple parameters

#### 2. **Error Handling in Cobra Commands**
```go
Run: func(command *cobra.Command, args []string) {
    // This function doesn't return errors!
    resp, err := client.ListImages()
    if err != nil {
        fmt.Printf("Error: %v\n", err)  // Print, not return
        return  // Exit function
    }
    // Safe to use resp here
}
```
- Cobra's `Run` function has no return type
- Use `fmt.Printf()` to **print** errors (not `fmt.Errorf()` which **creates** them)
- Always `return` after printing errors to stop execution

#### 3. **fmt.Errorf() vs fmt.Printf()**
```go
// WRONG in Cobra Run ❌
fmt.Errorf("error: %w", err)  // Creates error but discards it

// RIGHT ✅
fmt.Printf("Error: %v\n", err)  // Prints to user
return  // Stop execution
```
- `fmt.Errorf()` = creates errors for **returning** from functions
- `fmt.Printf()` = prints to stdout for **showing** users
- Use `fmt.Errorf()` in regular functions that return `error`
- Use `fmt.Printf()` in Cobra `Run` functions

#### 4. **Newlines in Output**
```go
// Two ways to add newlines:
fmt.Printf("message\n")   // Manual newline with \n
fmt.Println("message")    // Automatic newline
```
- Always add `\n` to `Printf()` calls or use `Println()`
- Prevents output from appearing on same line

#### 5. **String to Integer Conversion**
```go
id := args[0]  // String from command line
integer, err := strconv.Atoi(id)  // Convert to int
if err != nil {
    fmt.Printf("Error: invalid ID, must be a number\n")
    return
}
// Now safe to use integer
```
- Command line arguments are always strings
- Use `strconv.Atoi()` to convert to int
- Always check for conversion errors

#### 6. **JSON Unmarshaling**
```go
// In pkg/api/images.go
var images []models.Image  // For list (slice)
json.Unmarshal(resp, &images)

var image models.Image  // For single item
json.Unmarshal(resp, &image)
return &image, nil  // Return pointer
```
- Unmarshal JSON bytes into Go structs
- Use slice `[]Type` for arrays
- Use single `Type` for objects
- Pass pointer with `&` so Unmarshal can modify it

#### 7. **API Rate Limiting Understanding**
- Learned the API has two rate limits:
  - **Per-IP**: 1000 requests/day
  - **Per-API-Key**: 100 requests/hour (stricter)
- CLI requests count toward both limits
- Error handling already catches rate limit errors (429 status)

### Commands You Built:
```bash
# List all images
naturedopes-cli images list

# Get specific image
naturedopes-cli images get 5

# Help
naturedopes-cli images --help
```

### Code You Wrote:

**pkg/api/images.go** (complete):
```go
package api

import (
	"encoding/json"
	"fmt"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/models"
)

func (c *Client) ListImages() ([]models.Image, error) {
	var images []models.Image

	resp, err := c.doRequest("GET", "/images")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve images: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall to json: %w", err)
	}
	return images, nil
}

func (c *Client) GetImage(id int) (*models.Image, error) {
	var images models.Image

	resp, err := c.doRequest("GET", fmt.Sprintf("/images/%d", id))
	if err != nil {
		return nil, fmt.Errorf("Could not obtain image: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall to json: %w", err)
	}

	return &images, nil
}
```

**cmd/images.go** (complete):
```go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/api"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/config"
	"strconv"
)

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Get Images command",
}

var listImagesCmd = &cobra.Command{
	Use:   "list",
	Short: "Get list of images",
	Args:  cobra.ExactArgs(0),
	Run: func(command *cobra.Command, args []string) {
		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		resp, err := client.ListImages()
		if err != nil {
			fmt.Printf("could not retrieve images: %v\n", err)
			return
		}

		for _, image := range resp {
			fmt.Printf("name: %s, gps_long: %f, gps_lat: %f, image_path: %s\n",
				image.SpeciesName, image.GpsLong, image.GpsLat, image.ImagePath)
		}
	},
}

var getImageCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get individual image",
	Args:  cobra.ExactArgs(1),
	Run: func(command *cobra.Command, args []string) {
		id := args[0]
		integer, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Error, invalid ID, please check you've supplied an integer as argument: %v\n", err)
			return
		}

		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		image, err := client.GetImage(integer)
		if err != nil {
			fmt.Printf("could not retrieve image data: %v\n", err)
			return
		}

		fmt.Printf("id:%d name: %s, gps_long: %f, gps_lat: %f, image_path: %s\n",
			image.ID, image.SpeciesName, image.GpsLong, image.GpsLat, image.ImagePath)
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.AddCommand(listImagesCmd)
	imagesCmd.AddCommand(getImageCmd)
}
```

### Testing Results:
```bash
# Successfully tested with real API
✅ naturedopes-cli images list - shows all images
✅ naturedopes-cli images get 5 - shows single image
✅ naturedopes-cli images get abc - shows error for invalid ID
✅ API key authentication working
```

### Key Achievements:
- ✅ Built working image listing and retrieval
- ✅ Understood positional arguments vs flags trade-offs
- ✅ Mastered error handling in Cobra commands
- ✅ Learned difference between creating and printing errors
- ✅ Implemented proper input validation
- ✅ Successfully integrated with real API
- ✅ Understood rate limiting implications

---

## ✅ Phase 5: Search Functionality - COMPLETED

### What You Built:
1. ✅ Implemented `SearchImages()` method in `pkg/api/images.go`
2. ✅ Added URL query parameter support for filtering
3. ✅ Created `searchCmd` in `cmd/images.go` with positional arguments
4. ✅ Tested and verified CLI correctly sends query parameters

### ⚠️ API Limitation - TODO:
**The CLI search functionality is complete**, but the `naturedopesApi` backend doesn't support query parameter filtering yet!

**Current API behavior**: `GET /images?species_name=Globethistle` returns ALL images (ignores parameters)

**What needs to be added to the API**:
1. Modify `getImagesHandler` in `routes.go` to read query parameters (`r.URL.Query().Get()`)
2. Update `GetImages()` in `endpoints/image.go` to accept `species_name` and `user_id` parameters
3. Modify SQL query to use `WHERE` clauses for filtering

**Files to modify in naturedopesApi**:
- `/home/andrew/Code/2025/go/naturedopesApi/routes.go`
- `/home/andrew/Code/2025/go/naturedopesApi/endpoints/image.go`

**For now**: The CLI is ready and working correctly. API enhancement can be added later as a separate learning task.

### Files Modified:
```
naturedopes-cli/
├── pkg/api/
│   └── images.go       ✅  (Added SearchImages method)
└── cmd/
    └── images.go       ✅  (Added searchCmd)
```

### What You Learned:

#### 1. **URL Query Parameters**
```go
params := url.Values{}
params.Add("species", "Oak")
params.Add("user_id", "5")
queryString := params.Encode()  // "species=Oak&user_id=5"
```
- `url.Values{}` creates a map for query parameters
- `Add()` adds key-value pairs
- `Encode()` converts to URL-encoded string
- Automatically handles URL encoding (spaces, special chars)

#### 2. **Building Dynamic URLs**
```go
path := "/images"
if len(params) > 0 {
    path = path + "?" + params.Encode()
}
// Result: "/images?species=Oak&user_id=5"
```
- Start with base path
- Only add `?` and query string if there are parameters
- `len(params)` checks if any parameters were added

#### 3. **Optional Function Parameters in Go**
```go
func (c *Client) SearchImages(species string, userID int) ([]models.Image, error) {
    if species != "" {
        // Only add if provided
    }
    if userID > 0 {
        // Only add if provided
    }
}
```
- Go doesn't have optional parameters (unlike JavaScript/Python)
- Use zero values to indicate "not provided": `""` for string, `0` for int
- Check for zero values before using them
- Alternative: Use pointers (`*string`, `*int`) where `nil` means not provided

#### 4. **Converting Integers to Strings**
```go
userIDStr := strconv.Itoa(userID)  // int to string
// Example: 5 -> "5"
```
- `strconv.Itoa()` = "integer to ASCII"
- Needed for URL query parameters (must be strings)
- Opposite of `strconv.Atoi()` which you used before

#### 5. **Optional Flags in Cobra**
```go
var species string
var userID int

searchCmd.Flags().StringVar(&species, "species", "", "Filter by species name")
searchCmd.Flags().IntVar(&userID, "user-id", 0, "Filter by user ID")
```
- `.Flags()` creates command-specific flags (not persistent)
- Default values: `""` for string, `0` for int
- If user doesn't provide flag, variable has default value
- Both filters are optional - can use one, both, or neither

#### 6. **Difference Between `.Flags()` and `.PersistentFlags()`**
```go
// Regular flags - only for this command
searchCmd.Flags().StringVar(...)

// Persistent flags - inherited by all subcommands
rootCmd.PersistentFlags().StringVar(...)
```
- `.Flags()` = local to one command
- `.PersistentFlags()` = available to command and all children
- Use persistent for global settings (api-url, api-key)
- Use regular for command-specific options (species, user-id)

### Commands You Built:
```bash
# Search with species filter
naturedopes-cli images search --species Oak

# Search with user-id filter
naturedopes-cli images search --user-id 1

# Search with both filters
naturedopes-cli images search --species Oak --user-id 1

# Search with no filters (returns all)
naturedopes-cli images search

# Get help
naturedopes-cli images search --help
```

### Code You Wrote:

**pkg/api/images.go** (added SearchImages method):
```go
func (c *Client) SearchImages(species string, userID int) ([]models.Image, error) {
	var images []models.Image

	// Start with base path
	path := "/images"

	// Build query parameters if provided
	params := url.Values{}
	if species != "" {
		params.Add("species", species)
	}
	if userID > 0 {
		params.Add("user_id", strconv.Itoa(userID))
	}

	// Add query string to path if we have parameters
	if len(params) > 0 {
		path = path + "?" + params.Encode()
	}

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, fmt.Errorf("could not search images: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal to json: %w", err)
	}

	return images, nil
}
```

**cmd/images.go** (added searchCmd):
```go
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search images with optional filters",
	Run: func(cmd *cobra.Command, args []string) {
		species, _ := cmd.Flags().GetString("species")
		userID, _ := cmd.Flags().GetInt("user-id")

		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		images, err := client.SearchImages(species, userID)
		if err != nil {
			fmt.Printf("Error searching images: %v\n", err)
			return
		}

		if len(images) == 0 {
			fmt.Println("No images found")
			return
		}

		for _, image := range images {
			fmt.Printf("id: %d, name: %s, gps_long: %f, gps_lat: %f, image_path: %s\n",
				image.ID, image.SpeciesName, image.GpsLong, image.GpsLat, image.ImagePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.AddCommand(listImagesCmd)
	imagesCmd.AddCommand(getImageCmd)
	imagesCmd.AddCommand(searchCmd)

	searchCmd.Flags().String("species", "", "Filter by species name")
	searchCmd.Flags().Int("user-id", 0, "Filter by user ID")
}
```

### Testing (TODO for next session):
```bash
# Test different scenarios:
go run main.go images search --species Oak
go run main.go images search --user-id 1
go run main.go images search --species Oak --user-id 1
go run main.go images search
go run main.go images search --help
```

### Key Achievements:
- ✅ Learned URL query parameter construction
- ✅ Understood optional parameters in Go
- ✅ Used `url.Values{}` for building query strings
- ✅ Implemented command-specific flags in Cobra
- ✅ Learned difference between `.Flags()` and `.PersistentFlags()`
- ✅ Built flexible search with multiple optional filters
- ✅ Used `strconv.Itoa()` for integer to string conversion

---

## ✅ Phase 6: API Keys Commands - COMPLETED

### What You Built:
1. ✅ Modified `pkg/api/client.go` - Added request body support to `doRequest()`
2. ✅ Created `pkg/api/keys.go` - Implemented all three API methods
3. ✅ Created `cmd/keys.go` - Built CLI commands for key management
4. ✅ Documented API security improvements needed

### Files Created/Modified:
```
naturedopes-cli/
├── pkg/api/
│   ├── client.go                    ✅ Modified - Added body parameter to doRequest()
│   └── keys.go                      ✅ Created - GenerateKey, ListKeys, RevokeKey methods
├── cmd/
│   └── keys.go                      ✅ Created - list, generate, revoke commands
└── API_SECURITY_IMPROVEMENTS.md     ✅ Created - Documentation for API hardening
```

### What You've Learned:

#### 1. **POST Requests with JSON Body**
```go
// Create request body struct
requestBody := struct {
    Name string `json:"name"`
}{
    Name: name,
}

// Marshal to JSON
jsonData, err := json.Marshal(requestBody)

// Send in POST request
resp, err := client.doRequest("POST", "/api/keys", jsonData)
```
- Anonymous structs for one-time request bodies
- `json.Marshal()` converts Go struct → JSON bytes
- Request body passed as `[]byte` parameter

#### 2. **Modifying Existing Functions**
```go
// Before:
func (c *Client) doRequest(method string, path string) ([]byte, error)

// After:
func (c *Client) doRequest(method string, path string, body []byte) ([]byte, error)
```
- Added optional body parameter for POST/PUT requests
- Used `io.Reader` conversion: `bytes.NewBuffer(body)`
- Set `Content-Type: application/json` header when body exists

#### 3. **DELETE Requests**
```go
func (client *Client) RevokeKey(id int) error {
    _, err := client.doRequest("DELETE", fmt.Sprintf("/api/keys/%d", id), nil)
    return err
}
```
- DELETE often returns no data (204 No Content)
- Return only `error`, not data
- Use `_` to discard unused response body

#### 4. **Pointers vs Slices**
```go
// Single item - return pointer
func GenerateKey(name string) (*models.ApiKey, error)  // Returns *ApiKey

// Multiple items - return slice (NOT pointer to slice)
func ListKeys() ([]models.ApiKey, error)  // Returns []ApiKey, not *[]
```
**Why?**
- Structs are value types → use pointers to avoid copying
- Slices are reference types → already contain pointer internally

#### 5. **Separation of Concerns**
- **API layer** (`pkg/api/`): HTTP requests, return data/errors only
- **Command layer** (`cmd/`): User messages, formatting, interaction
- Don't print user messages from API methods!

### Code You Wrote:

**pkg/api/client.go** (modified doRequest):
```go
func (c *Client) doRequest(method string, path string, body []byte) ([]byte, error) {
    url := c.BaseUrl + path
    var reqBody io.Reader = nil
    if body != nil {
        reqBody = bytes.NewBuffer(body)
    }

    req, err := http.NewRequest(method, url, reqBody)
    if err != nil {
        return nil, fmt.Errorf("could not create http request err: %w", err)
    }

    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    if c.APIKey != "" {
        req.Header.Set("X-API-Key", c.APIKey)
    }

    // ... rest of request handling
}
```

**pkg/api/keys.go** (complete):
```go
func (client *Client) GenerateKey(name string) (*models.ApiKey, error) {
    var apiKey models.ApiKey

    requestBody := struct {
        Name string `json:"name"`
    }{
        Name: name,
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("could not create jsonData: %w", err)
    }

    resp, err := client.doRequest("POST", "/api/keys", jsonData)
    if err != nil {
        return nil, fmt.Errorf("could not create api keys from naturedopesApi: %w", err)
    }

    err = json.Unmarshal(resp, &apiKey)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshal response: %w", err)
    }

    return &apiKey, nil
}

func (client *Client) ListKeys() ([]models.ApiKey, error) {
    var apiKeys []models.ApiKey

    resp, err := client.doRequest("GET", "/api/keys", nil)
    if err != nil {
        return nil, fmt.Errorf("could not get apikeys: %w", err)
    }

    err = json.Unmarshal(resp, &apiKeys)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshall json: %w", err)
    }

    return apiKeys, nil
}

func (client *Client) RevokeKey(id int) error {
    _, err := client.doRequest("DELETE", fmt.Sprintf("/api/keys/%d", id), nil)
    if err != nil {
        return fmt.Errorf("could not delete api-key: %w", err)
    }

    return nil
}
```

#### 6. **API Authentication & The Bootstrap Problem**
```go
// API key endpoints are currently unprotected
// This allows first-time users to create their first key
// But creates security issues:
// - Anyone can list ALL keys
// - Anyone can revoke ANY key
```

**The Bootstrap Problem:**
- How do you authenticate to create an API key if you need an API key to authenticate?

**Current Solution:**
- `POST /api/keys` is unprotected (anyone can create)
- `GET /api/keys` is unprotected (security issue - shows all keys)
- `DELETE /api/keys/{id}` is unprotected (security issue - can revoke any key)

**Better Solution (documented in API_SECURITY_IMPROVEMENTS.md):**
- Add `user_id` to api_keys table
- Users can only see/revoke their own keys
- `POST /api/keys` remains unprotected for bootstrapping

### Commands You Built:

```bash
# List all API keys
naturedopes-cli keys list

# Generate new API key
naturedopes-cli keys generate "My Research Key"

# Revoke an API key by ID
naturedopes-cli keys revoke 5
```

### Code You Wrote:

**cmd/keys.go** (complete):
```go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/api"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/config"
	"strconv"
)

var keysCmnd = &cobra.Command{
	Use:   "keys",
	Short: "For api key management",
}

var listKeys = &cobra.Command{
	Use:   "list",
	Short: "List api keys",
	Args:  cobra.ExactArgs(0),
	Run: func(command *cobra.Command, args []string) {
		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		resp, error := client.ListKeys()
		if error != nil {
			fmt.Printf("could not get api keys: %v", error)
			return
		}

		for _, k := range resp {
			fmt.Printf("id: %v , name: %v, created: %v, expires: %v, last used: %v\n",
				k.ID, k.Name, k.CreatedAt, k.ExpiresAt, k.LastUsed)
		}
	},
}

var generateKey = &cobra.Command{
	Use:   "generate <name>",
	Short: "Create new api key",
	Args:  cobra.ExactArgs(1),
	Run: func(command *cobra.Command, args []string) {
		name := args[0]

		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		resp, error := client.GenerateKey(name)
		if error != nil {
			fmt.Printf("could not generate api key: %v", error)
			return
		}

		fmt.Printf("api key %v generated, key value: %v , please write this down. key will expire %v,",
			resp.Name, resp.Key, resp.ExpiresAt)
	},
}

var revokeKey = &cobra.Command{
	Use:   "revoke <id>",
	Short: "revoke api key by id",
	Args:  cobra.ExactArgs(1),
	Run: func(command *cobra.Command, args []string) {
		id := args[0]

		integer, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Error, invalid ID, please check you've supplied an integer as argument: %v\n", err)
			return
		}

		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		error := client.RevokeKey(integer)
		if error != nil {
			fmt.Printf("could not delete api key, %v", error)
			return
		}

		fmt.Printf("api key of id %v , has been successfully removed", integer)
	},
}

func init() {
	rootCmd.AddCommand(keysCmnd)
	keysCmnd.AddCommand(listKeys)
	keysCmnd.AddCommand(generateKey)
	keysCmnd.AddCommand(revokeKey)
}
```

### Key Achievements:
- ✅ Learned POST requests with JSON body
- ✅ Modified existing function to support request bodies
- ✅ Understood DELETE requests
- ✅ Mastered pointer vs slice return types
- ✅ Applied separation of concerns principle
- ✅ Successfully built complete API key management system
- ✅ Understood the bootstrap problem in API authentication
- ✅ Documented comprehensive security improvements for the API
- ✅ Built three working commands: list, generate, revoke

### Phase 6 Testing Results - 2025-12-17

**Testing Environment**: Live API on Railway (naturedopesapi-production.up.railway.app)

#### ✅ Tests Passed:
1. **Image Commands**:
   - `images list` - Successfully retrieves all images from live API
   - `images get <id>` - Successfully retrieves individual images
   - `images search` - Sends query parameters correctly (API filtering not yet implemented)

2. **API Key Commands**:
   - `keys generate <name>` - Successfully creates new API keys
   - `keys list` - Successfully lists all API keys
   - `keys revoke` - Successfully revokes keys (soft delete with revoked=true)

3. **Error Handling**:
   - Invalid image IDs are caught and reported
   - Invalid ID formats (non-integers) are validated

#### 🐛 Issue Discovered:
- **Missing Revoked Status Display**: The `keys list` command didn't show the `revoked` field, making it impossible to see which keys were revoked. (Fixed in Phase 7)

#### 📝 Notes:
- CLI successfully updated to match recent API changes:
  - ✅ `ListKeys()` now calls `/api/keys/list` (not `/api/keys`)
  - ✅ `RevokeKey()` now uses X-API-Key header (not path parameter)
  - ✅ New `GetKeyInfo()` function added for `/api/keys/get` endpoint
- API uses soft delete for keys (sets revoked=true, keeps record) - good design!
- All endpoints work correctly with the deployed Railway API

---

## 🚧 Phase 7: Polish & Error Handling - IN PROGRESS

**Started**: 2025-12-17
**Status**: 40% Complete

### Goals:
- ✅ Fix missing revoked field display
- ✅ Add confirmation prompts for destructive actions
- 🚧 Improve error messages with actionable guidance
- ⏳ Add input validation before API calls
- ⏳ Better output formatting

### What You've Built:

#### 1. ✅ Fixed Revoked Field Display
**File Modified**: `cmd/keys.go` (line 31)

**Problem**: Keys list didn't show revoked status, so users couldn't tell which keys were revoked.

**Solution**: Added `Revoked` field to output:
```go
fmt.Printf("id: %v , name: %v, key: %v, created: %v, expires: %v, last used: %v, revoked: %v\n",
    k.ID, k.Name, k.Key[:8], k.CreatedAt, k.ExpiresAt, k.LastUsed, k.Revoked)
```

**Result**: Users can now see `revoked: true` or `revoked: false` for each key.

#### 2. ✅ Added Confirmation Prompt for Revoke
**File Modified**: `cmd/keys.go` (revokeKey command)

**Problem**: Revoking a key is destructive and permanent - easy to do accidentally.

**Solution**: Added interactive confirmation prompt:
```go
fmt.Print("Are you sure you want to revoke your API key? This cannot be undone. (yes/no): ")

reader := bufio.NewReader(os.Stdin)
response, _ := reader.ReadString('\n')
response = strings.TrimSpace(strings.ToLower(response))

if response != "yes" {
    fmt.Println("Revoke cancelled.")
    return
}
```

**Result**: Users must type "yes" to confirm revocation, preventing accidental key deletion.

### What You've Learned (Phase 7):

#### 1. **Soft Delete Pattern**
```go
// Instead of deleting records:
DELETE FROM api_keys WHERE id = 1

// Mark as deleted (keeps history):
UPDATE api_keys SET revoked = true WHERE id = 1
```
- Maintains audit trail
- Allows recovery if needed
- Standard practice for important data

#### 2. **Buffered I/O (`bufio` package)**
```go
import "bufio"

reader := bufio.NewReader(os.Stdin)  // Create buffered reader
text, _ := reader.ReadString('\n')   // Read until Enter
text = strings.TrimSpace(text)       // Remove whitespace
```
- **Buffer** = Temporary storage area in memory
- Makes I/O operations more efficient (batch reads vs single character reads)
- Like using a shopping cart instead of carrying items one by one

#### 3. **Interactive User Prompts**
```go
fmt.Print("Question? (yes/no): ")     // Ask question
reader := bufio.NewReader(os.Stdin)   // Get ready to read
response, _ := reader.ReadString('\n') // Program PAUSES here waiting for input
```
- Program execution stops until user presses Enter
- Good for confirmations of destructive actions
- Improves user experience and prevents mistakes

#### 4. **String Manipulation for Input**
```go
response = strings.TrimSpace(response)  // Remove whitespace and \n
response = strings.ToLower(response)     // Convert to lowercase
// Now "YES", "Yes", "yes" all become "yes"
```
- Makes input comparison more flexible
- User-friendly (case-insensitive)

### Next Steps (When You Return):

#### 3. Improve Error Messages
Add helpful error messages when config is missing:
- Check if API key is set before making requests
- Provide actionable guidance: "Run: naturedopes-cli keys generate <name>"
- Better context in error messages

#### 4. Input Validation
- Validate URLs before saving to config
- Check ID values are positive integers
- Validate API key format

#### 5. Better Output Formatting
- Align columns in list outputs
- Consider using a table library
- Color coding (optional)

---

## 📊 Overall Progress

```
Phase 1: Foundation              ████████████████████ 100% ✅
Phase 2: Configuration           ████████████████████ 100% ✅
Phase 3: API Client              ████████████████████ 100% ✅
Phase 4: Images Commands         ████████████████████ 100% ✅
Phase 5: Search Functionality    ████████████████████ 100% ✅
Phase 6: API Keys Commands       ████████████████████ 100% ✅
Phase 7: Polish & Error Handling ████████░░░░░░░░░░░░  40%
Phase 8: Testing                 ░░░░░░░░░░░░░░░░░░░░   0%

Total Project: ████████████████░ 80% Complete
```

---

## 🎯 Quick Start for Next Session

**Current Status**: Phase 7 is 40% complete (2 of 5 tasks done)

When you're ready to continue Phase 7:

### Remaining Tasks:
1. **Add helpful error messages** - Check if API key is configured before API calls
2. **Add input validation** - Validate URLs, IDs, and other inputs before use
3. **Improve output formatting** - Better table alignment and visual clarity

### What You've Already Completed:
- ✅ Fixed revoked field display in keys list
- ✅ Added confirmation prompt for key revocation

### To Resume:
Say "Let's continue Phase 7" or "Ready to add error messages"

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
