package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadlaOrg/weaver/plugin"
	"github.com/spf13/cobra"
)

var (
	renderTemplatePath string
	renderFilePath     string
	renderOutputPath   string
	renderEngine       string

	pluginNew = plugin.New

	// RenderCmd delegates template rendering to the appropriate plugin.
	RenderCmd = &cobra.Command{
		Use:   "render",
		Short: "Render a template using a discovered plugin",
		Long:  "Detects the template engine from the file extension or --engine flag and delegates to the matching weaver-* plugin.",
		RunE:  runRender,
	}
)

func init() {
	RenderCmd.Flags().StringVarP(&renderTemplatePath, "template", "t", "", "Path to the template file (required)")
	RenderCmd.Flags().StringVarP(&renderFilePath, "file", "f", "", "Path to data file (JSON or YAML; default: stdin)")
	RenderCmd.Flags().StringVarP(&renderOutputPath, "output", "o", "", "Path to output file (default: stdout)")
	RenderCmd.Flags().StringVarP(&renderEngine, "engine", "e", "", "Plugin engine name (e.g. jinja2, go, mustache, qute, freemarker)")
	_ = RenderCmd.MarkFlagRequired("template")
}

func runRender(cmd *cobra.Command, args []string) error {
	svc := pluginNew()

	var pluginName string

	if renderEngine != "" {
		// Explicit engine: look for weaver-<engine>
		pluginName = "weaver-" + renderEngine
	} else {
		// Auto-detect from template file extension
		ext := filepath.Ext(renderTemplatePath)
		if ext == "" {
			return fmt.Errorf("cannot detect template engine: no file extension on %q (use --engine to specify)", renderTemplatePath)
		}

		name, err := svc.FindByExtension(ext)
		if err != nil {
			return fmt.Errorf("no plugin found for extension %q (use --engine to specify): %w", ext, err)
		}
		pluginName = name
	}

	err := svc.Render(pluginName, renderTemplatePath, renderFilePath, renderOutputPath, os.Stdin)
	if err != nil {
		return err
	}

	return nil
}
