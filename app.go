package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var desktopDirs = []string{
	"~/.nix-profile/share/applications",
	"/run/current-system/sw/share/applications",
}

func loadApplications() []app {
	var apps []app

	for _, dir := range desktopDirs {
		home, _ := os.UserHomeDir()
		expandedDir := strings.Replace(dir, "~", home, 1)
		files, err := filepath.Glob(filepath.Join(expandedDir, "*.desktop"))
		if err != nil {
			fmt.Println("Error while reading dir:", err)
			continue
		}

		for _, file := range files {
			parsedApp, err := parseDesktopFile(file)
			if err == nil && parsedApp.name != "" {
				apps = append(apps, parsedApp)
			}
		}

	}

	return apps
}

func parseDesktopFile(path string) (app, error) {
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

	return app{name: name, description: description, cmd: execCmd}, nil
}
