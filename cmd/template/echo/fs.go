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
			"github.com/labstack/echo/v4",
			"github.com/labstack/echo/v4/middleware",
		},
		ProjectType: "echo",
	}
}

func init() {
	template.RegisterProvider("echo", MakeProvider)
	steps.RegisterFrameworkItems(
		steps.Item{
			Title: "Echo",
			Desc:  "High performance, extensible, minimalist Go web framework",
		},
	)
}
