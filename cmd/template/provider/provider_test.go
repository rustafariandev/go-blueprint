package provider

import (
	"embed"
	"testing"
)

//go:embed internal/*
var test_fs embed.FS

func TestProvider(t *testing.T) {
	tp := TemplateProvider{FS: test_fs}
	tp.Create("bob", "bob")
}
