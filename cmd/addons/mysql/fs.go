package mysql

import (
	"embed"

	"github.com/melkeydev/go-blueprint/cmd/registry"
)

//go:embed blueprint.toml client
var fs embed.FS

func init() {
	err := registry.RegisterProviderFromFS(fs)
	if err != nil {
		panic(err.Error())
	}
}
