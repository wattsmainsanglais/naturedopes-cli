# Advanced Go Patterns for Nature Dopes CLI

**Created**: 2026-01-08
**Purpose**: Document advanced Go patterns to level up from junior to mid-level developer
**Projects**: naturedopes-cli & naturedopesApi

---

## Overview

This guide covers 7 advanced Go patterns that will take your projects to the next level. Each pattern includes:
- What it is and why it matters
- Where to apply it in your existing projects
- Practical code examples
- Common pitfalls

**Current Skill Level**: Junior-to-mid developer with solid fundamentals
**Target Skill Level**: Mid-level developer with production-ready patterns

---

## Pattern 1: Context (Most Important for Production Code)

### What is Context?

Context is a way to carry deadlines, cancellation signals, and request-scoped values across API boundaries and between processes.

**Think of it like**: A control panel that travels with your request, allowing you to:
- Set timeouts ("this request should finish in 5 seconds max")
- Cancel operations ("user clicked cancel, stop everything")
- Pass request-scoped data ("which user is making this request?")

### Why It Matters

**Problem Without Context:**
```go
// What if the API never responds?
resp, err := http.Get("https://slow-api.com/data")
// Your CLI hangs forever, user can't even Ctrl+C
```

**Solution With Context:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, "GET", "https://slow-api.com/data", nil)
resp, err := client.Do(req)
// Automatically stops after 5 seconds with error
```

### The Context Interface

```go
type Context interface {
    // Deadline returns when work should be canceled
    Deadline() (deadline time.Time, ok bool)

    // Done returns a channel that's closed when work should stop
    Done() <-chan struct{}

    // Err returns why this context was canceled
    Err() error

    // Value returns data associated with this context
    Value(key interface{}) interface{}
}
```

**You don't create contexts by implementing this interface.** Instead, use the standard library functions.

---

### Context Types

#### 1. `context.Background()`
The root context - use when starting a new chain.

```go
// In main functions, top-level requests, tests
ctx := context.Background()
```

#### 2. `context.TODO()`
Placeholder when you're not sure what context to use yet.

```go
// Temporary - replace with proper context later
ctx := context.TODO()
```

#### 3. `context.WithTimeout()`
Automatically cancels after a duration.

```go
// Timeout after 10 seconds
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel() // ALWAYS defer cancel to avoid leaks

// Use ctx in your operations...
```

#### 4. `context.WithDeadline()`
Cancels at a specific time.

```go
// Cancel at specific time
deadline := time.Now().Add(1 * time.Hour)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

#### 5. `context.WithCancel()`
Manual cancellation (for user-triggered stops).

```go
ctx, cancel := context.WithCancel(context.Background())

// Later, when user clicks "Cancel" button:
cancel()
```

#### 6. `context.WithValue()`
Pass request-scoped data (use sparingly!).

```go
// Pass user ID through context
ctx := context.WithValue(context.Background(), "userID", 123)

// Retrieve later
userID := ctx.Value("userID").(int)
```

---

### Context Rules (IMPORTANT!)

#### Rule 1: Always Pass Context as First Parameter
```go
// CORRECT ✅
func doRequest(ctx context.Context, url string) error

// WRONG ❌
func doRequest(url string, ctx context.Context) error
```

#### Rule 2: Never Store Context in a Struct
```go
// WRONG ❌
type Client struct {
    ctx context.Context  // Don't do this!
}

// CORRECT ✅
type Client struct {
    baseURL string
}

func (c *Client) GetData(ctx context.Context) error {
    // Pass context as parameter
}
```

**Why?** Contexts are request-scoped, not object-scoped. Storing them causes lifetime issues.

#### Rule 3: Always defer cancel()
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // MUST do this, even if context expires naturally

// Prevents resource leaks
```

#### Rule 4: Check for Cancellation in Long Operations
```go
func processLargeDataset(ctx context.Context, data []Item) error {
    for _, item := range data {
        // Check if context was cancelled
        select {
        case <-ctx.Done():
            return ctx.Err()  // Stop processing
        default:
            // Continue processing
        }

        processItem(item)
    }
    return nil
}
```

---

### Applying Context to naturedopes-cli

#### Step 1: Update `doRequest()` Method

**File**: `pkg/api/client.go`

**Current code:**
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

    // ... rest of function
}
```

**Updated with context:**
```go
func (c *Client) doRequest(ctx context.Context, method string, path string, body []byte) ([]byte, error) {
    url := c.BaseUrl + path
    var reqBody io.Reader = nil
    if body != nil {
        reqBody = bytes.NewBuffer(body)
    }

    // Use NewRequestWithContext instead of NewRequest
    req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
    if err != nil {
        return nil, fmt.Errorf("could not create http request err: %w", err)
    }

    // Set headers...
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    if c.APIKey != "" {
        req.Header.Set("X-API-Key", c.APIKey)
    }

    // Request now respects context timeout/cancellation
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        // Check if error was due to context cancellation
        if ctx.Err() != nil {
            return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
        }
        return nil, fmt.Errorf("could not send request: %w", err)
    }

    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("returned status code: %d, message: %s", resp.StatusCode, string(body))
    }

    return body, nil
}
```

