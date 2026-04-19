package cmd

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/AmadlaOrg/weaver/plugin"
	"github.com/stretchr/testify/assert"
)

type mockPluginService struct {
	discoverFunc      func() ([]string, error)
	getInfoFunc       func(string) (*plugin.Info, error)
	renderFunc        func(string, string, string, string, io.Reader) error
	findByExtFunc     func(string) (string, error)
}

func (m *mockPluginService) Discover() ([]string, error) {
	return m.discoverFunc()
}
func (m *mockPluginService) GetInfo(name string) (*plugin.Info, error) {
	return m.getInfoFunc(name)
}
func (m *mockPluginService) Render(name, tmpl, file, output string, stdin io.Reader) error {
	return m.renderFunc(name, tmpl, file, output, stdin)
}
func (m *mockPluginService) FindByExtension(ext string) (string, error) {
	return m.findByExtFunc(ext)
}

func TestRunRender_WithEngine(t *testing.T) {
	var capturedPlugin, capturedTmpl, capturedFile, capturedOutput string

	origFactory := pluginNew
	defer func() { pluginNew = origFactory }()

	pluginNew = func() plugin.Service {
		return &mockPluginService{
			renderFunc: func(name, tmpl, file, output string, _ io.Reader) error {
				capturedPlugin = name
				capturedTmpl = tmpl
				capturedFile = file
				capturedOutput = output
				return nil
			},
		}
	}

	// Reset flags
	renderEngine = "jinja2"
	renderTemplatePath = "config.j2"
	renderFilePath = "data.yaml"
	renderOutputPath = "out.conf"

	err := runRender(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "weaver-jinja2", capturedPlugin)
	assert.Equal(t, "config.j2", capturedTmpl)
	assert.Equal(t, "data.yaml", capturedFile)
	assert.Equal(t, "out.conf", capturedOutput)
}

func TestRunRender_AutoDetect(t *testing.T) {
	var capturedPlugin string

	origFactory := pluginNew
	defer func() { pluginNew = origFactory }()

	pluginNew = func() plugin.Service {
		return &mockPluginService{
			findByExtFunc: func(ext string) (string, error) {
				return "weaver-mustache", nil
			},
			renderFunc: func(name, tmpl, file, output string, _ io.Reader) error {
				capturedPlugin = name
				return nil
			},
		}
	}

	renderEngine = ""
	renderTemplatePath = "config.mustache"
	renderFilePath = ""
	renderOutputPath = ""

	err := runRender(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "weaver-mustache", capturedPlugin)
}

func TestRunRender_NoExtension(t *testing.T) {
	origFactory := pluginNew
	defer func() { pluginNew = origFactory }()

	pluginNew = func() plugin.Service {
		return &mockPluginService{}
	}

	renderEngine = ""
	renderTemplatePath = "Dockerfile"
	renderFilePath = ""
	renderOutputPath = ""

	err := runRender(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no file extension")
}

func TestRunRender_UnknownExtension(t *testing.T) {
	origFactory := pluginNew
	defer func() { pluginNew = origFactory }()

	pluginNew = func() plugin.Service {
		return &mockPluginService{
			findByExtFunc: func(ext string) (string, error) {
				return "", fmt.Errorf("no plugin found for extension %q", ext)
			},
		}
	}

	renderEngine = ""
	renderTemplatePath = "config.xyz"
	renderFilePath = ""
	renderOutputPath = ""

	err := runRender(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no plugin found")
}

func TestRunPlugins_Empty(t *testing.T) {
	origFactory := pluginsNew
	defer func() { pluginsNew = origFactory }()

	pluginsNew = func() plugin.Service {
		return &mockPluginService{
			discoverFunc: func() ([]string, error) {
				return nil, nil
			},
		}
	}

	// Capture stderr
	var buf bytes.Buffer
	oldStderr := pluginsStderr
	pluginsStderr = &buf
	defer func() { pluginsStderr = oldStderr }()

	err := runPlugins(nil, nil)
	assert.NoError(t, err)
}
