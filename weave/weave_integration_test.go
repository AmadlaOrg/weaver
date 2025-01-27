package weave

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const (
	testSimpleFixturePath = "../test/fixture/simple"
)

func TestSWeave_DoWithJSONStringToFile(t *testing.T) {
	input := `[
		{"Key": "Key1", "Value": "Value1"},
		{"Key": "Key2", "Value": "Value2"},
		{"Key": "Key3", "Value": "Value3"}
	]`

	tmplPath, err := filepath.Abs(filepath.Join(testSimpleFixturePath, "generic.tmpl"))
	if err != nil {
		t.Fatalf("Failed to get template path: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "weaver_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create a temporary file for output
	tempFile, err := os.CreateTemp(tempDir, "output*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}(tempFile.Name())

	// Call Do method
	err = NewWeaveService(tmplPath, bytes.NewReader([]byte(input)), tempFile).Do()
	if err != nil {
		t.Fatalf("Do method failed: %v", err)
	}
	err = tempFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Validate output file content
	output, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "Key: Key1, Value: Value1\nKey: Key2, Value: Value2\nKey: Key3, Value: Value3\n"
	if string(output) != expected {
		t.Errorf("Output mismatch\nExpected:\n%s\nGot:\n%s", expected, string(output))
	}
}

func TestSWeave_DoWithJSONStringToStdout(t *testing.T) {
	input := `[
		{"Key": "Key1", "Value": "Value1"},
		{"Key": "Key2", "Value": "Value2"},
		{"Key": "Key3", "Value": "Value3"}
	]`

	tmplPath, err := filepath.Abs(filepath.Join(testSimpleFixturePath, "generic.tmpl"))
	if err != nil {
		t.Fatalf("Failed to get template path: %v", err)
	}

	// Use a buffer to capture output
	var stdout bytes.Buffer

	// Call Do method
	err = NewWeaveService(tmplPath, bytes.NewReader([]byte(input)), &stdout).Do()
	if err != nil {
		t.Fatalf("Do method failed: %v", err)
	}

	// Validate stdout content
	expected := "Key: Key1, Value: Value1\nKey: Key2, Value: Value2\nKey: Key3, Value: Value3\n"
	if stdout.String() != expected {
		t.Errorf("Stdout mismatch\nExpected:\n%s\nGot:\n%s", expected, stdout.String())
	}
}

func TestSWeave_DoWithYAMLStringToStdout(t *testing.T) {
	input := `
- Key: Key1
  Value: Value1
- Key: Key2
  Value: Value2
- Key: Key3
  Value: Value3
`

	tmplPath, err := filepath.Abs(filepath.Join(testSimpleFixturePath, "generic.tmpl"))
	if err != nil {
		t.Fatalf("Failed to get template path: %v", err)
	}

	// Use a buffer to capture output
	var stdout bytes.Buffer

	// Call Do method
	err = NewWeaveService(tmplPath, bytes.NewReader([]byte(input)), &stdout).Do()
	if err != nil {
		t.Fatalf("Do method failed: %v", err)
	}

	// Validate stdout content
	expected := "Key: Key1, Value: Value1\nKey: Key2, Value: Value2\nKey: Key3, Value: Value3\n"
	if stdout.String() != expected {
		t.Errorf("Stdout mismatch\nExpected:\n%s\nGot:\n%s", expected, stdout.String())
	}
}

func TestSWeave_DoWithJSONFileToFile(t *testing.T) {
	// Create a temp JSON file
	input := `[
		{"Key": "Key1", "Value": "Value1"},
		{"Key": "Key2", "Value": "Value2"},
		{"Key": "Key3", "Value": "Value3"}
	]`

	tempDir, err := os.MkdirTemp("", "weaver_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	inputFile, err := os.CreateTemp(tempDir, "input*.json")
	if err != nil {
		t.Fatalf("Failed to create temp input file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}(inputFile.Name())

	if _, err := inputFile.WriteString(input); err != nil {
		t.Fatalf("Failed to write to temp input file: %v", err)
	}
	err = inputFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp input file: %v", err)
	}

	tmplPath, err := filepath.Abs(filepath.Join(testSimpleFixturePath, "generic.tmpl"))
	if err != nil {
		t.Fatalf("Failed to get template path: %v", err)
	}

	// Create a temp file for output
	tempFile, err := os.CreateTemp("", "output*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}(tempFile.Name())

	// Open the input file
	inputHandle, err := os.Open(inputFile.Name())
	if err != nil {
		t.Fatalf("Failed to open input file: %v", err)
	}
	defer func(inputHandle *os.File) {
		err := inputHandle.Close()
		if err != nil {
			t.Fatalf("Failed to close input file: %v", err)
		}
	}(inputHandle)

	// Call Do method
	err = NewWeaveService(tmplPath, inputHandle, tempFile).Do()
	if err != nil {
		t.Fatalf("Do method failed: %v", err)
	}
	err = tempFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Validate output file content
	output, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "Key: Key1, Value: Value1\nKey: Key2, Value: Value2\nKey: Key3, Value: Value3\n"
	if string(output) != expected {
		t.Errorf("Output mismatch\nExpected:\n%s\nGot:\n%s", expected, string(output))
	}
}

// Repeat similar tests for YAML files.
