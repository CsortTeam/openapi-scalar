# openapi-scalar

Shared Scalar API docs for Fiber apps. Builds OpenAPI 3.0 spec from route metadata and serves interactive docs.

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

Endpoints:

- `GET /api/docs` — Scalar UI
- `GET /api/docs/openapi.json` — OpenAPI spec (JSON)
- `GET /api/docs/openapi.yaml` — OpenAPI spec (YAML)
