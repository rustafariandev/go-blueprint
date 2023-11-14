package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/melkeydev/go-blueprint/cmd/program"
	"github.com/melkeydev/go-blueprint/cmd/provider"
	"github.com/melkeydev/go-blueprint/cmd/registry"
	"github.com/melkeydev/go-blueprint/cmd/steps"
	_ "github.com/melkeydev/go-blueprint/cmd/template/caddy"
	_ "github.com/melkeydev/go-blueprint/cmd/template/chi"
	_ "github.com/melkeydev/go-blueprint/cmd/template/cobra"
	_ "github.com/melkeydev/go-blueprint/cmd/template/echo"
	_ "github.com/melkeydev/go-blueprint/cmd/template/fiber"
	_ "github.com/melkeydev/go-blueprint/cmd/template/gin"
	_ "github.com/melkeydev/go-blueprint/cmd/template/gorilla"
	_ "github.com/melkeydev/go-blueprint/cmd/template/httprouter"
	_ "github.com/melkeydev/go-blueprint/cmd/template/standard-library"
	"github.com/melkeydev/go-blueprint/cmd/ui/multiInput"
	"github.com/melkeydev/go-blueprint/cmd/ui/textinput"
	"github.com/melkeydev/go-blueprint/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd2)

	createCmd2.Flags().StringP("name", "n", "", "Name of project to create")
	createCmd2.Flags().StringP("framework", "f", "", fmt.Sprintf("Framework to use. Allowed values: %s", strings.Join(allowedProjectTypes, ", ")))
}

// createCmd2 defines the "create2" command for the CLI
var createCmd2 = &cobra.Command{
	Use:   "create2",
	Short: "Create a Go project and don't worry about the structure",
	Long:  "Go Blueprint is a CLI tool that allows you to focus on the actual Go code, and not the project structure. Perfect for someone new to the Go language",

	Run: func(cmd *cobra.Command, args []string) {
		var tprogram *tea.Program

		options := steps.Options{
			ProjectName: &textinput.Output{},
		}

		isInteractive := !utils.HasChangedFlag(cmd.Flags())

		flagName := cmd.Flag("name").Value.String()
		flagFramework := cmd.Flag("framework").Value.String()

		project := &program.Project{
			FrameworkMap: make(map[string]program.Framework),
			ProjectName:  flagName,
			ProjectType:  strings.ReplaceAll(flagFramework, "-", " "),
		}
		ProjectName := flagName
		ProjectType := flagFramework

		fmt.Printf("%s\n", logoStyle.Render(logo))

		if ProjectName == "" {
			tprogram := tea.NewProgram(textinput.InitialTextInputModel(options.ProjectName, "What is the name of your project?", project))
			if _, err := tprogram.Run(); err != nil {
				log.Printf("Name of project contains an error: %v", err)
				cobra.CheckErr(err)
			}
			project.ExitCLI(tprogram)

			ProjectName = options.ProjectName.Output
			err := cmd.Flag("name").Value.Set(ProjectName)
			if err != nil {
				log.Fatal("failed to set the name flag value", err)
			}
		}

		steps := steps.GetSteps(&options)
		if ProjectType == "" {
			for _, step := range steps.Steps {
				s := &multiInput.Selection{}
				tprogram = tea.NewProgram(multiInput.InitialModelMulti(step.Options, s, step.Headers, project))
				if _, err := tprogram.Run(); err != nil {
					cobra.CheckErr(err)
				}
				project.ExitCLI(tprogram)

				*step.Field = s.Choice
			}

			ProjectType = options.ProjectType
			err := cmd.Flag("framework").Value.Set(ProjectType)
			if err != nil {
				log.Fatal("failed to set the framework flag value", err)
			}
		}

		tp, err := registry.GetFramework(ProjectType)
		if err != nil {
			log.Printf("Problem getting provider for project. %v", err)
			cobra.CheckErr(err)
		}

		currentWorkingDir, err := os.Getwd()
		if err != nil {
			log.Printf("could not get current working directory: %v", err)
			cobra.CheckErr(err)
		}

		project.AbsolutePath = currentWorkingDir
		err = tp.Create(&provider.Project{
			ProjectName:  ProjectName,
			AbsolutePath: currentWorkingDir,
			ProjectType:  ProjectType,
			PackageNames: tp.PackageNames,
		})

		// This calls the templates
		//		err = project.CreateMainFile()
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
