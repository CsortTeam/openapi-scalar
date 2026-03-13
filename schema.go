package docs

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

// SchemaFromType reflects a Go struct into a JSON Schema map.
// Useful for building request/response schemas from domain types.
func SchemaFromType(v any) map[string]any {
	reflector := jsonschema.Reflector{DoNotReference: true}
	schema := reflector.Reflect(v)
	content, err := json.Marshal(schema)
	if err != nil {
		return map[string]any{"type": "object"}
	}
	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		return map[string]any{"type": "object"}
	}
	delete(result, "$schema")
	delete(result, "$id")
	return result
}

// SchemaArray reflects a Go struct as an array-of-items schema.
func SchemaArray(item any) map[string]any {
	return map[string]any{
		"type":  "array",
		"items": SchemaFromType(item),
	}
}
