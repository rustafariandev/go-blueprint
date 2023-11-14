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
		PackageNames: []string{"github.com/julienschmidt/httprouter"},
		ProjectType:  "httprouter",
	}
}

func init() {
	registry.RegisterFramework("httprouter", MakeProvider)
	steps.RegisterFrameworkItems(
		steps.Item{
			Title: "HttpRouter",
			Desc:  "HttpRouter is a lightweight high performance HTTP request router for Go",
		},
	)
}
