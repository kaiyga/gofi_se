package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var desktopDirs = []string{
	// Default linux .desktop files
	"/usr/share/applications/",
	"~/.local/share/applications/",
	// NixOS .desktop files
	"~/.nix-profile/share/applications",
	"/run/current-system/sw/share/applications",
}

// Object of application in launcher list
type app struct {
	name        string
	cmd         string
	description string
	// .desktop file
	file   string
	weight int
}

// Get count how many times user runned this app
func (a app) Weight() int { return a.weight }

// Get application .desktop file
func (a app) File() string { return a.file }

func (a app) Title() string       { return a.name }
func (a app) Description() string { return a.description }
func (a app) FilterValue() string { return a.name }

var appsWeights = make(map[string]int)
var appsWeightFile string = filepath.Join(os.Getenv("HOME"), ".config", "gofi-launcher", "gofi-drun.json")

func saveAppsWeight() { writeToJSON(appsWeightFile, appsWeights) }
func loadAppsWeight() { loadFromJSON(appsWeightFile, &appsWeights) }

func loadApplications() []app {
	var apps []app
	loadAppsWeight()

	for file, weight := range appsWeights {
		parsedApp, err := parseDesktopFile(file, weight)
		if err == nil && parsedApp.name != "" {
			apps = append(apps, parsedApp)
		}
	}

	for _, dir := range desktopDirs {
		home, _ := os.UserHomeDir()
		expandedDir := strings.Replace(dir, "~", home, 1)
		files, err := filepath.Glob(filepath.Join(expandedDir, "*.desktop"))
		if err != nil {
			fmt.Println("Error while reading dir:", err)
			continue
		}

		for _, file := range files {
			parsedApp, err := parseDesktopFile(file, 0)
			_, ok := appsWeights[parsedApp.File()]
			if err == nil && parsedApp.name != "" && !ok {
				apps = append(apps, parsedApp)
			}
		}

	}

	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Weight() > apps[j].Weight()
	})

	return apps
}

func parseDesktopFile(path string, weight int) (app, error) {
	file, err := os.Open(path)
	if err != nil {
		return app{}, err
	}
	defer file.Close()

	var name, execCmd, description string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Name=") {
			name = strings.TrimPrefix(line, "Name=")
		} else if strings.HasPrefix(line, "Exec=") {
			execCmd = strings.Trim(execCmd, `"`)
			execCmd = strings.TrimPrefix(line, "Exec=")
			execCmd = strings.Fields(execCmd)[0]
			execCmd = strings.ReplaceAll(execCmd, "%U", "")
			execCmd = strings.ReplaceAll(execCmd, "%F", "")
		} else if strings.HasPrefix(line, "Comment=") {
			description = strings.TrimPrefix(line, "Comment=")
		}

		if name != "" && execCmd != "" {
			break
		}
	}
	if description == "" {
		description = execCmd
	}

	if err := scanner.Err(); err != nil {
		return app{}, err
	}

	return app{name: name, description: description, cmd: execCmd, weight: weight, file: path}, nil
}
