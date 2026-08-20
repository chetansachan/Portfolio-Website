# How to Write Test Cases for a Route

This guide explains how to write unit tests for a route in this Go project, based on
the pattern already used across every test file in this folder.

---

## 1. The Big Picture: How a Route Is Wired

Every route in this project is a thin 3-layer pipeline. Tests target the **controller**
layer in isolation by injecting a **mock service**.

```
HTTP request
   ↓
Router  (internals/router/routes.go)        ← matches URL + method, picks handler
   ↓
Controller  (internals/controller/*.go)      ← validates input, maps errors → status codes
   ↓
Service interface  (BlogServiceInterface)   ← business logic (we mock THIS in tests)
   ↓
Repository / DB                            ← not exercised in unit tests
```

Key insight: every controller depends on an **interface** (e.g. `BlogServiceInterface`),
not a concrete service. That means tests can build a `MockBlogService` that satisfies
the interface, hand it to the controller via `NewBlogController(mock)`, and exercise
**only the controller's behavior** — no DB, no network, no globals.

**You do not need to touch the router to test a route.** The controller method is
called directly with a synthetic `*http.Request`.

---

## 2. The Pattern: Table-Driven Tests

Every test file in `test/` uses the same shape:

```go
func TestXxx_Controller(t *testing.T) {
    // 1. (optional) shared sample fixtures
    sample := model.Blog{ ... }

    // 2. the table — one entry per scenario
    tests := []struct {
        name           string
        method         string
        body           string
        pathID         string
        mockFunc       func(ctx, ...) (...)
        expectedStatus int
        expectedMessage string
        expectData     bool
    }{
        { name: "Success — ...",         ... },
        { name: "Bad Request — ...",     ... },
        { name: "Method Not Allowed —",  ... },
        // ... one entry per scenario
    }

    // 3. the runner — same setup/asserts for every row
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := &MockXxxService{ XxxFunc: tt.mockFunc }
            ctrl := controller.NewXxxController(mock)

            req := httptest.NewRequest(tt.method, "/route"+tt.pathID, bodyReader)
            req.SetPathValue("id", tt.pathID)        // if route has {id}
            w := httptest.NewRecorder()
            ctrl.SomeHandler(w, req)

            // assert status, content-type, response body
        })
    }
}
```

Three blocks. That's the entire pattern. Everything below is "what to put in each block."

---

## 3. What Files Already Exist (and Why)

| File | Purpose |
| --- | --- |
| `test/mocks.go` | Mock implementations of every service interface. **One file, shared by every test.** You rarely need to edit it. |
| `test/<verb>_<entity>_test.go` | One file per controller method. |
| `internals/controller/*_controller.go` | The code under test. Read this to know what to assert. |
| `internals/router/routes.go` | Lists every route. Use it to confirm method + path. |

The mocks in `mocks.go` are **function-field mocks** — each mock method checks if its
`*Func` field is set and calls it, otherwise returns a zero value. This lets each
test case inject its own behavior with one line: `MyFunc: func(...) { ... }`.

---

## 4. The Case Categories — What to Cover for ANY Route

Every route should have cases in **all** of these buckets. The table below maps each
bucket to a real example in the existing tests so you can copy the pattern.

### 4.1 Happy Path (at least 1 case)
The "everything correct" scenario. Method matches, input is valid, mock returns success.
- **Files to look at:** `get_all_blogs_test.go` "Success — default pagination",
  `post_blog_test.go` "Success — all fields provided".

### 4.2 Wrong HTTP Method (≥ 1 per disallowed method)
Send a method the route doesn't allow and expect `405 Method Not Allowed`.
- Try **every** method the route does *not* accept. For a GET route, test POST, PUT,
  PATCH, DELETE.
- **Files to look at:** `delete_blog_test.go` "Method Not Allowed — GET/POST/PUT/PATCH
  instead of DELETE" (covers all four in one block).

