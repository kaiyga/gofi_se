package main

import (
	"os"
	"path/filepath"
)

var pinnedFile = filepath.Join(os.Getenv("HOME"), ".config", "gofi-launcher", "pinned.json")

func savePinnedApps() {
	writeToJSON(pinnedFile, pinnedApps)
}

func loadPinnedApps() {
	loadFromJSON(pinnedFile, &pinnedApps)
}
