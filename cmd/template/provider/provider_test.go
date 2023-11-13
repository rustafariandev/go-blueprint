package provider

import (
	"embed"
	"fmt"
	"os"
	"testing"
)

//go:embed internal/*
var test_fs embed.FS

func TestProvider(t *testing.T) {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	//	defer os.RemoveAll(dir)

	tp := TemplateProvider{TempateFS: test_fs}
	p := &Project{ProjectName: "test", AbsolutePath: dir}
	fmt.Printf("create dir %s\n", dir)
	tp.Create(p)
}
