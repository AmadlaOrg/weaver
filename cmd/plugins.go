package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AmadlaOrg/weaver/plugin"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	pluginsNew    = plugin.New
	pluginsStdout io.Writer = os.Stdout
	pluginsStderr io.Writer = os.Stderr

	pluginsOutputFlag string
	pluginsHeryFlag   bool

	// PluginsCmd lists all discovered weaver plugins.
	PluginsCmd = &cobra.Command{
		Use:   "plugins",
		Short: "List discovered weaver plugins",
		RunE:  runPlugins,
	}
)

func init() {
	PluginsCmd.Flags().StringVarP(&pluginsOutputFlag, "output", "o", "table", "Output format: table, json, yaml")
	PluginsCmd.Flags().BoolVar(&pluginsHeryFlag, "hery", false, "Wrap output in HERY envelope (_type, _body)")
}

type pluginRow struct {
	Plugin      string `json:"plugin" yaml:"plugin"`
	Engine      string `json:"engine" yaml:"engine"`
	Version     string `json:"version" yaml:"version"`
	Extensions  string `json:"extensions" yaml:"extensions"`
	Description string `json:"description" yaml:"description"`
}

type heryEnvelope struct {
	Type string `json:"_type" yaml:"_type"`
	Body any    `json:"_body" yaml:"_body"`
}

func runPlugins(cmd *cobra.Command, args []string) error {
	svc := pluginsNew()

	plugins, err := svc.Discover()
	if err != nil {
		return fmt.Errorf("failed to discover plugins: %w", err)
	}

	if len(plugins) == 0 {
		fmt.Fprintln(pluginsStderr, "No weaver plugins found in PATH.")
		return nil
	}

	var rows []pluginRow
	for _, name := range plugins {
		info, err := svc.GetInfo(name)
		if err != nil {
			rows = append(rows, pluginRow{
				Plugin:      name,
				Engine:      "?",
				Version:     "?",
				Extensions:  "?",
				Description: fmt.Sprintf("error: %v", err),
			})
			continue
		}
		rows = append(rows, pluginRow{
			Plugin:      name,
			Engine:      info.Engine,
			Version:     info.Version,
			Extensions:  strings.Join(info.FileExtensions, ", "),
			Description: info.Description,
		})
	}

	switch pluginsOutputFlag {
	case "json":
		return renderJSON(pluginsStdout, rows, pluginsHeryFlag)
	case "yaml":
		return renderYAML(pluginsStdout, rows, pluginsHeryFlag)
	default:
		return renderTable(pluginsStdout, rows, pluginsHeryFlag)
	}
}

func renderJSON(w io.Writer, rows []pluginRow, hery bool) error {
	var data any = rows
	if hery {
		data = heryEnvelope{
			Type: "amadla.org/entity/tools/plugins@v1.0.0",
			Body: rows,
		}
	}
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Fprintln(w, string(out))
	return nil
}

func renderYAML(w io.Writer, rows []pluginRow, hery bool) error {
	var data any = rows
	if hery {
		data = heryEnvelope{
			Type: "amadla.org/entity/tools/plugins@v1.0.0",
			Body: rows,
		}
	}
	out, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	fmt.Fprint(w, string(out))
	return nil
}

func renderTable(w io.Writer, rows []pluginRow, hery bool) error {
	if hery {
		fmt.Fprintln(w, "_type: amadla.org/entity/tools/plugins@v1.0.0")
	}
	table := tablewriter.NewWriter(w)
	table.Header("Plugin", "Engine", "Version", "Extensions", "Description")
	for _, r := range rows {
		table.Append(r.Plugin, r.Engine, r.Version, r.Extensions, r.Description)
	}
	table.Render()
	return nil
}
