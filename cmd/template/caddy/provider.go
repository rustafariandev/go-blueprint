package caddy

import (
	"embed"

	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/template"
)

//go:embed internal cmd all:*.tmpl
var fs embed.FS

func MakeProvider() *provider.TemplateProvider {
	return &provider.TemplateProvider{
		TempateFS: fs,
		PackageNames: []string{
			"github.com/caddyserver/caddy/v2",
			"github.com/caddyserver/caddy/v2/cmd",
		},
		ProjectType: "caddy",
	}
}

func init() {
	template.RegisterProvider("caddy", MakeProvider)
}
