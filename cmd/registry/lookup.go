package registry

import (
	"embed"
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/steps"
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
	Packages    []string `toml:"packages"`
}

func RegisterFramework(name string, f func() *provider.TemplateProvider) {
	lookup[name] = f
}

func RegisterProviderFromFS(fs embed.FS) error {
	bytes, err := fs.ReadFile("blueprint.toml")
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
			TempateFS:    fs,
			PackageNames: config.Packages,
			ProjectType:  config.Name,
		}
	}

	steps.RegisterFrameworkItems(
		steps.Item{
			Title: config.Title,
			Desc:  config.Description,
		},
	)

	return nil
}
