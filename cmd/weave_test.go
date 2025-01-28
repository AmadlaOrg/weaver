package cmd

import (
	"github.com/AmadlaOrg/weaver/weave"
	"github.com/spf13/cobra"
	"io"
	"os"
	"testing"
)

func TestRunWeave(t *testing.T) {
	tests := []struct {
		name                         string
		inputCmd                     *cobra.Command
		inputArgs                    []string
		internalFileIsFile           func(string) (bool, error)
		internalOsOpen               func(string) (*os.File, error)
		internalOsCreate             func(string) (*os.File, error)
		internalWeaveNewWeaveService func(string, io.Reader, io.Writer) weave.IWeave
		expectedCmd                  *cobra.Command
	}{
		{
			name:      "Test",
			inputCmd:  &cobra.Command{},
			inputArgs: []string{},
			internalFileIsFile: func(string) (bool, error) {
				return true, nil
			},
			internalOsOpen: func(string) (*os.File, error) {
				return nil, nil
			},
			internalOsCreate: func(string) (*os.File, error) {
				return nil, nil
			},
			internalWeaveNewWeaveService: func(string, io.Reader, io.Writer) weave.IWeave {
				return nil
			},
			expectedCmd: &cobra.Command{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runWeave(tt.inputCmd, tt.inputArgs)
			print(tt.inputCmd.Flags())
		})
	}
}
