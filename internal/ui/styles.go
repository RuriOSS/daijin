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

import "github.com/charmbracelet/lipgloss"

// Theme colors
var (
	ThemeColor   = lipgloss.Color("#FEE4D0")
	ErrorColor   = lipgloss.Color("#FF6B6B")
	SuccessColor = lipgloss.Color("#51CF66")
	InfoColor    = lipgloss.Color("#4DABF7")
	MutedColor   = lipgloss.Color("#868E96")
)

// Styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ThemeColor).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginBottom(1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ThemeColor).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			PaddingLeft(1).
			PaddingRight(1)

	CursorStyle = lipgloss.NewStyle().
			Foreground(ThemeColor).
			Bold(true)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ThemeColor).
			Padding(1, 2)

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginTop(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor)

	KeyStyle = lipgloss.NewStyle().
			Foreground(ThemeColor).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))
)
