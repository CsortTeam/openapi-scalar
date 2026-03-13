package docs

import (
	"encoding/json"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v3"
	"gopkg.in/yaml.v3"
)

func Spec(routes []RouteInfo, opts Options) map[string]any {
	return buildSpec(routes, opts)
}

func Register(app *fiber.App, routes []RouteInfo, opts Options) {
	prefix := opts.PathPrefix
	if prefix == "" {
		prefix = "/api/docs"
	}

	spec := buildSpec(routes, opts)
	specJSON, _ := json.Marshal(spec)

	app.Get(prefix, func(c fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: string(specJSON),
			CustomOptions: scalar.CustomOptions{
				PageTitle: opts.Title,
			},
			DarkMode: opts.DarkMode,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "failed to render docs UI"})
		}
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		c.Set(fiber.HeaderCacheControl, "no-store")
		return c.SendString(htmlContent)
	})

	app.Get(prefix+"/openapi.json", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		c.Set(fiber.HeaderCacheControl, "no-store")
		return c.Send(specJSON)
	})

	app.Get(prefix+"/openapi.yaml", func(c fiber.Ctx) error {
		content, err := yaml.Marshal(spec)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "failed to generate OpenAPI spec"})
		}
		c.Set(fiber.HeaderContentType, "application/yaml; charset=utf-8")
		c.Set(fiber.HeaderCacheControl, "no-store")
		return c.Send(content)
	})
}
