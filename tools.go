package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func loadFromJSON(path string, pointer_obj any) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("\n%s not found, make new", path)
		return
	}

	err = json.Unmarshal(data, &pointer_obj)
	if err != nil {
		fmt.Printf("\nLoading %s error: %s", path, err)
	}
}

func writeToJSON(path string, obj any) {
	os.MkdirAll(filepath.Dir(path), 0755)

	data, err := json.MarshalIndent(obj, "", " ")

	if err != nil {
		fmt.Printf("\nSaving %s apps error: %s", path, err.Error())
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Printf("\nWriting %s apps error: %s", path, err.Error())
	}
}
