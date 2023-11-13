// Package steps provides utility for creating
// each step of the CLI
package steps

import (
	"cmp"
	"slices"

	textinput "github.com/melkeydev/go-blueprint/cmd/ui/textinput"
)

// A StepSchema contains the data that is used
// for an individual step of the CLI
type StepSchema struct {
	StepName string  // The name of a given step
	Options  []Item  // The slice of each option for a given step
	Headers  string  // The title displayed at the top of a given step
	Field    *string // The pointer to the string to be overwritten with the selected Item
}

// Steps contains a slice of steps
type Steps struct {
	Steps []StepSchema
}

// An Item contains the data for each option
// in a StepSchema.Options
type Item struct {
	Title, Desc string
}

// Options contains the name and type of the created project
type Options struct {
	ProjectName *textinput.Output
	ProjectType string
}

var registeredFramworkOptions = []Item{}

// InitSteps initializes and returns the *Steps to be used in the CLI program
func InitSteps(options *Options) *Steps {
	steps := &Steps{
		[]StepSchema{
			{
				StepName: "Go Project Framework",
				Options: []Item{
					{
						Title: "Standard library",
						Desc:  "The built-in Go standard library HTTP package",
					},
					{
						Title: "Chi",
						Desc:  "A lightweight, idiomatic and composable router for building Go HTTP services",
					},
					{
						Title: "Gin",
						Desc:  "Features a martini-like API with performance that is up to 40 times faster thanks to httprouter",
					},
					{
						Title: "Fiber",
						Desc:  "An Express inspired web framework built on top of Fasthttp",
					},
					{
						Title: "Gorilla/Mux",
						Desc:  "Package gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler",
					},
					{
						Title: "HttpRouter",
						Desc:  "HttpRouter is a lightweight high performance HTTP request router for Go",
					},
					{Title: "Echo",
						Desc: "High performance, extensible, minimalist Go web framework",
					},
					{Title: "Caddy",
						Desc: "Fast and extensible multi-platform HTTP/1-2-3 web server with automatic HTTPS",
					},
				},
				Headers: "What framework do you want to use in your Go project?",
				Field:   &options.ProjectType,
			},
		},
	}

	return steps
}

func GetSteps(options *Options) *Steps {
	items := make([]Item, len(registeredFramworkOptions))
	copy(items, registeredFramworkOptions)
	steps := &Steps{
		[]StepSchema{
			{
				StepName: "Go Project Framework",
				Options:  items,
				Headers:  "What framework do you want to use in your Go project?",
				Field:    &options.ProjectType,
			},
		},
	}

	return steps
}

func RegisterFrameworkItems(items ...Item) {
	registeredFramworkOptions = append(registeredFramworkOptions, items...)
	slices.SortFunc(registeredFramworkOptions, func(a, b Item) int {
		return cmp.Compare(a.Title, b.Title)
	})

}
