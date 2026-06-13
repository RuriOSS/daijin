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
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moe-hacker/daijin/internal/container"
	"github.com/moe-hacker/daijin/internal/system"
)

type registerDialog struct {
	step        int
	nameInput   textinput.Model
	pathInput   textinput.Model
	backend     string
	err         error
	manager     *container.Manager
	parent      tea.Model
}

func NewRegisterDialog(mgr *container.Manager, parent tea.Model) *registerDialog {
	nameInput := textinput.New()
	nameInput.Placeholder = "container-name"
	nameInput.Focus()

	pathInput := textinput.New()
	pathInput.Placeholder = "/path/to/container"

	return &registerDialog{
		step:      0,
		nameInput: nameInput,
		pathInput: pathInput,
		manager:   mgr,
		parent:    parent,
	}
}

func (d *registerDialog) Init() tea.Cmd {
	return textinput.Blink
}

func (d *registerDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return d.parent, nil

		case "enter":
			switch d.step {
			case 0: // Name input
				if d.nameInput.Value() != "" {
					d.step = 1
					d.pathInput.Focus()
					return d, textinput.Blink
				}
			case 1: // Path input
				if d.pathInput.Value() != "" {
					d.step = 2
					return d, nil
				}
			case 2: // Backend selection
				// Already handled in number keys
			case 3: // Done
				return d.parent, nil
			}

		case "1":
			if d.step == 2 {
				d.backend = "proot"
				return d, d.performRegister()
			}

		case "2":
			if d.step == 2 && system.HasRootAccess() {
				d.backend = "ruri"
				return d, d.performRegister()
			}
		}
	}

	// Update active input
	switch d.step {
	case 0:
		d.nameInput, cmd = d.nameInput.Update(msg)
	case 1:
		d.pathInput, cmd = d.pathInput.Update(msg)
	}

	return d, cmd
}

func (d *registerDialog) performRegister() tea.Cmd {
	return func() tea.Msg {
		name := d.nameInput.Value()
		path := d.pathInput.Value()

		// Get absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			d.err = fmt.Errorf("invalid path: %w", err)
			return nil
		}

		// Register container
		if err := d.manager.Register(name, absPath, d.backend); err != nil {
			d.err = err
			return nil
		}

		d.step = 3
		return nil
	}
}

func (d *registerDialog) View() string {
	title := TitleStyle.Render(fmt.Sprintf("Register Container (%d/3)", d.step+1))

	var content string

	switch d.step {
	case 0:
		content = title + "\n\n" +
			SubtitleStyle.Render("Container name:") + "\n" +
			d.nameInput.View() + "\n\n" +
			HelpStyle.Render("enter next • esc cancel")

	case 1:
		content = title + "\n\n" +
			SubtitleStyle.Render("Container directory path:") + "\n" +
			d.pathInput.View() + "\n\n" +
			HelpStyle.Render("enter next • esc cancel")

	case 2:
		backendOpts := "  [1] proot (no root required)\n"
		if system.HasRootAccess() {
			backendOpts += "  [2] ruri (requires root)\n"
		} else {
			backendOpts += "  [2] ruri ✗ Root not available\n"
		}

		content = title + "\n\n" +
			SubtitleStyle.Render("Choose backend:") + "\n\n" +
			backendOpts + "\n" +
			HelpStyle.Render("1/2 select • esc cancel")

	case 3:
		if d.err != nil {
			content = title + "\n\n" +
				ErrorStyle.Render(fmt.Sprintf("Error: %v", d.err)) + "\n\n" +
				HelpStyle.Render("enter continue")
		} else {
			content = title + "\n\n" +
				SuccessStyle.Render(fmt.Sprintf("✓ Container '%s' registered successfully!", d.nameInput.Value())) + "\n\n" +
				HelpStyle.Render("enter continue")
		}
	}

	return BorderStyle.Render(content)
}
