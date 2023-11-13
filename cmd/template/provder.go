package template

import (
	"errors"

	"github.com/melkeydev/go-blueprint/cmd/provider"
)

var lookup = map[string]func() *provider.TemplateProvider{}

func GetProvider(name string) (*provider.TemplateProvider, error) {
	if f, ok := lookup[name]; ok {
		return f(), nil
	}

	return nil, errors.New("Provider Not Found")
}

func RegisterProvider(name string, f func() *provider.TemplateProvider) {
	lookup[name] = f
}
