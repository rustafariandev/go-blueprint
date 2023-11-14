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
		TempateFS: FS,
		PackageNames: []string{
			"github.com/gorilla/mux",
		},
		ProjectType: "gorilla/mux",
	}
}

func init() {
	registry.RegisterFramework("gorilla/mux", MakeProvider)
	steps.RegisterFrameworkItems(
		steps.Item{
			Title: "Gorilla/Mux",
			Desc:  "Package gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler",
			Value: "gorilla/mux",
		},
	)
}
