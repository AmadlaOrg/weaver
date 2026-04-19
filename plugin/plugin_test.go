package plugin

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Discover(t *testing.T) {
	// Create temp dir with fake plugin binaries
	tmpDir, err := os.MkdirTemp("", "weaver-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create fake weaver-foo and weaver-bar executables
	for _, name := range []string{"weaver-foo", "weaver-bar", "other-tool"} {
		f, err := os.Create(filepath.Join(tmpDir, name))
		require.NoError(t, err)
		require.NoError(t, f.Chmod(0755))
		f.Close()
	}

	// Override PATH
	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string {
		if key == "PATH" {
			return tmpDir
		}
		return ""
	}

	svc := &service{}
	plugins, err := svc.Discover()
	require.NoError(t, err)

	assert.Contains(t, plugins, "weaver-foo")
	assert.Contains(t, plugins, "weaver-bar")
	assert.NotContains(t, plugins, "other-tool")
}

func TestService_DiscoverEmpty(t *testing.T) {
	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string { return "" }

	svc := &service{}
	plugins, err := svc.Discover()
	require.NoError(t, err)
	assert.Nil(t, plugins)
}

func TestService_DiscoverSkipsDirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "weaver-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a directory named weaver-subdir
	require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "weaver-subdir"), 0755))

	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string {
		if key == "PATH" {
			return tmpDir
		}
		return ""
	}

	svc := &service{}
	plugins, err := svc.Discover()
	require.NoError(t, err)
	assert.Empty(t, plugins)
}

func TestService_DiscoverSkipsNonExecutable(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "weaver-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create non-executable file
	f, err := os.Create(filepath.Join(tmpDir, "weaver-noexec"))
	require.NoError(t, err)
	require.NoError(t, f.Chmod(0644))
	f.Close()

	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string {
		if key == "PATH" {
			return tmpDir
		}
		return ""
	}

	svc := &service{}
	plugins, err := svc.Discover()
	require.NoError(t, err)
	assert.Empty(t, plugins)
}

func TestService_DiscoverDeduplicates(t *testing.T) {
	tmpDir1, err := os.MkdirTemp("", "weaver-test1-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir1)

	tmpDir2, err := os.MkdirTemp("", "weaver-test2-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir2)

	// Same plugin name in both dirs
	for _, dir := range []string{tmpDir1, tmpDir2} {
		f, err := os.Create(filepath.Join(dir, "weaver-dup"))
		require.NoError(t, err)
		require.NoError(t, f.Chmod(0755))
		f.Close()
	}

	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string {
		if key == "PATH" {
			return tmpDir1 + string(os.PathListSeparator) + tmpDir2
		}
		return ""
	}

	svc := &service{}
	plugins, err := svc.Discover()
	require.NoError(t, err)
	assert.Len(t, plugins, 1)
}

func TestService_GetInfo(t *testing.T) {
	expectedInfo := Info{
		Name:           "weaver-test",
		Version:        "1.0.0",
		Engine:         "test",
		Description:    "Test plugin",
		FileExtensions: []string{".test", ".tst"},
	}
	infoJSON, _ := json.Marshal(expectedInfo)

	origLookPath := execLookPath
	origCommand := execCommand
	defer func() {
		execLookPath = origLookPath
		execCommand = origCommand
	}()

	execLookPath = func(name string) (string, error) {
		return "/usr/bin/" + name, nil
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo", string(infoJSON))
	}

	svc := &service{}
	info, err := svc.GetInfo("weaver-test")
	require.NoError(t, err)

	assert.Equal(t, "weaver-test", info.Name)
	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, "test", info.Engine)
	assert.Equal(t, []string{".test", ".tst"}, info.FileExtensions)
}

func TestService_GetInfoNotFound(t *testing.T) {
	origLookPath := execLookPath
	defer func() { execLookPath = origLookPath }()

	execLookPath = func(name string) (string, error) {
		return "", exec.ErrNotFound
	}

	svc := &service{}
	_, err := svc.GetInfo("weaver-nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in PATH")
}

func TestService_RenderNotFound(t *testing.T) {
	origLookPath := execLookPath
	defer func() { execLookPath = origLookPath }()

	execLookPath = func(name string) (string, error) {
		return "", exec.ErrNotFound
	}

	svc := &service{}
	err := svc.Render("weaver-nonexistent", "t.tmpl", "", "", &bytes.Buffer{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in PATH")
}

func TestService_FindByExtensionNotFound(t *testing.T) {
	origGetenv := osGetenv
	defer func() { osGetenv = origGetenv }()
	osGetenv = func(key string) string { return "" }

	svc := &service{}
	_, err := svc.FindByExtension(".xyz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no plugin found")
}