### 4.3 Input Validation — Path Parameters
For routes with `{id}` or `{slug}` in the path, cover:
- Valid value
- Wrong type / wrong format (e.g. `not-a-uuid`, `12345`, `"   "`, `1 OR 1=1`)
- Empty string
- For UUIDs specifically: truncated UUIDs
- **Files to look at:** `get_blog_by_id_test.go` and `delete_blog_test.go` (the
  delete test is the most thorough — it covers SQL injection as a path value).

### 4.4 Input Validation — Query Parameters (if any)
For `?page=`, `?size=` etc., cover:
- Valid value (custom and default)
- Zero
- Negative number
- Non-numeric (`abc`, `xyz`)
- **Boundary cases the controller enforces** (e.g. `size > 100` gets capped at 100)
- **Files to look at:** `get_all_blogs_test.go` "Invalid query parameters" + "size
  capped at maximum of 100".

### 4.5 Input Validation — Request Body (POST/PATCH only)
- Completely invalid JSON (`this is not json`)
- Empty body
- Truncated JSON (`{"blog_title": "Hello`)
- JSON array instead of object (`[{...}]`)
- Empty `{}`
- Missing individual required fields
- All required fields empty strings
- **Files to look at:** `post_blog_test.go` "Malformed / Invalid JSON body" +
  "Validation failures".

### 4.6 Domain Errors Returned by the Service
Map each `errors.Is(err, model.Err...)` branch in the controller to one test case.
Open the controller file, find every `case errors.Is(...)` switch arm, and write
one test per arm.
- **Files to look at:** `get_blog_by_id_test.go` covers `ErrIdValidation`,
  `ErrBlogNotFound`. `post_blog_test.go` covers `ErrValidation` and
  `ErrBlogAlreadyExist`.

### 4.7 Internal Server Error (at least 1)
Make the mock return a plain `errors.New("anything")` so the controller hits its
`default` switch arm.
- **Files to look at:** every test file ends with one or two "Internal Server Error"
  cases. Use a different error message in each so it's obvious which arm ran.

### 4.8 Edge Cases (good to have, not required)
- Unknown / extra fields in JSON (should be ignored)
- Very long strings
- Special characters / Unicode / emoji
- Idempotency (e.g. delete the same thing twice)
- **Files to look at:** `post_blog_test.go` "Edge Case" block. `delete_blog_test.go`
  "already deleted (idempotency check)".

---

## 5. Step-by-Step Recipe: Adding a Test for a New Route

Let's say you just added a new route `PATCH /blog/publish/{id}` and want to test it.

**Step 1. Confirm the route exists.** Open `internals/router/routes.go` and verify:
```go
mux.HandleFunc("PATCH /blog/publish/{id}", blogController.PublishBlog_Controller)
```
Note the method (`PATCH`), path pattern (`/blog/publish/{id}`), and the handler
(`PublishBlog_Controller`).

**Step 2. Read the handler.** Open `internals/controller/blogs_controller.go` and
find `PublishBlog_Controller`. Make a list of every branch:
- method mismatch → 405
- `uuid.Parse` fails → 400 "Invalid ID format"
- `service.PublishBlog_Service` returns `ErrIdValidation` → 400
- `service.PublishBlog_Service` returns `ErrBlogNotFound` → 404
- `service.PublishBlog_Service` returns other error → 500
- success → 200

This list **is** your test table.

**Step 3. Check the mock.** Open `test/mocks.go`. Does `MockBlogService` have a
`PublishBlogFunc` field and a `PublishBlog_Service` method? If not, add them now
(usually just copy the pattern from a sibling method). Update
`BlogServiceInterface` in the controller if you haven't already.

**Step 4. Create the file.** Name it after the handler:
`test/publish_blog_test.go`.

**Step 5. Write the test file.** Copy the skeleton from `test/update_blog_test.go`
(PATCH is the closest cousin). Then fill the table with one row per branch from
Step 2, plus a couple of edge cases.

**Step 6. Build & run.**
```bash
go test ./test/... -run TestPublishBlog -v
```

---

## 6. Skeleton Template

Drop this into a new `test/<thing>_test.go` and fill in the marked spots.

