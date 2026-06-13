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

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type confirmDialog struct {
	title   string
	message string
	onYes   func() tea.Msg
	onNo    func() tea.Msg
	parent  tea.Model
}

func NewConfirmDialog(title, message string, onYes, onNo func() tea.Msg, parent tea.Model) *confirmDialog {
	return &confirmDialog{
		title:   title,
		message: message,
		onYes:   onYes,
		onNo:    onNo,
		parent:  parent,
	}
}

func (d *confirmDialog) Init() tea.Cmd {
	return nil
}

func (d *confirmDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("y", "Y"))):
			if d.onYes != nil {
				d.onYes()
			}
			return d.parent, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("n", "N", "esc"))):
			if d.onNo != nil {
				d.onNo()
			}
			return d.parent, nil
		}
	}

	return d, nil
}

func (d *confirmDialog) View() string {
	content := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
		TitleStyle.Render(d.title),
		d.message,
		SuccessStyle.Render("[Y]")+" Yes   "+ErrorStyle.Render("[N]")+" No",
		HelpStyle.Render("y yes • n/esc no"),
	)

	return BorderStyle.Render(content)
}
