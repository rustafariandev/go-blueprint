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
			"github.com/gin-gonic/gin",
		},
		ProjectType: "gin",
	}
}

func init() {
	template.RegisterProvider("gin", MakeProvider)
	steps.RegisterFrameworkItems(
		steps.Item{
			Title: "Gin",
			Desc:  "Features a martini-like API with performance that is up to 40 times faster thanks to httprouter",
		},
	)
}
