package provider

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateProvider struct {
	TempateFS    embed.FS
	PackageNames []string
}

// A Project contains the data for the project folder
// being created, and methods that help with that process
type Project struct {
	ProjectName  string
	AbsolutePath string
	ProjectType  string
	PackageNames []string
}

func (p *Project) CreateFile(path string) (*os.File, error) {
	os.Create(filepath.Join(p.Absolute, path))
}

func (tp *TemplateProvider) Create(p *Project) error {
	// check if s
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			log.Printf("Could not create directory: %v", err)
			return err
		}
	}

	fs.WalkDir(
		tp.TempateFS,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				newDir := filepath.Join(p.AbsolutePath, path)
				if err := os.Mkdir(newDir, 0754); err != nil {
					log.Printf("Could not create directory %s: %v", newDir, err)
					return err
				}
				return err
			}
			if strings.HasSuffix(path, ".tmpl") {
				err := tp.CreateFileFromTemplate(p, path)
				if err != nil {
					return err
				}
			}
			return
		},
	)

	return nil
}

func (tp *TemplateProvider) CreateFileFromTemplate(p *Project, templ string) error {
	path := strings.TrimSuffix(templ, ".tmpl")
	createdFile, err := p.CreateFile(path)
	if err != nil {
		return err
	}

	defer createdFile.Close()
	data, err := tp.TempateFS.ReadFile(templ)
	if err != nil {
		return err
	}

	createdTemplate := template.Must(data)
	return createdTemplate.Execute(createdFile, p)
}
