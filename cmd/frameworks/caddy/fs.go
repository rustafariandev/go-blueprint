package caddy

import (
	"embed"

	"github.com/melkeydev/go-blueprint/cmd/frameworks/common"
	"github.com/melkeydev/go-blueprint/cmd/registry"
	"github.com/yalue/merged_fs"
)

//go:embed blueprint.toml internal cmd all:*.tmpl
var fs embed.FS

func init() {
	err := registry.RegisterProviderFromFS(merged_fs.NewMergedFS(fs, common.FS))
	if err != nil {
		panic(err.Error())
	}
}
