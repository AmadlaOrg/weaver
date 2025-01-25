package fs

/*package main

import (
"encoding/json"
"os"
"text/template"
)

func main() {
	// Open the template file
	tmplFile := "generic.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	// Open the JSON data file
	dataFile, err := os.Open("data.json")
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	// Open the output file for writing
	outputFile, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Create a JSON decoder for streaming
	decoder := json.NewDecoder(dataFile)

	// Begin decoding an array
	if _, err := decoder.Token(); err != nil {
		panic(err)
	}

	// Iterate over each object in the JSON array
	for decoder.More() {
		var item map[string]interface{}
		if err := decoder.Decode(&item); err != nil {
			panic(err)
		}

		// Execute the template for each item and write to the output file
		err := tmpl.Execute(outputFile, item)
		if err != nil {
			panic(err)
		}

		// Optionally add a separator (e.g., a newline) between entries
		outputFile.WriteString("\n")
	}
}*/
