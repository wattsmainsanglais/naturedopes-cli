# Nature Dopes CLI - Build Guide

**Project**: Command-line tool for Nature Dopes API
**Language**: Go
**Framework**: Cobra CLI
**Learning Goal**: Hands-on Go development with guidance

## ⚠️ IMPORTANT: Your Role
**YOU write the code. Claude provides guidance, explanations, and reviews.**
- Claude will NOT write code for you (that defeats the learning purpose!)
- Claude will explain concepts, suggest structure, and help debug
- You learn by doing, not by watching Claude code

---

## 📋 Project Overview

Building a CLI tool to interact with the Nature Dopes Go API. You'll learn Go by writing the code yourself with step-by-step guidance.

### What You're Building:
```bash
naturedopes-cli images list
naturedopes-cli images get --id 1
naturedopes-cli images search --species "Oak"
naturedopes-cli keys list
naturedopes-cli keys create --name "Research Key"
naturedopes-cli config set api-key <key>
```

---

## 🎯 Learning Objectives

By the end of this project, you'll understand:
- ✅ Go project structure and organization
- ✅ Package management with `go.mod`
- ✅ Command-line parsing with Cobra
- ✅ HTTP client usage (`net/http`)
- ✅ JSON encoding/decoding
- ✅ Error handling patterns in Go
- ✅ File I/O for configuration
- ✅ Structs and interfaces
- ✅ Testing in Go

---

## 🏗️ Project Structure

```
naturedopes-cli/
├── main.go                 # Entry point - calls cmd.Execute()
├── go.mod                  # Dependencies
├── go.sum                  # Dependency checksums (auto-generated)
├── README.md               # Project documentation
├── BUILD_GUIDE.md          # This file
│
├── cmd/                    # Commands package
│   ├── root.go            # Root command + global flags
│   ├── images.go          # Image-related commands
│   ├── keys.go            # API key management commands
│   └── config.go          # Configuration commands
│
├── pkg/                    # Shared packages
│   ├── api/               # API client
│   │   ├── client.go      # HTTP client wrapper
│   │   ├── images.go      # Image-related API calls
│   │   └── keys.go        # Key-related API calls
│   │
│   ├── config/            # Configuration management
│   │   └── config.go      # Load/save user config
│   │
│   └── models/            # Data structures
│       └── types.go       # Shared types (Image, ApiKey, etc.)
│
└── tests/                  # Tests (future)
    └── ...
```

---

## 📚 Phase 1: Foundation (Session 1)

### Goals:
- Set up project structure ✅ (DONE)
- Create root command with help
- Test basic CLI works

### Files to Create:
1. `main.go` - Entry point ✅ (DONE)
2. `cmd/root.go` - Root command
3. Test it runs

### What You'll Learn:
- How Go modules work
- Cobra command structure
- Go package imports

---

## 📚 Phase 2: Configuration Management (Session 2)

### Goals:
- Create config file system
- Store API key and URL
- Learn file I/O in Go

### Files to Create:
1. `pkg/config/config.go` - Config struct and functions
2. `cmd/config.go` - Config commands

### What You'll Learn:
- Structs in Go
- JSON marshaling/unmarshaling
- File system operations (`os` package)
- Error handling patterns
- User home directory handling

### Commands You'll Build:
```bash
naturedopes-cli config set api-key <key>
naturedopes-cli config set api-url <url>
naturedopes-cli config get api-key
naturedopes-cli config list
```

---

## 📚 Phase 3: API Client Foundation (Session 3)

### Goals:
- Create HTTP client wrapper
- Define data models
- Test connection to API

### Files to Create:
1. `pkg/models/types.go` - Data structures
2. `pkg/api/client.go` - HTTP client wrapper

### What You'll Learn:
- Go structs and tags (JSON)
- HTTP client usage
- Request/response handling
- Type definitions
- Pointer vs value semantics

