package docs

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

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

func SchemaArray(item any) map[string]any {
	return map[string]any{
		"type":  "array",
		"items": SchemaFromType(item),
	}
}
