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
		TempateFS: FS,
		PackageNames: []string{
			"github.com/go-chi/chi/v5",
		},
		ProjectType: "chi",
	}
}

func init() {
	template.RegisterProvider("chi", MakeProvider)
	steps.RegisterFrameworkItems(
		steps.Item{
			Title: "Chi",
			Desc:  "A lightweight, idiomatic and composable router for building Go HTTP services",
		},
	)
}
