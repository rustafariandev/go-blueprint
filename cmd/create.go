package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/registry"
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
	ListAddons     bool
	ProjectName    string
	Framework      string
	Output         string
	Addons         []string
}

func (o *CreateOptions) DoListing() bool {
	return o.ListFrameworks || o.ListAddons
}

func (o *CreateOptions) SetFlags(f *pflag.FlagSet) {
	f.StringVarP(&o.ProjectName, "name", "n", o.ProjectName, "Name of project to create")
	f.StringVarP(&o.Output, "output", "o", o.Output, "Ouput")
	f.StringVarP(&o.Framework, "framework", "f", o.Framework, "Framework to use - To see aviailable options use create -l")
	f.StringArrayVarP(&o.Addons, "addon", "a", o.Addons, "Addon to use can be used multiple times - To see aviailable options use create -A")
	f.BoolVarP(&o.ListFrameworks, "list", "l", o.ListFrameworks, "List aviailable frameworks")
	f.BoolVarP(&o.ListAddons, "list-addons", "A", o.ListFrameworks, "List aviailable addons")
}

func init() {
	rootCmd.AddCommand(createCmd)
	options.SetFlags(createCmd.Flags())
}

func (o *CreateOptions) List() {

	if o.ListFrameworks {
		fmt.Printf("Frameworks\n")
		for _, i := range registry.GetFrameworkItems() {
			fmt.Printf("* %s\n", i.Value)
		}
	}

	if o.ListAddons {
		fmt.Printf("Addons\n")
		for _, i := range registry.GetAddonItems() {
			fmt.Printf("* %s\n", i.Value)
		}
	}

}

func (o *CreateOptions) Verify() error {
	if o.Framework == "" {
		return nil
	}

	if !registry.HasFramework(o.Framework) {
		return fmt.Errorf("Framework %s is not invalid", o.Framework)
	}

	var addonsNotFound = []string{}
	for _, a := range o.Addons {
		if !registry.HasAddon(a) {
			addonsNotFound = append(addonsNotFound, a)
		}
	}
	if len(addonsNotFound) > 0 {
		return fmt.Errorf("Addons %s not found", strings.Join(addonsNotFound, ", "))
	}

	return nil
}

func (o *CreateOptions) AskForOptions() bool {
	return o.Framework == "" || o.ProjectName == ""
}

func (o *CreateOptions) GetModelOption() inputoptions.ModelOptions {
	options := inputoptions.ModelOptions{
		Items:      registry.GetFrameworkItems(),
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

		if options.DoListing() {
			options.List()
			os.Exit(0)
		}

		if err := options.Verify(); err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		flagFramework := cmd.Flag("framework").Value.String()
		flagOutputDir := cmd.Flag("output").Value.String()
		_ = flagOutputDir

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

		if options.Output == "" {
			currentWorkingDir, err := os.Getwd()
			if err != nil {
				log.Printf("could not get current working directory: %v", err)
				cobra.CheckErr(err)
			}

			options.Output = filepath.Join(currentWorkingDir, options.ProjectName)
		}

		if err := tp.Create(
			&provider.Project{
				ProjectName:  options.ProjectName,
				AbsolutePath: options.Output,
				ProjectType:  options.Framework,
				PackageNames: tp.PackageNames,
			},
			&provider.RunOptions{
				SkipInitialization: false,
			},
		); err != nil {
			log.Printf("Problem creating files for project. %v", err)
			cobra.CheckErr(err)
		}

		for _, a := range options.Addons {
			tp, err := registry.GetAddon(a)
			if err != nil {
				log.Printf("Problem getting a. %v", err)
				cobra.CheckErr(err)
			}
			if err := tp.Create(
				&provider.Project{
					ProjectName:  options.ProjectName,
					AbsolutePath: options.Output,
					ProjectType:  options.Framework,
					PackageNames: tp.PackageNames,
				},
				&provider.RunOptions{
					SkipInitialization: true,
				},
			); err != nil {
				log.Printf("Problem creating files for project. %v", err)
				cobra.CheckErr(err)
			}
		}

		fmt.Println(endingMsgStyle.Render("\nNext steps cd into the newly created project with:"))
		fmt.Println(endingMsgStyle.Render(fmt.Sprintf("• cd %s\n", options.Output)))

		if isInteractive {
			nonInteractiveCommand := utils.NonInteractiveCommand(cmd.Flags())
			fmt.Println(tipMsgStyle.Render("Tip: Repeat the equivalent Blueprint with the following non-interactive command:"))
			fmt.Println(tipMsgStyle.Italic(false).Render(fmt.Sprintf("• %s\n", nonInteractiveCommand)))
		}
	},
}
