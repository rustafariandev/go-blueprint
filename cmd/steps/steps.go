// Package steps provides utility for creating
// each step of the CLI
package steps

import (
	"cmp"
	"slices"

	"github.com/melkeydev/go-blueprint/cmd/ui/inputoptions"
)

var registeredFrameworkItems = []*inputoptions.Item{}

// InitSteps initializes and returns the *Steps to be used in the CLI program

func GetItems() []*inputoptions.Item {
	return registeredFrameworkItems
}

func RegisterFrameworkItems(items ...*inputoptions.Item) {
	registeredFrameworkItems = append(registeredFrameworkItems, items...)
	slices.SortFunc(registeredFrameworkItems, func(a, b *inputoptions.Item) int {
		return cmp.Compare(a.Name, b.Name)
	})

}
