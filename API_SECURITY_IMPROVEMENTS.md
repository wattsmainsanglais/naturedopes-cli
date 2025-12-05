# API Security Improvements Needed

**Project**: naturedopesApi
**Issue**: API key management endpoints are currently unprotected
**Priority**: HIGH (security vulnerability)
**Status**: TODO

---

## Current Security Issues

### 1. All API Key Endpoints Are Unprotected

From `naturedopesApi/main.go` lines 43-46:

```go
// API key management endpoints (unprotected for now)
router.HandleFunc("/api/keys", createApiKeyHandler).Methods("POST")
router.HandleFunc("/api/keys", getApiKeysHandler).Methods("GET")
router.HandleFunc("/api/keys/{id}", revokeApiKeyHandler).Methods("DELETE")
```

**Problems:**
- ❌ Anyone can list ALL API keys (privacy issue)
- ❌ Anyone can revoke ANY API key (security vulnerability)
- ❌ No ownership concept - can't tell which user owns which key

---

## Recommended Changes

### Phase 1: Add User Association (Database)

**1.1. Add `user_id` column to `api_keys` table**

```sql
-- Migration to add user ownership
ALTER TABLE api_keys ADD COLUMN user_id INTEGER;

-- For existing keys, you'll need to decide:
-- Option A: Assign all to admin user
UPDATE api_keys SET user_id = 1 WHERE user_id IS NULL;

-- Option B: Delete existing keys and start fresh
DELETE FROM api_keys;

-- Make column required going forward
ALTER TABLE api_keys ALTER COLUMN user_id SET NOT NULL;

-- Add foreign key if you have a users table
-- ALTER TABLE api_keys ADD CONSTRAINT fk_user
--   FOREIGN KEY (user_id) REFERENCES users(id);
```

**1.2. Update `GenerateApiKey()` function**

File: `endpoints/apikey.go` line 26

```go
// Before:
func GenerateApiKey(conn *pgx.Conn, name string, ipAddress string) (*ApiKey, error)

// After:
func GenerateApiKey(conn *pgx.Conn, name string, ipAddress string, userID int) (*ApiKey, error)
```

Update the INSERT query (line 43):
```go
// Before:
"INSERT INTO api_keys (key, name, created_at, expires_at, revoked, created_ip)
 VALUES ($1, $2, NOW(), $3, false, $4)
 RETURNING id, key, name, created_at, expires_at, revoked, created_ip",
key, name, expiresAt, ipAddress,

// After:
"INSERT INTO api_keys (key, name, created_at, expires_at, revoked, created_ip, user_id)
 VALUES ($1, $2, NOW(), $3, false, $4, $5)
 RETURNING id, key, name, created_at, expires_at, revoked, created_ip, user_id",
key, name, expiresAt, ipAddress, userID,
```

**1.3. Update `GetApiKeys()` function**

File: `endpoints/apikey.go` line 104

Add a parameter to filter by user:

```go
// Before:
func GetApiKeys(conn *pgx.Conn) ([]ApiKey, error)

// After:
func GetApiKeys(conn *pgx.Conn, userID *int) ([]ApiKey, error)
```

Update the query (line 110):
```go
// Before:
"SELECT id, key, name, created_at, expires_at, last_used, revoked, created_ip
 FROM api_keys
 ORDER BY created_at DESC"

// After - if userID provided, filter by it:
var query string
var args []interface{}

if userID != nil {
    query = `SELECT id, key, name, created_at, expires_at, last_used, revoked, created_ip, user_id
             FROM api_keys
             WHERE user_id = $1
             ORDER BY created_at DESC`
    args = append(args, *userID)
} else {
    // Admin view - show all keys (or remove this branch)
    query = `SELECT id, key, name, created_at, expires_at, last_used, revoked, created_ip, user_id
             FROM api_keys
             ORDER BY created_at DESC`
}

rows, err := conn.Query(ctx, query, args...)
```

**1.4. Update `RevokeApiKey()` function**

File: `endpoints/apikey.go` line 135

