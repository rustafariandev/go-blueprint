package provider

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/melkeydev/go-blueprint/cmd/utils"
	"github.com/spf13/cobra"
)

type TemplateProvider struct {
	TempateFS    embed.FS
	PackageNames []string
	ProjectType  string
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
	return os.Create(filepath.Join(p.AbsolutePath, p.ProjectName, path))
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

	projectPath := filepath.Join(p.AbsolutePath, p.ProjectName)
	if err := os.Mkdir(projectPath, 0754); err != nil {
		log.Printf("Could not create directory: %v", err)
		return err
	}
	// Create go.mod
	err := utils.InitGoMod(p.ProjectName, projectPath)
	if err != nil {
		log.Printf("Could not initialize go.mod in new project %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, p.PackageNames)
	if err != nil {
		log.Printf("Could not install go dependency for the chosen framework %v\n", err)
		cobra.CheckErr(err)
	}

	fs.WalkDir(
		tp.TempateFS,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				if path == "." {
					return err
				}

				newDir := filepath.Join(projectPath, path)
				if err := os.Mkdir(newDir, 0754); err != nil {
					log.Printf("Could not create directory %s %s: %v", path, newDir, err)
					return err
				}
				return err
			}

			if strings.HasSuffix(path, ".tmpl") {
				err := tp.CreateFileFromTemplate(p, path)
				if err != nil {
					return err
				}

				return nil
			}

			return tp.CopyFile(p, path)
		},
	)

	// Initialize git repo
	err = utils.ExecuteCmd("git", []string{"init"}, projectPath)
	if err != nil {
		log.Printf("Error initializing git repo: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = utils.GoModTidy(projectPath)
	if err != nil {
		log.Printf("Could not go mod tidy the new project %v\n", err)
		cobra.CheckErr(err)
		return err
	}

	err = utils.GoFmt(projectPath)
	if err != nil {
		log.Printf("Could not gofmt in new project %v\n", err)
		cobra.CheckErr(err)
		return err
	}

	return nil
}

func (tp *TemplateProvider) CopyFile(p *Project, path string) error {
	createdFile, err := p.CreateFile(path)
	if err != nil {
		return err
	}
	defer createdFile.Close()

	data, err := tp.TempateFS.ReadFile(path)
	if err != nil {
		return err
	}

	createdFile.Write(data)
	return err
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

	createdTemplate := template.Must(template.New(templ).Parse(string(data)))
	return createdTemplate.Execute(createdFile, p)
}
