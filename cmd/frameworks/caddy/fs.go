package caddy

import (
	"embed"

	"github.com/melkeydev/go-blueprint/cmd/registry"
)

//go:embed blueprint.toml internal cmd all:*.tmpl
var fs embed.FS

func init() {
	err := registry.RegisterProviderFromFS(fs)
	if err != nil {
		panic(err.Error())
	}
}
