package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"os/exec"

	"github.com/sahilm/fuzzy"
	"golang.org/x/term"
)

type model struct {
	list       list.Model
	apps       []app
	filtered   []app
	searchText string
	searching  bool
}

var pinnedApps = make(map[string]bool)

func filterApps(input string, apps []app) []app {

	if input == "" {
		return apps
	}

	var appNames []string
	for _, a := range apps {
		appNames = append(appNames, a.name)
	}

	matches := fuzzy.Find(input, appNames)

	var filtered []app
	for _, match := range matches {
		filtered = append(filtered, apps[match.Index])
	}
	return filtered
}

func (m model) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pin/unpin"),
		),
	}
}

func (m model) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

func (m *model) updateList() {
	var pinned, regular []app

	for _, a := range m.filtered {
		displayName := a.name
		if pinnedApps[a.name] {
			displayName = "⭐ " + a.name
		}

		newApp := app{
			name:        displayName,
			description: a.description,
			cmd:         a.cmd,
			weight:      a.weight,
			file:        a.file,
		}

		if pinnedApps[a.name] {
			pinned = append(pinned, newApp)
		} else {
			regular = append(regular, newApp)
		}
	}
	m.list.SetItems(convertToListItems(append(pinned, regular...)))
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		// for russian
		// if len(msg.Runes) > 0 {
		// 	switch msg.Runes[0] {
		// 	case 1088:
		// 		key = "h"
		// 	case 1086:
		// 		key = "j"
		// 	case 1083:
		// 		key = "k"
		// 	case 1076:
		// 		key = "l"
		// 	}
		// }

		switch key {
		case "q":
			return m, tea.Quit

		case "/":
			m.searching = true
			m.searchText = ""

		case "esc":
			m.searching = false
			m.searchText = ""

		case "p":
			selected, ok := m.list.SelectedItem().(app)
			if ok {
				originalName := strings.TrimPrefix(selected.name, "⭐ ")
				if pinnedApps[originalName] {
					delete(pinnedApps, originalName)
				} else {
					pinnedApps[originalName] = true
				}
				savePinnedApps()
				m.updateList()
			}

		case "enter":
			selected, ok := m.list.SelectedItem().(app)
			if ok {
				appsWeights[selected.File()] = selected.Weight() + 1

				m.list.Title = fmt.Sprint(selected)
				saveAppsWeight()

				cmd := exec.Command("setsid", selected.cmd)
				cmd.Start()

				return m, tea.Quit
			}

		case "backspace":
			if len(m.searchText) > 0 {
				m.searchText = m.searchText[:len(m.searchText)-1]
			}

		default:
			if m.searching && len(key) == 1 {
				m.searchText += key
			}
		}
		if m.searching {
			m.filtered = filterApps(m.searchText, m.apps)
			m.list.SetItems(convertToListItems(m.filtered))
		}
	}
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	// style := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	// return fmt.Sprintf("%s\nSearch: %s\n\n%s", style.Render("TUI Launcher"), m.searchText, m.list.View())
	return m.list.View()
}

func convertToListItems(apps []app) []list.Item {
	items := make([]list.Item, len(apps))
	for i, a := range apps {
		items[i] = a
	}
	return items
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getTerminalSize() (int, int, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	return width, height, err
}

func main() {
	clearScreen()
	loadPinnedApps()
	allApps := loadApplications()

	width, height, err := getTerminalSize()
	if err != nil {
		width, height = 80, 24
	}

	l := list.New(convertToListItems(allApps), list.NewDefaultDelegate(), width-5, height)
	l.Title = "Apps"
	l.SetShowHelp(true)
	// l.Styles.Title = titleStyle

	m := model{list: l, apps: allApps, filtered: allApps}
	m.list.SetItems(convertToListItems(m.filtered))
	m.updateList()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
	}
}