```go
package test

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "portfolio-website/internals/controller"
    "portfolio-website/internals/model"

    "github.com/google/uuid"
)

// =============================================================================
// Tests for <METHOD> <PATH>            <-- e.g. PATCH /blog/publish/{id}
// =============================================================================

func Test<HandlerName>_Controller(t *testing.T) {
    // --- shared fixtures (optional) ---
    testUUID := uuid.New()
    sample := model.Blog{
        Blog_ID:    testUUID,
        Blog_Title: "Sample",
        // ... fill the rest
        Created_At: time.Now(),
    }

    tests := []struct {
        name string

        // Request configuration — only the fields your route actually uses
        method string
        pathID string                  // for /{id}
        pathSlug string                // for /{slug}
        body   string                  // for POST/PATCH
        queryParams string             // for ?page=&size=

        // Mock configuration — one func per service method this handler calls
        handlerFunc func(ctx context.Context, /* match the interface */) (/* return type */, error)

        // Expected response
        expectedStatus  int
        expectedSuccess bool
        expectedMessage string
        expectData      bool         // true if success path returns a body
    }{
        // ===== HAPPY PATH =====
        {
            name:   "Success — <describe the scenario>",
            method: http.MethodGet,
            pathID: testUUID.String(),
            handlerFunc: func(ctx context.Context, /* ... */) (/* ... */, error) {
                return sample, nil
            },
            expectedStatus:  http.StatusOK, // or StatusFound, StatusCreated, etc.
            expectedSuccess: true,
            expectedMessage: "<exact string the controller writes>",
            expectData:      true,
        },

        // ===== METHOD NOT ALLOWED (one per wrong method) =====
        {
            name:            "Method Not Allowed — POST instead of GET",
            method:          http.MethodPost,
            pathID:          testUUID.String(),
            handlerFunc:     nil,
            expectedStatus:  http.StatusMethodNotAllowed,
            expectedSuccess: false,
            expectedMessage: "Method not allowed",
        },

        // ===== INPUT VALIDATION (path / query / body) =====
        // ... add rows for every invalid input the controller rejects

        // ===== DOMAIN ERRORS (one per errors.Is branch) =====
        // ... add a row per ErrXxx the service can return

        // ===== INTERNAL SERVER ERROR =====
        {
            name:   "Internal Server Error — <what went wrong>",
            method: http.MethodGet,
            pathID: testUUID.String(),
            handlerFunc: func(ctx context.Context, /* ... */) (/* ... */, error) {
                return model.Blog{}, errors.New("anything")
            },
            expectedStatus:  http.StatusInternalServerError,
            expectedSuccess: false,
            expectedMessage: "Internal server error. Please try again later.",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockService := &MockBlogService{
                HandlerFunc: tt.handlerFunc, // wire YOUR service method here
            }
            ctrl := controller.NewBlogController(mockService)

            // Build the request — only set what your route needs
            var body *strings.Reader
            if tt.body != "" {
                body = strings.NewReader(tt.body)
            } else {
                body = strings.NewReader("")
            }
            req := httptest.NewRequest(tt.method, "/route/path"+tt.queryParams, body)
            req.Header.Set("Content-Type", "application/json")
            if tt.pathID != "" {
                req.SetPathValue("id", tt.pathID)
            }
            if tt.pathSlug != "" {
                req.SetPathValue("slug", tt.pathSlug)
            }

            w := httptest.NewRecorder()
            ctrl.<YourHandler>(w, req)

            // --- Assert status code ---
            if w.Code != tt.expectedStatus {
                t.Errorf("Status code: got %d, want %d", w.Code, tt.expectedStatus)
            }

            // --- Assert Content-Type ---
            if got := w.Header().Get("Content-Type"); got != "application/json" {
                t.Errorf("Content-Type: got %q, want %q", got, "application/json")
            }

            // --- Assert body ---
            if tt.expectData {
                // Success path: choose the response struct that matches your handler
                var resp model.<ResponseStructForSuccess>
                if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
                    t.Fatalf("Failed to decode response: %v", err)
                }
                if resp.Success != tt.expectedSuccess {
                    t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
                }
                if resp.Message != tt.expectedMessage {
                    t.Errorf("Message: got %q, want %q", resp.Message, tt.expectedMessage)
                }
                // assert Data fields specific to this handler
            } else {
                // Error path
                var resp model.Response
                if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
                    t.Fatalf("Failed to decode error response: %v", err)
                }
                if resp.Success != tt.expectedSuccess {
                    t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
                }
                if resp.Message != tt.expectedMessage {
                    t.Errorf("Message: got %q, want %q", resp.Message, tt.expectedMessage)
                }
            }
        })
    }
}
```

