package chi

import (
	"embed"

	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/registry"
	"github.com/melkeydev/go-blueprint/cmd/steps"
)

//go:embed internal cmd all:*.tmpl
var FS embed.FS

func MakeProvider() *provider.TemplateProvider {
	return &provider.TemplateProvider{
		TempateFS:    FS,
		PackageNames: []string{},
		ProjectType:  "standard-library",
	}
}

func init() {
	registry.RegisterFramework("standard library", MakeProvider)
	steps.RegisterFrameworkItems(
		steps.Item{
			Title: "Standard library",
			Desc:  "The built-in Go standard library HTTP package",
		},
	)
}
