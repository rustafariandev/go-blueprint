package provider

import (
	"embed"
	"fmt"
	"os"
	"testing"
)

//go:embed internal/* TEST
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
	err = tp.Create(p, &RunOptions{})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