---

## 7. Conventions Used Throughout the Test Suite

These are not rules — they're what's already in every test file, so follow them so
the next person reading your test feels at home.

| Convention | Why |
| --- | --- |
| File name: `<verb>_<entity>_test.go` | e.g. `post_blog_test.go`, `delete_project_test.go` |
| One test func per controller handler | `Test<HandlerName>_Controller` |
| Test names start with the bucket, then `—` and a one-line description | e.g. `"Bad Request — empty ID string"`. They read like headlines in test output. |
| Use `http.MethodGet`, `http.StatusBadRequest` etc. | Not raw `"GET"`, `400` strings. |
| `expectedMessage` must match the controller's exact string | Don't paraphrase. If you change the controller's message, update the test. |
| Set `Content-Type: application/json` in every request that has a body | Mirrors real clients. |
| Use `req.SetPathValue("id", ...)` for path params | This is what `http.ServeMux` does in production; using it keeps the test honest. |
| Mock returns the **exact** `model.Err*` sentinel the controller checks | Use `errors.New("...")` for the catch-all 500 case only. |
| For paginated routes, also assert `Meta.CurrentPage`, `Meta.PageSize`, `Meta.TotalBlogs` | See `get_all_blogs_test.go`. |
| Group cases with `// ===== BUCKET NAME =====` comments | Makes the table scannable. |

---

## 8. Response Struct Cheat Sheet

Different handlers return different shapes. Match the struct to the handler:

| Handler returns | Decode as |
| --- | --- |
| Single blog + meta | `model.PostBlogResponse` |
| Single project + meta | `model.PostProjectResponse` |
| List of blogs | `model.GetBlogsResponse` |
| Paginated blogs | `model.PaginatedResponse` |
| Paginated projects | `model.PaginatedProjectResponse` |
| Plain success/error | `model.Response` |
| Project plain success/error | `model.ProjectResponse` |

When in doubt, open the controller, find the `json.NewEncoder(w).Encode(...)` line,
and use the struct it encodes.

---

## 9. Common Pitfalls

1. **Don't hit a real database.** The whole point of the mock layer is to keep tests
   hermetic. If your test needs a DB, you want an integration test, not a unit test,
   and it doesn't belong in this folder.

2. **Don't test the service or the repository here.** Those have their own concerns.
   If you want to test them, write a separate test file that exercises the concrete
   service against a real or test DB.

3. **Don't forget the `Content-Type` assertion.** Several routes set it in the
   handler AND in the test. Keep both for symmetry with the real handler flow.

4. **Wrong method cases should leave the mock function as `nil`.** If you set a
   mock function on a 405 test, you've made the mock reachable when it shouldn't be.

5. **`expectedMessage` is brittle by design.** It catches accidental message changes
   in the controller. If you intentionally change a message, update the test
   *and* the doc/README that references it.

6. **`uuid.Nil.String()` is `"00000000-0000-0000-0000-000000000000"`.** Use it to
   trigger the `ErrIdValidation` branch in the service.

---

## 10. Quick Checklist Before You Commit

- [ ] Every test name is a single readable sentence
- [ ] At least one happy-path case
- [ ] Method-not-allowed cases cover all disallowed methods
- [ ] Every `errors.Is` arm in the controller has a corresponding case
- [ ] At least one Internal Server Error case
- [ ] Body validation cases (if route accepts a body)
- [ ] Query param validation cases (if route accepts query params)
- [ ] Path param validation cases (if route has `{id}` or `{slug}`)
- [ ] `expectedMessage` matches the controller string **exactly**
- [ ] `go test ./test/...` is green
- [ ] `go test ./test/... -v` reads like a checklist of behavior

That's the whole playbook. Pick a handler, follow the recipe, ship the table.