**Key changes:**
1. Added `ctx context.Context` as first parameter
2. Changed `http.NewRequest()` to `http.NewRequestWithContext(ctx, ...)`
3. Added check for context cancellation in error handling

---

#### Step 2: Update All API Methods

**File**: `pkg/api/images.go`

**Before:**
```go
func (c *Client) ListImages() ([]models.Image, error) {
    var images []models.Image

    resp, err := c.doRequest("GET", "/images", nil)
    if err != nil {
        return nil, fmt.Errorf("could not retrieve images: %w", err)
    }

    err = json.Unmarshal(resp, &images)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshall to json: %w", err)
    }
    return images, nil
}
```

**After:**
```go
func (c *Client) ListImages(ctx context.Context) ([]models.Image, error) {
    var images []models.Image

    // Pass context to doRequest
    resp, err := c.doRequest(ctx, "GET", "/images", nil)
    if err != nil {
        return nil, fmt.Errorf("could not retrieve images: %w", err)
    }

    err = json.Unmarshal(resp, &images)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshall to json: %w", err)
    }
    return images, nil
}

func (c *Client) GetImage(ctx context.Context, id int) (*models.Image, error) {
    var image models.Image

    resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/images/%d", id), nil)
    if err != nil {
        return nil, fmt.Errorf("could not obtain image: %w", err)
    }

    err = json.Unmarshal(resp, &image)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshall to json: %w", err)
    }

    return &image, nil
}

func (c *Client) SearchImages(ctx context.Context, species string, userID int) ([]models.Image, error) {
    var images []models.Image

    path := "/images"

    params := url.Values{}
    if species != "" {
        params.Add("species", species)
    }
    if userID > 0 {
        params.Add("user_id", strconv.Itoa(userID))
    }

    if len(params) > 0 {
        path = path + "?" + params.Encode()
    }

    resp, err := c.doRequest(ctx, "GET", path, nil)
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

**File**: `pkg/api/keys.go`

```go
func (client *Client) GenerateKey(ctx context.Context, name string) (*models.ApiKey, error) {
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

    resp, err := client.doRequest(ctx, "POST", "/api/keys", jsonData)
    if err != nil {
        return nil, fmt.Errorf("could not create api keys: %w", err)
    }

    err = json.Unmarshal(resp, &apiKey)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshal response: %w", err)
    }

    return &apiKey, nil
}

func (client *Client) ListKeys(ctx context.Context) ([]models.ApiKey, error) {
    var apiKeys []models.ApiKey

    resp, err := client.doRequest(ctx, "GET", "/api/keys", nil)
    if err != nil {
        return nil, fmt.Errorf("could not get apikeys: %w", err)
    }

    err = json.Unmarshal(resp, &apiKeys)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshall json: %w", err)
    }

    return apiKeys, nil
}