Add user ownership check:

```go
// Before:
func RevokeApiKey(conn *pgx.Conn, id int) error

// After:
func RevokeApiKey(conn *pgx.Conn, id int, userID int) error
```

Update the query (line 140):
```go
// Before:
"UPDATE api_keys SET revoked = true WHERE id = $1", id

// After - only allow revoking your own keys:
"UPDATE api_keys SET revoked = true WHERE id = $1 AND user_id = $2", id, userID
```

---

### Phase 2: Add Authentication Middleware

**2.1. Create helper to extract user ID from API key**

File: `middleware/auth.go` (create new file or add to existing)

```go
package middleware

import (
    "context"
    "net/http"
    "github.com/jackc/pgx/v4"
)

type contextKey string

const UserIDKey contextKey = "userID"

// ApiKeyAuthWithUser validates API key and adds user ID to context
func ApiKeyAuthWithUser(getDBConnection func() (*pgx.Conn, error)) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            apiKey := r.Header.Get("X-API-Key")

            if apiKey == "" {
                http.Error(w, "API key required", http.StatusUnauthorized)
                return
            }

            conn, err := getDBConnection()
            if err != nil {
                http.Error(w, "Database error", http.StatusInternalServerError)
                return
            }
            defer conn.Close(context.Background())

            // Get user ID from API key
            var userID int
            var revoked bool
            var expiresAt time.Time

            err = conn.QueryRow(
                context.Background(),
                "SELECT user_id, revoked, expires_at FROM api_keys WHERE key = $1",
                apiKey,
            ).Scan(&userID, &revoked, &expiresAt)

            if err != nil {
                http.Error(w, "Invalid API key", http.StatusUnauthorized)
                return
            }

            if revoked {
                http.Error(w, "API key has been revoked", http.StatusUnauthorized)
                return
            }

            if time.Now().After(expiresAt) {
                http.Error(w, "API key has expired", http.StatusUnauthorized)
                return
            }

            // Add user ID to request context
            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(r *http.Request) (int, bool) {
    userID, ok := r.Context().Value(UserIDKey).(int)
    return userID, ok
}
```

**2.2. Update route handlers**

File: `routes.go`

```go
func createApiKeyHandler(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from context
    userID, ok := middleware.GetUserIDFromContext(r)
    if !ok {
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }

    var req struct {
        Name string `json:"name"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Name == "" {
        http.Error(w, "Name is required", http.StatusBadRequest)
        return
    }

    ip := r.Header.Get("X-Forwarded-For")
    if ip == "" {
        ip = r.Header.Get("X-Real-IP")
    }
    if ip == "" {
        ip = r.RemoteAddr
    }

    conn, err := connectToDB()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    // Pass userID to GenerateApiKey
    apiKey, err := endpoints.GenerateApiKey(conn, req.Name, ip, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(apiKey)
}