### Data Models:
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
    ID        int       `json:"id"`
    Key       string    `json:"key"`
    Name      string    `json:"name"`
    CreatedAt string    `json:"created_at"`
    ExpiresAt string    `json:"expires_at"`
    LastUsed  *string   `json:"last_used"`
    Revoked   bool      `json:"revoked"`
}
```

---

## 📚 Phase 4: Images Commands (Session 4)

### Goals:
- Implement image listing
- Implement image retrieval
- Format output nicely

### Files to Create:
1. `pkg/api/images.go` - Image API calls
2. `cmd/images.go` - Image commands

### What You'll Learn:
- HTTP GET requests
- JSON decoding
- Command flags and arguments
- Table formatting
- Error handling and user feedback

### Commands You'll Build:
```bash
naturedopes-cli images list [--limit 10]
naturedopes-cli images get --id <id>
```

---

## 📚 Phase 5: Images Search (Session 5)

### Goals:
- Add search functionality
- Handle query parameters
- Filter results

### Files to Modify:
1. `pkg/api/images.go` - Add search method
2. `cmd/images.go` - Add search command

### What You'll Learn:
- URL query parameters
- String filtering
- Working with slices
- Optional flags

### Commands You'll Build:
```bash
naturedopes-cli images search --species "Oak"
naturedopes-cli images search --user-id 5
```

---

## 📚 Phase 6: API Keys Commands (Session 6)

### Goals:
- List API keys
- Create new keys
- Revoke keys

### Files to Create:
1. `pkg/api/keys.go` - API key operations
2. `cmd/keys.go` - Key management commands

### What You'll Learn:
- HTTP POST requests
- HTTP DELETE requests
- Request body handling
- Authentication headers
- Confirmation prompts

### Commands You'll Build:
```bash
naturedopes-cli keys list
naturedopes-cli keys create --name "My Key"
naturedopes-cli keys revoke --id <id>
```

---

## 📚 Phase 7: Polish & Error Handling (Session 7)

### Goals:
- Improve error messages
- Add validation
- Better output formatting

### What You'll Learn:
- Go error wrapping
- Input validation
- Custom error types
- Output formatting (tables, colors)

### Improvements:
- Clear error messages
- Input validation before API calls
- Colored output (optional)
- Progress indicators

---

## 📚 Phase 8: Testing (Session 8)

### Goals:
- Write unit tests
- Test API client
- Mock HTTP responses

### What You'll Learn:
- Go testing framework
- Table-driven tests
- HTTP mocking
- Test organization

### Files to Create:
```
tests/
├── api_test.go
├── config_test.go
└── ...
```

---

## 🎓 Go Concepts You'll Master

### 1. **Project Structure**
- How to organize Go code
- Package visibility rules
- Import paths

### 2. **Data Types**
```go
// Basic types
string, int, float64, bool

// Composite types
struct, slice, map

// Pointers
*string (pointer to string)
&variable (address of variable)
```

### 3. **Error Handling**
```go
// Go's explicit error handling
result, err := doSomething()
if err != nil {
    return err
}
```

### 4. **HTTP Client**
```go
// Making HTTP requests
resp, err := http.Get(url)
defer resp.Body.Close()
```

### 5. **JSON Handling**
```go
// Decode JSON
json.NewDecoder(resp.Body).Decode(&data)

// Encode JSON
json.NewEncoder(file).Encode(data)
```

### 6. **File I/O**
```go
// Read file
data, err := os.ReadFile(path)

// Write file
err := os.WriteFile(path, data, 0644)
```

---

## 🛠️ Development Workflow

### For Each Session:
1. **Review** - Look at what we're building
2. **Plan** - Understand the structure
3. **Code** - You write the code with guidance
4. **Test** - Run and verify it works
5. **Debug** - Fix any issues together
6. **Commit** - Save your progress

### Testing Your Code:
```bash
# Run the CLI
go run main.go

# Build executable
go build -o naturedopes-cli

# Run built executable
./naturedopes-cli

# Run with arguments
./naturedopes-cli images list
```

### Common Go Commands:
```bash
# Install dependencies
go get <package>

# Format code
go fmt ./...

# Run tests
go test ./...

# Build
go build

# Run without building
go run main.go
```

---

## 📖 Go Resources

### Official Docs:
- https://go.dev/tour/ - Interactive Go tutorial
- https://go.dev/doc/effective_go - Best practices
- https://pkg.go.dev/ - Package documentation

### Cobra Docs:
- https://cobra.dev/ - Cobra user guide
- https://github.com/spf13/cobra - GitHub repo

### HTTP Client:
- https://pkg.go.dev/net/http - Standard library HTTP

---

## 🎯 Current Status

### ✅ Completed:
- [x] Project folder created
- [x] Go module initialized
- [x] Cobra installed
- [x] Directory structure created
- [x] `main.go` created

### 🚧 Next Steps:
- [ ] Create `cmd/root.go`
- [ ] Test basic CLI works
- [ ] Create config management

---

## 💡 Tips for Success

1. **Read Error Messages** - Go errors are descriptive, read them carefully
2. **Use `go fmt`** - Format your code automatically
3. **Check Types** - Go is strongly typed, pay attention to types
4. **Test Often** - Run your code frequently
5. **Ask Questions** - If something is unclear, ask!

---

## 🚀 Ready to Start?

Your first task is in **Phase 1: Foundation**. We'll create the root command together.

When you're ready, say "I'm ready" and I'll guide you through creating `cmd/root.go`!

---

**Last Updated**: 2025-11-07
**Your Guide**: Claude
**Your Role**: Developer (you write the code!)
