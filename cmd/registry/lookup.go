package registry

import (
	"cmp"
	"errors"
	"io/fs"
	"slices"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/ui/inputoptions"
)

type TemplateProviderFunc func() *provider.TemplateProvider

type lookupTemplateProvider struct {
	lock   sync.RWMutex
	lookup map[string]TemplateProviderFunc
	items  []*inputoptions.Item
}

func newLookupTemplateProvider() *lookupTemplateProvider {
	return &lookupTemplateProvider{
		lookup: make(map[string]TemplateProviderFunc),
		items:  make([]*inputoptions.Item, 0),
	}
}

func (l *lookupTemplateProvider) Get(name string) (*provider.TemplateProvider, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if f, ok := l.lookup[name]; ok {
		return f(), nil
	}

	return nil, errors.New("Provider Not Found")
}

func (l *lookupTemplateProvider) Set(name string, f TemplateProviderFunc, item *inputoptions.Item) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lookup[name] = f
	l.items = append(l.items, item)
	slices.SortFunc(l.items, func(a, b *inputoptions.Item) int {
		return cmp.Compare(a.Name, b.Name)
	})
}

func (l *lookupTemplateProvider) Has(name string) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	_, b := l.lookup[name]
	return b
}

func (l *lookupTemplateProvider) GetItems() []*inputoptions.Item {
	l.lock.Lock()
	defer l.lock.Unlock()
	items := make([]*inputoptions.Item, len(l.items))
	copy(items, l.items)
	return items
}

var frameworks = newLookupTemplateProvider()
var addons = newLookupTemplateProvider()

// InitSteps initializes and returns the *Steps to be used in the CLI program

func GetFrameworkItems() []*inputoptions.Item {
	return frameworks.GetItems()
}

func GetAddonItems() []*inputoptions.Item {
	return addons.GetItems()
}

func GetFramework(name string) (*provider.TemplateProvider, error) {
	return frameworks.Get(name)
}

func GetAddon(name string) (*provider.TemplateProvider, error) {
	return addons.Get(name)
}

type BlueprintConfig struct {
	Name        string   `toml:"name"`
	Title       string   `toml:"title"`
	Group       string   `toml:"group"`
	Description string   `toml:"description"`
	Type        string   `toml:"type"`
	Packages    []string `toml:"packages"`
}

func HasFramework(name string) bool {
	return frameworks.Has(name)
}

func HasAddon(name string) bool {
	return addons.Has(name)
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

	f := func() *provider.TemplateProvider {
		return &provider.TemplateProvider{
			TempateFS:    filesys,
			PackageNames: config.Packages,
			ProjectType:  config.Name,
		}
	}

	item := &inputoptions.Item{
		Name:  config.Title,
		Desc:  config.Description,
		Value: config.Name,
	}

	if config.Type == "addon" {
		addons.Set(config.Name, f, item)
	} else {
		frameworks.Set(config.Name, f, item)
	}

	return nil
}
