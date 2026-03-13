package docs

// RouteInfo describes a route for OpenAPI spec generation.
type RouteInfo struct {
	Method string
	Path   string
	Doc    *DocInfo
}

// DocInfo holds OpenAPI metadata for a route.
type DocInfo struct {
	Summary     string
	Description string
	Tags        []string
	Parameters  []ParamInfo
	RequestBody *RequestBodyInfo
	Responses   map[string]ResponseInfo
	Security    []map[string][]string
}

// ParamInfo describes a request parameter.
type ParamInfo struct {
	Name        string
	In          string
	Description string
	Required    bool
	Schema      map[string]any
}

// RequestBodyInfo describes a request body.
type RequestBodyInfo struct {
	Required    bool
	ContentType string
	Schema      map[string]any
}

// ResponseInfo describes a response.
type ResponseInfo struct {
	Description string
	ContentType string
	Schema      map[string]any
}

// Options configures the docs handler.
type Options struct {
	Title       string
	Version     string
	PathPrefix  string
	APIPrefix   string
	DarkMode    bool
}
