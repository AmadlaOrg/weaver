package cmd

import (
	"errors"
	"io"
	"testing"

	"github.com/AmadlaOrg/weaver/weave"
	"github.com/stretchr/testify/assert"
)

type mockWeaver struct {
	err error
}

func (m *mockWeaver) Do() error {
	return m.err
}

func resetWeaveFlags() {
	templatePath = ""
	outputPath = ""
	entityPath = ""
}

func TestWeaveCmd_FlagsRegisteredBeforeRun(t *testing.T) {
	for _, name := range []string{"template", "output", "entity"} {
		assert.NotNilf(t, WeaveCmd.Flags().Lookup(name),
			"flag %q must be registered at init time, not inside Run", name)
	}
}

func TestWeaveCmd_Execute(t *testing.T) {
	origWeaveNew := weaveNew
	defer func() { weaveNew = origWeaveNew }()
	defer resetWeaveFlags()

	var capturedTmpl string
	weaveNew = func(tmpl string, input io.Reader, output io.Writer) weave.Weaver {
		capturedTmpl = tmpl
		return &mockWeaver{}
	}

	WeaveCmd.SetArgs([]string{"-t", "config.tmpl"})
	err := WeaveCmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "config.tmpl", capturedTmpl)
}

func TestWeaveCmd_ExecuteWeaveError(t *testing.T) {
	origWeaveNew := weaveNew
	defer func() { weaveNew = origWeaveNew }()
	defer resetWeaveFlags()

	weaveNew = func(tmpl string, input io.Reader, output io.Writer) weave.Weaver {
		return &mockWeaver{err: errors.New("template parse failed")}
	}

	WeaveCmd.SetArgs([]string{"-t", "config.tmpl"})
	err := WeaveCmd.Execute()
	assert.Error(t, err, "weave failures must propagate as errors, not exit 0")
	assert.Contains(t, err.Error(), "template parse failed")
}

func TestWeaveCmd_ExecuteMissingEntityFile(t *testing.T) {
	origWeaveNew := weaveNew
	defer func() { weaveNew = origWeaveNew }()
	defer resetWeaveFlags()

	weaveNew = func(tmpl string, input io.Reader, output io.Writer) weave.Weaver {
		t.Fatal("weaveNew must not be called when the entity file cannot be opened")
		return nil
	}

	WeaveCmd.SetArgs([]string{"-t", "config.tmpl", "-e", "/nonexistent/entity.yaml"})
	err := WeaveCmd.Execute()
	assert.Error(t, err)
}
