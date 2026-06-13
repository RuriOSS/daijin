// Copyright 2024 moe-hacker
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const version = "2.0.0-dev"

var (
	themeColor = lipgloss.Color("#FEE4D0")

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(themeColor).
		MarginBottom(1)

	menuStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(themeColor).
		Padding(1, 2).
		Width(60)
)

type model struct {
	cursor int
	choices []string
}

func initialModel() model {
	return model{
		cursor: 0,
		choices: []string{
			"Install",
			"Run",
			"Remove",
			"Register",
			"Clean rootfs",
			"Exit",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.cursor == len(m.choices)-1 {
				return m, tea.Quit
			}
			// TODO: Handle menu actions
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	title := titleStyle.Render(fmt.Sprintf("Daijin v%s", version))

	menu := ""
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = "▸"
			choice = lipgloss.NewStyle().
				Foreground(themeColor).
				Bold(true).
				Render(choice)
		}
		menu += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\n↑/k up • ↓/j down • enter select • q quit")

	content := title + "\n\n" + menu + help

	return "\n" + menuStyle.Render(content) + "\n"
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("Daijin v%s\n", version)
			return
		case "--help", "-h":
			fmt.Println("Daijin - Lightweight container manager for Termux")
			fmt.Println("\nUsage:")
			fmt.Println("  daijin           Start interactive TUI")
			fmt.Println("  daijin --version Show version")
			fmt.Println("  daijin --help    Show this help")
			return
		}
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