func (client *Client) RevokeKey(ctx context.Context, id int) error {
    _, err := client.doRequest(ctx, "DELETE", fmt.Sprintf("/api/keys/%d", id), nil)
    if err != nil {
        return fmt.Errorf("could not delete api-key: %w", err)
    }

    return nil
}
```

---

#### Step 3: Update Commands to Create Contexts

**File**: `cmd/images.go`

**Example: List Images Command**

```go
var listImagesCmd = &cobra.Command{
    Use:   "list",
    Short: "Get list of images",
    Args:  cobra.ExactArgs(0),
    Run: func(command *cobra.Command, args []string) {
        baseUrl, _ := config.Get("api-url")
        key, _ := config.Get("api-key")

        if !checkApiKey(key) {
            return
        }

        client := api.NewClient(baseUrl, key)

        // Create context with 30 second timeout
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        // Pass context to API method
        resp, err := client.ListImages(ctx)
        if err != nil {
            // Check if timeout occurred
            if ctx.Err() == context.DeadlineExceeded {
                fmt.Println("Error: Request timed out after 30 seconds")
                fmt.Println("The API might be slow or unavailable. Try again later.")
                return
            }
            fmt.Printf("could not retrieve images: %v\n", err)
            return
        }

        // Display results...
        w := tablewriter.NewWriter(os.Stdout)
        w.SetHeader([]string{"ID", "Species", "GPS Long", "GPS Lat", "Image Path"})

        for _, image := range resp {
            w.Append([]string{
                strconv.Itoa(image.ID),
                image.SpeciesName,
                fmt.Sprintf("%.6f", image.GpsLong),
                fmt.Sprintf("%.6f", image.GpsLat),
                image.ImagePath,
            })
        }

        w.Render()
    },
}
```

**Example: Get Image Command**

```go
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

        if !validatePositiveInt(integer) {
            return
        }

        baseUrl, _ := config.Get("api-url")
        key, _ := config.Get("api-key")

        if !checkApiKey(key) {
            return
        }

        client := api.NewClient(baseUrl, key)

        // Create context with 10 second timeout
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        image, err := client.GetImage(ctx, integer)
        if err != nil {
            if ctx.Err() == context.DeadlineExceeded {
                fmt.Println("Error: Request timed out after 10 seconds")
                return
            }
            fmt.Printf("could not retrieve image data: %v\n", err)
            return
        }

        // Display result...
        fmt.Printf("ID: %d\n", image.ID)
        fmt.Printf("Species: %s\n", image.SpeciesName)
        fmt.Printf("GPS: (%.6f, %.6f)\n", image.GpsLong, image.GpsLat)
        fmt.Printf("Path: %s\n", image.ImagePath)
    },
}
```

---

### Context Error Handling

#### Common Context Errors

```go
// Check what kind of error occurred
if err != nil {
    switch {
    case errors.Is(err, context.DeadlineExceeded):
        fmt.Println("Request timed out")
    case errors.Is(err, context.Canceled):
        fmt.Println("Request was cancelled")
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

#### Or use ctx.Err() directly:

```go
if ctx.Err() == context.DeadlineExceeded {
    fmt.Println("Timeout!")
} else if ctx.Err() == context.Canceled {
    fmt.Println("Cancelled!")
}
```

---

### Testing Context Behavior

Create a test file to experiment:

**File**: `experiments/context_test.go`

```go
package experiments

import (
    "context"
    "fmt"
    "testing"
    "time"
)

func TestContextTimeout(t *testing.T) {
    // Create context that times out in 2 seconds
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    // Simulate slow operation
    select {
    case <-time.After(5 * time.Second):
        t.Fatal("Should have timed out!")
    case <-ctx.Done():
        fmt.Println("Context timed out as expected:", ctx.Err())
        // Output: context deadline exceeded
    }
}

func TestContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())

    // Simulate work happening in goroutine
    go func() {
        time.Sleep(1 * time.Second)
        cancel() // Cancel after 1 second
    }()

    // Wait for cancellation
    <-ctx.Done()
    fmt.Println("Context was cancelled:", ctx.Err())
    // Output: context canceled
}

func TestContextValue(t *testing.T) {
    // Create context with value
    ctx := context.WithValue(context.Background(), "requestID", "abc-123")

    // Retrieve value
    requestID := ctx.Value("requestID").(string)
    fmt.Println("Request ID:", requestID)
    // Output: Request ID: abc-123
}
```

Run with: `go test -v ./experiments`

---

### Real-World Context Patterns

#### Pattern 1: Timeout with Fallback

```go
func getImageWithFallback(client *api.Client, id int) (*models.Image, error) {
    // Try with short timeout first
    ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel1()

    image, err := client.GetImage(ctx1, id)
    if err == nil {
        return image, nil
    }

    // If timeout, try slower backup server
    if ctx1.Err() == context.DeadlineExceeded {
        fmt.Println("Primary server slow, trying backup...")

        ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel2()

        backupClient := api.NewClient("https://backup-api.com", key)
        return backupClient.GetImage(ctx2, id)
    }

    return nil, err
}
```

#### Pattern 2: Parent-Child Contexts

```go
func processMultipleImages(imageIDs []int) error {
    // Parent context - 1 minute for entire operation
    parentCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
    defer cancel()

    for _, id := range imageIDs {
        // Child context - 5 seconds per image
        childCtx, childCancel := context.WithTimeout(parentCtx, 5*time.Second)

        image, err := client.GetImage(childCtx, id)
        childCancel() // Clean up immediately after use

        if err != nil {
            // If parent context expired, stop everything
            if parentCtx.Err() != nil {
                return fmt.Errorf("operation took too long: %w", parentCtx.Err())
            }
            // Otherwise just skip this image
            fmt.Printf("Skipping image %d: %v\n", id, err)
            continue
        }

        processImage(image)
    }

    return nil
}
```

---

### Summary: Context Best Practices

✅ **DO:**
- Pass context as first parameter
- Create context with timeout for external calls (API, database)
- Always `defer cancel()`
- Check `ctx.Err()` to distinguish timeout vs cancellation
- Use `context.Background()` as root for new request chains

❌ **DON'T:**
- Store context in structs
- Pass nil context (use `context.TODO()` if unsure)
- Ignore context in long-running operations
- Use `context.WithValue()` for everything (only for request-scoped data)

---

### Next Steps

Once you've implemented context:
1. Test timeout behavior with slow network
2. Try cancelling requests with Ctrl+C
3. Move on to **Pattern 2: Interfaces** (enables proper testing)

---

**Status**: 📝 Documented | ⏳ Not Implemented Yet
**Next Pattern**: Interfaces (for testability and mocking)
