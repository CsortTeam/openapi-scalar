package docs

import (
	"sort"
	"strings"
)

func buildSpec(routes []RouteInfo, opts Options) map[string]any {
	paths := map[string]any{}
	tagsSet := map[string]struct{}{}

	apiPrefix := opts.APIPrefix
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}

	for _, route := range routes {
		if !strings.HasPrefix(route.Path, apiPrefix) {
			continue
		}

		method := strings.ToLower(route.Method)
		if method == "" || method == "use" {
			continue
		}

		openAPIPath := toOpenAPIPath(route.Path)

		pathItem, ok := paths[openAPIPath].(map[string]any)
		if !ok {
			pathItem = map[string]any{}
			paths[openAPIPath] = pathItem
		}

		doc := route.Doc
		tag := inferTag(openAPIPath, apiPrefix)
		if doc != nil && len(doc.Tags) > 0 && doc.Tags[0] != "" {
			tag = doc.Tags[0]
		}
		tagsSet[tag] = struct{}{}

		operation := map[string]any{"tags": []string{tag}}
		if doc != nil {
			if doc.Summary != "" {
				operation["summary"] = doc.Summary
			}
			if doc.Description != "" {
				operation["description"] = doc.Description
			}
			if len(doc.Tags) > 0 {
				operation["tags"] = doc.Tags
			}
			if len(doc.Security) > 0 {
				operation["security"] = doc.Security
			}
		}

		params := buildOperationParameters(route.Path, doc)
		if len(params) > 0 {
			operation["parameters"] = params
		}

		operation["responses"] = buildOperationResponses(method, doc)
		if doc != nil && doc.RequestBody != nil {
			operation["requestBody"] = buildRequestBody(doc.RequestBody)
		}

		pathItem[method] = operation
	}

	tagNames := make([]string, 0, len(tagsSet))
	for tag := range tagsSet {
		tagNames = append(tagNames, tag)
	}
	sort.Strings(tagNames)

	tags := make([]map[string]any, 0, len(tagNames))
	for _, tag := range tagNames {
		tags = append(tags, map[string]any{"name": tag})
	}

	title := opts.Title
	if title == "" {
		title = "API"
	}
	version := opts.Version
	if version == "" {
		version = "1.0.0"
	}

	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":   title,
			"version": version,
		},
		"servers": []map[string]any{
			{"url": "/", "description": "Default"},
		},
		"paths": paths,
		"tags":  tags,
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
		},
	}
}

func buildOperationParameters(path string, doc *DocInfo) []map[string]any {
	result := make([]map[string]any, 0)
	paramSet := map[string]struct{}{}

	for _, name := range extractPathParams(path) {
		key := "path:" + name
		paramSet[key] = struct{}{}
		result = append(result, map[string]any{
			"name":     name,
			"in":       "path",
			"required": true,
			"schema":   map[string]any{"type": "string"},
		})
	}

	if doc != nil {
		for _, p := range doc.Parameters {
			name := strings.TrimSpace(p.Name)
			in := strings.TrimSpace(p.In)
			if name == "" || in == "" {
				continue
			}
			key := in + ":" + name
			if _, exists := paramSet[key]; exists {
				continue
			}
			paramSet[key] = struct{}{}

			required := p.Required
			if in == "path" {
				required = true
			}
			schema := p.Schema
			if schema == nil {
				schema = map[string]any{"type": "string"}
			}

			item := map[string]any{
				"name":     name,
				"in":       in,
				"required": required,
				"schema":   schema,
			}
			if p.Description != "" {
				item["description"] = p.Description
			}
			result = append(result, item)
		}
	}

	return result
}

func buildOperationResponses(method string, doc *DocInfo) map[string]any {
	if doc != nil && len(doc.Responses) > 0 {
		result := map[string]any{}
		for status, r := range doc.Responses {
			response := map[string]any{"description": r.Description}
			if r.Schema != nil {
				ct := r.ContentType
				if ct == "" {
					ct = "application/json"
				}
				response["content"] = map[string]any{
					ct: map[string]any{"schema": r.Schema},
				}
			}
			result[status] = response
		}
		return result
	}

	status := defaultStatusByMethod(method)
	description := defaultDescriptionByMethod(method)
	response := map[string]any{"description": description}
	if status != "204" {
		response["content"] = map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{"type": "object"},
			},
		}
	}
	return map[string]any{status: response}
}

func buildRequestBody(body *RequestBodyInfo) map[string]any {
	ct := body.ContentType
	if ct == "" {
		ct = "application/json"
	}
	return map[string]any{
		"required": body.Required,
		"content": map[string]any{
			ct: map[string]any{"schema": body.Schema},
		},
	}
}

func toOpenAPIPath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	parts := strings.Split(path, "/")
	for i := range parts {
		if strings.HasPrefix(parts[i], ":") {
			parts[i] = "{" + strings.TrimPrefix(parts[i], ":") + "}"
		}
	}
	return normalizePath(strings.Join(parts, "/"))
}

func extractPathParams(path string) []string {
	if path == "" {
		return nil
	}
	var params []string
	for _, part := range strings.Split(path, "/") {
		if strings.HasPrefix(part, ":") {
			params = append(params, strings.TrimPrefix(part, ":"))
		}
	}
	return params
}

func inferTag(path, apiPrefix string) string {
	trimmed := strings.TrimPrefix(path, apiPrefix+"/")
	if trimmed == "" {
		return "api"
	}
	segment := strings.Split(trimmed, "/")[0]
	if segment == "" || strings.HasPrefix(segment, "{") {
		return "api"
	}
	return segment
}

func normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func defaultStatusByMethod(method string) string {
	switch method {
	case "post":
		return "201"
	case "delete":
		return "204"
	default:
		return "200"
	}
}

func defaultDescriptionByMethod(method string) string {
	switch method {
	case "post":
		return "Created"
	case "delete":
		return "No Content"
	default:
		return "Success"
	}
}

