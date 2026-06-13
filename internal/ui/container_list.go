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

package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moe-hacker/daijin/internal/container"
	"github.com/moe-hacker/daijin/internal/rootfs"
	"github.com/moe-hacker/daijin/internal/system"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Delete   key.Binding
	Install  key.Binding
	Refresh  key.Binding
	Register key.Binding
	Help     key.Binding
	Quit     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "start"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Register: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "register"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type containerListModel struct {
	manager    *container.Manager
	containers []*container.Config
	cursor     int
	width      int
	height     int
	err        error
	message    string
}

func NewContainerListModel() containerListModel {
	manager := container.NewManager(system.ContainersDir)
	containers, err := manager.List()

	model := containerListModel{
		manager:    manager,
		containers: containers,
		cursor:     0,
	}

	if err != nil {
		model.err = err
	}

	return model
}

func (m containerListModel) Init() tea.Cmd {
	return nil
}

func (m containerListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.containers)-1 {
				m.cursor++
			}

		case key.Matches(msg, keys.Enter):
			if len(m.containers) > 0 {
				return m, m.startContainer()
			}

		case key.Matches(msg, keys.Delete):
			if len(m.containers) > 0 {
				c := m.containers[m.cursor]
				dialog := NewConfirmDialog(
					"Delete Container",
					fmt.Sprintf("Are you sure you want to delete '%s'?\nThis will remove the container directory: %s", c.Name, c.ContainerDir),
					func() tea.Msg {
						// Perform deletion
						if err := m.manager.Remove(c.Name); err != nil {
							m.err = fmt.Errorf("failed to delete: %w", err)
						} else {
							m.message = fmt.Sprintf("Deleted container: %s", c.Name)
							// Reload list
							m.containers, m.err = m.manager.List()
							if m.cursor >= len(m.containers) && m.cursor > 0 {
								m.cursor--
							}
						}
						return nil
					},
					func() tea.Msg {
						m.message = "Deletion cancelled"
						return nil
					},
					m,
				)
				return dialog, nil
			}

		case key.Matches(msg, keys.Refresh):
			return m, m.refresh()

		case key.Matches(msg, keys.Register):
			dialog := NewRegisterDialog(m.manager, m)
			return dialog, dialog.Init()

		case key.Matches(msg, keys.Install):
			// Check if rurima is available
			if client := rootfs.NewRurimaClient(); client == nil {
				m.message = "Error: rurima not found. Please install rurima first."
			} else {
				// Switch to install wizard
				wizard := NewInstallWizard(m.manager)
				return wizard, wizard.Init()
			}
		}
	}

	return m, nil
}

func (m containerListModel) View() string {
	var b strings.Builder

	// Title
	title := TitleStyle.Render("Daijin v2.0")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Error message
	if m.err != nil {
		b.WriteString(ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	// Info message
	if m.message != "" {
		b.WriteString(InfoStyle.Render(m.message))
		b.WriteString("\n\n")
	}

	// Container list
	if len(m.containers) == 0 {
		b.WriteString(SubtitleStyle.Render("No containers found."))
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("Press 'i' to install a new container."))
	} else {
		subtitle := SubtitleStyle.Render(fmt.Sprintf("Containers (%d)", len(m.containers)))
		b.WriteString(subtitle)
		b.WriteString("\n\n")

		for i, c := range m.containers {
			cursor := "  "
			style := NormalStyle

			if i == m.cursor {
				cursor = CursorStyle.Render("▸ ")
				style = SelectedStyle
			}

			// Container line: name, backend, path
			line := fmt.Sprintf("%-30s %-8s %s",
				c.Name,
				c.Backend,
				m.truncatePath(c.ContainerDir, 40),
			)

			b.WriteString(cursor)
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}

		// Container details (if selected)
		if len(m.containers) > 0 && m.cursor < len(m.containers) {
			b.WriteString("\n")
			b.WriteString(m.renderDetails(m.containers[m.cursor]))
		}
	}

	// Help
	b.WriteString("\n")
	b.WriteString(m.helpView())

	return BorderStyle.Render(b.String())
}

func (m containerListModel) renderDetails(c *container.Config) string {
	var b strings.Builder

	b.WriteString(SubtitleStyle.Render("Details:"))
	b.WriteString("\n")

	details := []struct {
		key   string
		value string
	}{
		{"Backend", c.Backend},
		{"Path", c.ContainerDir},
		{"Config", c.ConfigPath},
	}

	if c.ExtraArgs != "" {
		details = append(details, struct {
			key   string
			value string
		}{"Extra Args", c.ExtraArgs})
	}

	for _, d := range details {
		b.WriteString(KeyStyle.Render(d.key + ": "))
		b.WriteString(ValueStyle.Render(d.value))
		b.WriteString("\n")
	}

	return b.String()
}

func (m containerListModel) helpView() string {
	helpKeys := []string{
		"↑/k up",
		"↓/j down",
		"enter start",
		"d delete",
		"i install",
		"R register",
		"r refresh",
		"q quit",
	}

	return HelpStyle.Render(strings.Join(helpKeys, " • "))
}

func (m containerListModel) truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}

	// Show "...tail"
	return "..." + path[len(path)-maxLen+3:]
}

func (m containerListModel) startContainer() tea.Cmd {
	return func() tea.Msg {
		if m.cursor >= len(m.containers) {
			return nil
		}

		c := m.containers[m.cursor]

		// Exit TUI before starting container
		tea.Quit()

		// Start the container
		if err := m.manager.Start(c.Name, nil); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start container: %v\n", err)
			os.Exit(1)
		}

		return nil
	}
}

func (m containerListModel) refresh() tea.Cmd {
	return func() tea.Msg {
		containers, err := m.manager.List()
		m.containers = containers
		if err != nil {
			m.err = err
		} else {
			m.message = "Refreshed container list"
		}
		return nil
	}
}
