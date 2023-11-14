package cobra

import (
	"embed"

	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/registry"
	"github.com/melkeydev/go-blueprint/cmd/steps"
)

//go:embed cmd all:*.tmpl
var FS embed.FS

func MakeProvider() *provider.TemplateProvider {
	return &provider.TemplateProvider{
		TempateFS:    FS,
		PackageNames: []string{"github.com/spf13/cobra"},
		ProjectType:  "cobra",
	}
}

func init() {
	registry.RegisterFramework("cobra", MakeProvider)
	steps.RegisterFrameworkItems(
		steps.Item{
			Title: "Cobra",
			Desc:  "A library for creating powerful modern CLI applications",
		},
	)
}
