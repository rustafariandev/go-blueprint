package registry

import (
	"errors"
	"io/fs"

	"github.com/BurntSushi/toml"
	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/steps"
	"github.com/melkeydev/go-blueprint/cmd/ui/inputoptions"
)

var lookup = map[string]func() *provider.TemplateProvider{}

func GetFramework(name string) (*provider.TemplateProvider, error) {
	if f, ok := lookup[name]; ok {
		return f(), nil
	}

	return nil, errors.New("Provider Not Found")
}

type BlueprintConfig struct {
	Name        string   `toml:"name"`
	Title       string   `toml:"title"`
	Group       string   `toml:"group"`
	Description string   `toml:"description"`
	Type        string   `toml:"type"`
	Packages    []string `toml:"packages"`
}

func RegisterFramework(name string, f func() *provider.TemplateProvider) {
	lookup[name] = f
}

func HasFramework(name string) bool {
	_, b := lookup[name]
	return b
}

func RegisterProviderFromFS(filesys fs.FS) error {
	bytes, err := fs.ReadFile(filesys, "blueprint.toml")
	if err != nil {
		return err
	}

	config := BlueprintConfig{}
	_, err = toml.Decode(string(bytes), &config)
	if err != nil {
		return err
	}

	lookup[config.Name] = func() *provider.TemplateProvider {
		return &provider.TemplateProvider{
			TempateFS:    filesys,
			PackageNames: config.Packages,
			ProjectType:  config.Name,
		}
	}

	steps.RegisterFrameworkItems(
		&inputoptions.Item{
			Name:  config.Title,
			Desc:  config.Description,
			Value: config.Name,
		},
	)

	return nil
}
