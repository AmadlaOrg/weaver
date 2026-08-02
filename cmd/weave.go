package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/AmadlaOrg/weaver/weave"
	"github.com/spf13/cobra"
)

var (
	// Cmd flags package global variable
	templatePath string
	outputPath   string
	entityPath   string

	// Functions
	osOpen   = os.Open
	osCreate = os.Create
	weaveNew = weave.New

	// WeaveCmd sets up the weave command
	WeaveCmd = &cobra.Command{
		Use:   "weave",
		Short: "From entity to output using a template",
		//Long:  `Execute the weave process using a specified template and data`,
		RunE: runWeave,
	}
)

func init() {
	WeaveCmd.Flags().StringVarP(
		&templatePath,
		"template",
		"t",
		"",
		"Specify the template file path (required)",
	)
	WeaveCmd.Flags().StringVarP(
		&outputPath,
		"output",
		"o",
		"",
		"Specify the output file path (optional, defaults to stdout)",
	)
	WeaveCmd.Flags().StringVarP(
		&entityPath,
		"entity",
		"e",
		"",
		"Specify the entity file path (optional)",
	)
	_ = WeaveCmd.MarkFlagRequired("template")
}

// runWeave
func runWeave(cmd *cobra.Command, args []string) error {
	// 1. Handle input
	//
	// - By default the input is `os.Stdin`
	var input io.Reader = os.Stdin
	if entityPath != "" {
		// 1.1 Opens the entity path given
		entityFile, err := osOpen(entityPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("entity file does not exist: %s", entityPath)
			}
			return fmt.Errorf("failed to open entity file: %w", err)
		}
		defer func(file *os.File) {
			if err := file.Close(); err != nil {
				cmd.PrintErrf("Failed to close entity file: %v\n", err)
			}
		}(entityFile)

		// 1.2 Sets the input variable with the `os.File`
		input = entityFile
	}

	// 2. Handle output
	//
	// - By default the output is `os.Stdout`
	// - The output flag is a path to a file; the directory in the path must exist
	var output io.Writer = os.Stdout
	if outputPath != "" {
		// 2.1 Create the file for the output
		outputFile, err := osCreate(outputPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("output directory does not exist: %s", outputPath)
			}
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func(file *os.File) {
			if err := file.Close(); err != nil {
				cmd.PrintErrf("Failed to close output file: %v\n", err)
			}
		}(outputFile)

		// 2.2 Sets the output variable to the `os.File`
		output = outputFile
	}

	// 3. Execute weaving process
	if err := weaveNew(templatePath, input, output).Do(); err != nil {
		return fmt.Errorf("weave process failed: %w", err)
	}

	return nil
}