func getApiKeysHandler(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from context
    userID, ok := middleware.GetUserIDFromContext(r)
    if !ok {
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }

    conn, err := connectToDB()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    // Only get keys for this user
    apiKeys, err := endpoints.GetApiKeys(conn, &userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(apiKeys)
}

func revokeApiKeyHandler(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from context
    userID, ok := middleware.GetUserIDFromContext(r)
    if !ok {
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }

    params := mux.Vars(r)
    id := params["id"]
    idInt, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    conn, err := connectToDB()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    // Pass userID to RevokeApiKey (only revokes if owned by user)
    err = endpoints.RevokeApiKey(conn, idInt, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
```

**2.3. Update routes to use authentication**

File: `main.go`

```go
func SetupRoutes(router *mux.Router) {
    // Create auth middleware that adds user ID to context
    apiKeyAuth := middleware.ApiKeyMiddleware(connectToDB)
    apiKeyAuthWithUser := middleware.ApiKeyAuthWithUser(connectToDB)

    // Protected image endpoints (require API key)
    imageRouter := router.PathPrefix("/images").Subrouter()
    imageRouter.Use(apiKeyAuth)
    imageRouter.HandleFunc("", getImagesHandler).Methods("GET")
    imageRouter.HandleFunc("/{id}", getImageHandler).Methods("GET")

    // API key management endpoints
    // Option A: Protect all endpoints (users manage their own keys)
    keyRouter := router.PathPrefix("/api/keys").Subrouter()
    keyRouter.Use(apiKeyAuthWithUser) // Requires auth, adds user ID to context
    keyRouter.HandleFunc("", createApiKeyHandler).Methods("POST")
    keyRouter.HandleFunc("", getApiKeysHandler).Methods("GET")
    keyRouter.HandleFunc("/{id}", revokeApiKeyHandler).Methods("DELETE")

    // Option B: Allow unauthenticated key creation (bootstrap problem)
    // Keep POST /api/keys unprotected
    router.HandleFunc("/api/keys", createApiKeyHandler).Methods("POST")

    keyRouter := router.PathPrefix("/api/keys").Subrouter()
    keyRouter.Use(apiKeyAuthWithUser)
    keyRouter.HandleFunc("", getApiKeysHandler).Methods("GET")
    keyRouter.HandleFunc("/{id}", revokeApiKeyHandler).Methods("DELETE")
}
```

---

## The Bootstrap Problem

### Problem:
If all key endpoints require authentication, how does a user get their first key?

### Solutions:

**Option 1: Public Key Creation (Current approach - keep it)**
- `POST /api/keys` remains unprotected
- Anyone can create a key
- `GET` and `DELETE` require authentication
- Pro: Easy onboarding
- Con: Potential abuse (rate limiting helps)

**Option 2: Registration Endpoint**
- Create separate `/auth/register` endpoint
- Returns initial API key
- All `/api/keys` endpoints protected
- Pro: More controlled
- Con: More complex

**Option 3: Admin-Generated Keys**
- Keys created via admin panel or CLI tool
- No public key creation
- All endpoints protected
- Pro: Most secure
- Con: Requires manual provisioning

**Recommendation**: Option 1 (what you have) with good rate limiting

---

## Summary of Changes Needed

### Database:
- [ ] Add `user_id` column to `api_keys` table
- [ ] Add foreign key constraint (if users table exists)
- [ ] Migrate existing data

### Code Changes:
- [ ] `endpoints/apikey.go`:
  - [ ] Update `GenerateApiKey()` to accept `userID`
  - [ ] Update `GetApiKeys()` to filter by `userID`
  - [ ] Update `RevokeApiKey()` to check ownership
- [ ] `middleware/auth.go`:
  - [ ] Create `ApiKeyAuthWithUser()` middleware
  - [ ] Create `GetUserIDFromContext()` helper
- [ ] `routes.go`:
  - [ ] Update all handlers to extract user ID from context
  - [ ] Pass user ID to endpoint functions
- [ ] `main.go`:
  - [ ] Apply `apiKeyAuthWithUser` middleware to protected routes
  - [ ] Decide on bootstrap approach

### Testing:
- [ ] Test unauthenticated key creation works
- [ ] Test users can only see their own keys
- [ ] Test users can only revoke their own keys
- [ ] Test expired/revoked keys are rejected
- [ ] Test rate limiting still works

---

## Impact on CLI

The naturedopes-cli will continue to work because:
- ✅ It already sends `X-API-Key` header (via `api.NewClient()`)
- ✅ `generate` command works without auth (if bootstrap Option 1)
- ✅ `list` command will show only user's keys (better privacy)
- ✅ `revoke` command will only revoke user's keys (better security)

**No CLI changes needed!** The security improvements are backward-compatible.

---

## Timeline

**Phase 1** (Database + Basic Auth): 1-2 hours
**Phase 2** (Middleware + Route Protection): 1-2 hours
**Testing**: 1 hour

**Total**: 3-5 hours of focused work

---

**Created**: 2025-12-05
**Author**: Claude (via naturedopes-cli project)
**Related**: Phase 6 - API Keys Commands
