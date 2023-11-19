package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/melkeydev/go-blueprint/cmd/frameworks/caddy"
	_ "github.com/melkeydev/go-blueprint/cmd/frameworks/chi"
	_ "github.com/melkeydev/go-blueprint/cmd/frameworks/cobra"
	_ "github.com/melkeydev/go-blueprint/cmd/frameworks/echo"
	_ "github.com/melkeydev/go-blueprint/cmd/frameworks/fiber"
	_ "github.com/melkeydev/go-blueprint/cmd/frameworks/gin"
	_ "github.com/melkeydev/go-blueprint/cmd/frameworks/gorilla"
	_ "github.com/melkeydev/go-blueprint/cmd/frameworks/httprouter"
	_ "github.com/melkeydev/go-blueprint/cmd/frameworks/standard-library"
	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/registry"
	"github.com/melkeydev/go-blueprint/cmd/steps"
	"github.com/melkeydev/go-blueprint/cmd/ui/inputoptions"
	"github.com/melkeydev/go-blueprint/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const logo = `

 ____  _                       _       _   
|  _ \| |                     (_)     | |  
| |_) | |_   _  ___ _ __  _ __ _ _ __ | |_ 
|  _ <| | | | |/ _ \ '_ \| '__| | '_ \| __|
| |_) | | |_| |  __/ |_) | |  | | | | | |_ 
|____/|_|\__,_|\___| .__/|_|  |_|_| |_|\__|
                                  | |
                                  |_|

`

var (
	logoStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	tipMsgStyle         = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("190")).Italic(true)
	endingMsgStyle      = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
	allowedProjectTypes = []string{"chi", "gin", "fiber", "gorilla/mux", "httprouter", "standard-library", "echo", "caddy", "cobra"}
	options             = CreateOptions{}
)

type CreateOptions struct {
	ListFrameworks bool
	ProjectName    string
	Framework      string
}

func (o *CreateOptions) SetFlags(f *pflag.FlagSet) {
	f.StringVarP(&o.ProjectName, "name", "n", o.ProjectName, "Name of project to create")
	f.StringVarP(&o.Framework, "framework", "f", o.Framework, "Framework to use - To see aviailable options use create -l")
	f.BoolVarP(&o.ListFrameworks, "list", "l", o.ListFrameworks, "List aviailable frameworks")
}

func init() {
	rootCmd.AddCommand(createCmd)
	options.SetFlags(createCmd.Flags())
}

func (o *CreateOptions) ListAllowFrameworks() {
	for _, i := range steps.GetItems() {
		fmt.Printf("%s\n", i.Value)
	}

	os.Exit(0)
}

func (o *CreateOptions) Verify() error {
	if o.Framework == "" {
		return nil
	}

	if !registry.HasFramework(o.Framework) {
		return fmt.Errorf("Framework %s is not invalid", o.Framework)
	}

	return nil
}

func (o *CreateOptions) AskForOptions() bool {
	return o.Framework == "" || o.ProjectName == ""
}

func (o *CreateOptions) GetModelOption() inputoptions.ModelOptions {
	options := inputoptions.ModelOptions{
		Items:      steps.GetItems(),
		Header:     "What is the name of your project?",
		ListHeader: "What framework do you want to use in your Go project?",
	}

	if o.Framework != "" {
		options.SkipList = true
	}

	if o.ProjectName != "" {
		options.ShowList = true
	}
	return options
}

// createCmd defines the "create" command for the CLI
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Go project and don't worry about the structure",
	Long:  "Go Blueprint is a CLI tool that allows you to focus on the actual Go code, and not the project structure. Perfect for someone new to the Go language",

	Run: func(cmd *cobra.Command, args []string) {

		isInteractive := !utils.HasChangedFlag(cmd.Flags())

		if options.ListFrameworks {
			options.ListAllowFrameworks()
		}

		if err := options.Verify(); err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		flagName := cmd.Flag("name").Value.String()
		flagFramework := cmd.Flag("framework").Value.String()

		ProjectName := flagName
		ProjectType := flagFramework

		if options.AskForOptions() {
			modelOptions := options.GetModelOption()
			model1 := inputoptions.NewModel(modelOptions)
			p := tea.NewProgram(model1, tea.WithAltScreen())
			m, err := p.Run()
			if err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}

			model2 := m.(inputoptions.Model)
			if model2.Quit() {
				if err := p.ReleaseTerminal(); err != nil {
					log.Fatal(err)
				}
				os.Exit(1)
			}

			output := model2.GetOutput()
			if options.Framework == "" {
				options.Framework = output.Framework
			}

			if options.ProjectName == "" {
				options.ProjectName = output.Name
			}
		}

		tp, err := registry.GetFramework(ProjectType)
		if err != nil {
			log.Printf("Problem getting framework. %v", err)
			cobra.CheckErr(err)
		}

		currentWorkingDir, err := os.Getwd()
		if err != nil {
			log.Printf("could not get current working directory: %v", err)
			cobra.CheckErr(err)
		}

		err = tp.Create(&provider.Project{
			ProjectName:  options.ProjectName,
			AbsolutePath: currentWorkingDir,
			ProjectType:  options.Framework,
			PackageNames: tp.PackageNames,
		})

		if err != nil {
			log.Printf("Problem creating files for project. %v", err)
			cobra.CheckErr(err)
		}

		fmt.Println(endingMsgStyle.Render("\nNext steps cd into the newly created project with:"))
		fmt.Println(endingMsgStyle.Render(fmt.Sprintf("• cd %s\n", ProjectName)))

		if isInteractive {
			nonInteractiveCommand := utils.NonInteractiveCommand(cmd.Flags())
			fmt.Println(tipMsgStyle.Render("Tip: Repeat the equivalent Blueprint with the following non-interactive command:"))
			fmt.Println(tipMsgStyle.Italic(false).Render(fmt.Sprintf("• %s\n", nonInteractiveCommand)))
		}
	},
}
