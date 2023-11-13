package chi

import (
	"embed"

	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/steps"
	"github.com/melkeydev/go-blueprint/cmd/template"
)

//go:embed internal cmd all:*.tmpl
var FS embed.FS

func MakeProvider() *provider.TemplateProvider {
	return &provider.TemplateProvider{
		TempateFS:    FS,
		PackageNames: []string{"github.com/gofiber/fiber/v2"},
		ProjectType:  "fiber",
	}
}

func init() {
	template.RegisterProvider("fiber", MakeProvider)
	steps.RegisterFrameworkItems(
		steps.Item{
			Title: "Fiber",
			Desc:  "An Express inspired web framework built on top of Fasthttp",
		},
	)
}
