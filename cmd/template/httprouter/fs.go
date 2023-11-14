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
		PackageNames: []string{"github.com/julienschmidt/httprouter"},
		ProjectType:  "httprouter",
	}
}

func init() {
	template.RegisterProvider("httprouter", MakeProvider)
	steps.RegisterFrameworkItems(
		steps.Item{
			Title: "HttpRouter",
			Desc:  "HttpRouter is a lightweight high performance HTTP request router for Go",
		},
	)
}
