package docs

type RouteInfo struct {
	Method string
	Path   string
	Doc    *DocInfo
}

type DocInfo struct {
	Summary     string
	Description string
	Tags        []string
	Parameters  []ParamInfo
	RequestBody *RequestBodyInfo
	Responses   map[string]ResponseInfo
	Security    []map[string][]string
}

type ParamInfo struct {
	Name        string
	In          string
	Description string
	Required    bool
	Schema      map[string]any
}

type RequestBodyInfo struct {
	Required    bool
	ContentType string
	Schema      map[string]any
}

type ResponseInfo struct {
	Description string
	ContentType string
	Schema      map[string]any
}

type Options struct {
	Title       string
	Version     string
	PathPrefix  string
	APIPrefix   string
	DarkMode    bool
}
