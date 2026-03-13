# openapi-scalar

Scalar API docs for Fiber. Builds OpenAPI 3.0 from route metadata.

## Usage

```go
import (
    docs "csort.ru/openapi-scalar"
    "github.com/gofiber/fiber/v3"
)

routes := []docs.RouteInfo{
    {Method: "GET", Path: "/api/v1/health"},
    {Method: "GET", Path: "/api/v1/users/:id", Doc: &docs.DocInfo{
        Summary: "Get user by ID",
        Tags:    []string{"users"},
    }},
}

docs.Register(app, routes, docs.Options{
    Title:      "My Service API",
    Version:    "1.0.0",
    PathPrefix: "/api/docs",
    APIPrefix:  "/api/v1",
    DarkMode:   true,
})
```

Endpoints: `GET /api/docs` (Scalar UI), `GET /api/docs/openapi.json`, `GET /api/docs/openapi.yaml`

## Response format

```go
{Doc: &docs.DocInfo{
    Summary: "Get user",
    Tags:    []string{"users"},
    Responses: map[string]docs.ResponseInfo{
        "200": {Description: "Success", Schema: docs.SchemaFromType(User{})},
        "404": {Description: "Not found", Schema: map[string]any{"type": "object"}},
    },
}}
```

ContentType defaults to application/json.

```go
"200": {Description: "OK", ContentType: "text/plain", Schema: map[string]any{"type": "string"}}
```

## Request body

```go
{Doc: &docs.DocInfo{
    Summary: "Create user",
    RequestBody: &docs.RequestBodyInfo{
        Required: true,
        Schema:   docs.SchemaFromType(CreateUserRequest{}),
    },
}}
```

## Parameters

Path params from `:id` are auto-included. Use Parameters for query, header, etc:

```go
{Doc: &docs.DocInfo{
    Parameters: []docs.ParamInfo{
        {Name: "id", In: "path", Description: "User ID", Required: true},
        {Name: "X-Custom-Header", In: "header", Description: "Optional header", Required: false},
        {Name: "sort", In: "query", Description: "Sort field", Schema: map[string]any{"type": "string", "enum": []any{"name", "created"}}},
    },
}}
```

## Security (Bearer)

```go
{Doc: &docs.DocInfo{
    Summary:   "Get profile",
    Security:  []map[string][]string{{"bearerAuth": {}}},
}}
```

## Schema helpers

- `docs.SchemaFromType(v)` — JSON Schema from Go struct
- `docs.SchemaArray(item)` — array of item type
