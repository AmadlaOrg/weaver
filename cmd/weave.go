package cmd

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/LibraryUtils/file"
	"github.com/AmadlaOrg/weaver/weave"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

var (
	// Cmd flags package global variable
	templatePath string
	outputPath   string
	entityPath   string

	// Functions
	osOpen               = os.Open
	osCreate             = os.Create
	osStat               = os.Stat
	fileIsFile           = file.IsFile
	weaveNewWeaveService = weave.NewWeaveService

	// WeaveCmd sets up the weave command
	WeaveCmd = &cobra.Command{
		Use:   "weave",
		Short: "Execute the weave process using a specified template and data",
		Long:  `Execute the weave process using a specified template and data`,
		Run: func(cmd *cobra.Command, args []string) {

			// 1. Setup of the `weave` flags
			cmd.Flags().StringVarP(
				&templatePath,
				"template",
				"t",
				"",
				"Specify the template file path (required)",
			)
			cmd.Flags().StringVarP(
				&outputPath,
				"output",
				"o",
				"",
				"Specify the output file path (optional, defaults to stdout)",
			)
			cmd.Flags().StringVarP(
				&entityPath,
				"entity",
				"e",
				"",
				"Specify the entity file path (optional)",
			)

			// 2. The template flag is required
			err := cmd.MarkFlagRequired("template")
			if err != nil {
				fmt.Println(err)
				return
			}

			// 3. Validates what was passed in the template flag
			if _, err := fileIsFile(templatePath); err != nil {
				if errors.Is(err, file.ErrorIsDir) {
					cmd.Println("Template path given is a directory, expected a template file.")
				} else {
					cmd.Println("The template file does not exist.")
				}
				return
			}

			// 4. Handle input
			//
			// - By default the input is `os.Stdin`
			// -
			var input io.Reader = os.Stdin
			if entityPath != "" {

				// 4.1 Validates if the file exist and that the path is not just a directory
				if _, err := fileIsFile(entityPath); err != nil {
					if errors.Is(err, file.ErrorIsDir) {
						cmd.Println("Entity path given is a directory, expected a template file.")
					} else {
						cmd.Println("The entity file does not exist.")
					}
					return
				}

				// 4.2 Opens the entity path given
				file, err := osOpen(entityPath)
				if err != nil {
					cmd.Println("Failed to open entity file: %v\n", err)
					return
				}
				defer func(file *os.File) {
					err := file.Close()
					if err != nil {
						cmd.Println("Failed to close entity file: %v\n", err)
					}
				}(file)

				// 4.3 Sets the input variable with the `os.File`
				input = file
			}

			// 5. Handle output
			//
			// - By default the output is `os.Stdout`
			// - The output flag is a path to a file that does not exist (that needs to be created), but it validates that the directory in the path given is valid
			// - If the file already exist then it asks if the existing output file can be overwritten
			var output io.Writer = os.Stdout
			if outputPath != "" {
				// 5.1 Extracts the base of file path (gives the path without the file name)
				outputBasePath := filepath.Base(outputPath)

				// 5.2 Validates that the path is of an existing directory
				info, err := osStat(outputPath)
				if err != nil {
					cmd.Println("Failed to stat output directory: %v\n", err)
				}
				if !info.IsDir() {
					cmd.Println("Output directory path does not exists: %v\n", outputBasePath)
					return
				}

				// 5.3 Create the file for the output
				file, err := osCreate(outputPath)
				if err != nil {
					cmd.Println("Failed to create output file: %v\n", err)
					return
				}
				defer func(file *os.File) {
					err := file.Close()
					if err != nil {
						cmd.Println("Failed to close output file: %v\n", err)
					}
				}(file)

				// 5.4 Sets the output variable to the `os.File`
				output = file
			}

			// 6. Execute weaving process
			err = weaveNewWeaveService(templatePath, input, output).Do()
			if err != nil {
				cmd.Println("Weave process failed: %v\n", err)
			}
		},
	}
)
