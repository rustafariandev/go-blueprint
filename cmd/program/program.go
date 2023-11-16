// Package program provides the
// main functionality of Blueprint
package program

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// A Project contains the data for the project folder
// being created, and methods that help with that process
type Project struct {
	ProjectName  string
	Exit         bool
	AbsolutePath string
	ProjectType  string
}

// A Framework contains the name and templater for a
// given Framework
type Framework struct {
	packageName []string
	templater   Templater
}

// A Templater has the methods that help build the files
// in the Project folder, and is specific to a Framework
type Templater interface {
	Main() []byte
	Server() []byte
	Routes() []byte
	Plugin() []byte
}

var (
	chiPackage     = []string{"github.com/go-chi/chi/v5"}
	gorillaPackage = []string{"github.com/gorilla/mux"}
	routerPackage  = []string{"github.com/julienschmidt/httprouter"}
	ginPackage     = []string{"github.com/gin-gonic/gin"}
	fiberPackage   = []string{"github.com/gofiber/fiber/v2"}
	echoPackage    = []string{"github.com/labstack/echo/v4", "github.com/labstack/echo/v4/middleware"}
	caddyPackage   = []string{"github.com/caddyserver/caddy/v2", "github.com/caddyserver/caddy/v2/cmd"}

	cmdApiPath         = "cmd/api"
	internalServerPath = "internal/server"
	internalPluginPath = "internal/plugin"
)

// ExitCLI checks if the Project has been exited, and closes
// out of the CLI if it has
func (p *Project) ExitCLI(tprogram *tea.Program) {
	if p.Exit {
		// logo render here
		if err := tprogram.ReleaseTerminal(); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}

func (p *Project) RemoveEmptyDirs(projectPath string, dirs ...string) {
	for _, d := range dirs {
		path := filepath.Join(projectPath, d)
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		_, err = f.Readdirnames(1)
		f.Close()
		if err == io.EOF {
			os.Remove(path)
		}
	}

}

// CreatePath creates the given directory in the projectPath
func (p *Project) CreatePath(pathToCreate string, projectPath string) error {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", projectPath, pathToCreate)); os.IsNotExist(err) {
		err := os.MkdirAll(fmt.Sprintf("%s/%s", projectPath, pathToCreate), 0751)
		if err != nil {
			log.Printf("Error creating directory %v\n", err)
			return err
		}
	}

	return nil
}
