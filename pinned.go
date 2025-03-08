package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var pinnedFile = filepath.Join(os.Getenv("HOME"), ".config", "tui-launcher", "pinned.json")

func savePinnedApps() {
	os.MkdirAll(filepath.Dir(pinnedFile), 0755)

	data, err := json.MarshalIndent(pinnedApps, "", " ")

	if err != nil {
		fmt.Println("Saving pinned apps error:", err)
	}

	err = os.WriteFile(pinnedFile, data, 0644)
	if err != nil {
		fmt.Println("Writing pinned apps error:", err)
	}
}

func loadPinnedApps() {
	data, err := os.ReadFile(pinnedFile)
	if err != nil {
		fmt.Println("pinned.json not found, make new")
		return
	}

	err = json.Unmarshal(data, &pinnedApps)
	if err != nil {
		fmt.Println("Loading pinned.json error:", err)
	}
}
