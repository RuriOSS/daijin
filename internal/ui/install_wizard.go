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
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moe-hacker/daijin/internal/container"
	"github.com/moe-hacker/daijin/internal/rootfs"
	"github.com/moe-hacker/daijin/internal/system"
)

type installStep int

const (
	stepSource installStep = iota
	stepBackend
	stepSearch
	stepTag
	stepInstall
	stepDone
)

type installWizard struct {
	step         installStep
	source       string // "docker" or "lxc"
	backend      string // "proot" or "ruri"
	selectedImg  int
	selectedTag  int
	searchInput  textinput.Model
	images       []interface{} // []DockerImage or []LXCDistro
	tags         []interface{} // []DockerTag or []LXCVersion
	progress     progress.Model
	progressVal  float64
	err          error
	message      string
	manager      *container.Manager
	rurima       *rootfs.RurimaClient
}

type progressMsg int

type installDoneMsg struct {
	name string
	err  error
}

func NewInstallWizard(mgr *container.Manager) *installWizard {
	ti := textinput.New()
	ti.Placeholder = "Search for image..."
	ti.Focus()
	ti.CharLimit = 50

	prog := progress.New(progress.WithDefaultGradient())

	return &installWizard{
		step:        stepSource,
		searchInput: ti,
		progress:    prog,
		manager:     mgr,
		rurima:      rootfs.NewRurimaClient(),
	}
}

func (w *installWizard) Init() tea.Cmd {
	return textinput.Blink
}

func (w *installWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return w.handleKey(msg)

	case progressMsg:
		w.progressVal = float64(msg) / 100.0
		if w.progressVal >= 1.0 {
			w.step = stepDone
			w.message = "Installation complete!"
		}
		return w, nil

	case installDoneMsg:
		if msg.err != nil {
			w.err = msg.err
		} else {
			w.message = fmt.Sprintf("Container '%s' installed successfully!", msg.name)
			w.step = stepDone
		}
		return w, nil
	}

	return w, nil
}

func (w *installWizard) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch w.step {
	case stepSource:
		return w.handleSourceStep(msg)
	case stepBackend:
		return w.handleBackendStep(msg)
	case stepSearch:
		return w.handleSearchStep(msg)
	case stepTag:
		return w.handleTagStep(msg)
	case stepDone:
		if key.Matches(msg, keys.Enter) || key.Matches(msg, keys.Quit) {
			// Return to main list
			return NewContainerListModel(), nil
		}
	}

	if key.Matches(msg, keys.Quit) {
		return NewContainerListModel(), nil
	}

	return w, nil
}

func (w *installWizard) handleSourceStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		w.source = "docker"
		w.step = stepBackend
	case "2":
		w.source = "lxc"
		w.step = stepBackend
	case "esc":
		return NewContainerListModel(), nil
	}
	return w, nil
}

func (w *installWizard) handleBackendStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		w.backend = "proot"
		w.step = stepSearch
	case "2":
		if system.HasRootAccess() {
			w.backend = "ruri"
			w.step = stepSearch
		} else {
			w.message = "Root access not available"
		}
	case "esc":
		w.step = stepSource
	}
	return w, textinput.Blink
}

func (w *installWizard) handleSearchStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "enter":
		return w, w.performSearch()
	case "up", "k":
		if w.selectedImg > 0 {
			w.selectedImg--
		}
	case "down", "j":
		if w.selectedImg < len(w.images)-1 {
			w.selectedImg++
		}
	case "tab":
		if len(w.images) > 0 {
			return w, w.loadTags()
		}
	case "esc":
		w.step = stepBackend
		w.images = nil
	default:
		w.searchInput, cmd = w.searchInput.Update(msg)
		return w, cmd
	}

	return w, nil
}

func (w *installWizard) handleTagStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if w.selectedTag > 0 {
			w.selectedTag--
		}
	case "down", "j":
		if w.selectedTag < len(w.tags)-1 {
			w.selectedTag++
		}
	case "enter":
		return w, w.startInstall()
	case "esc":
		w.step = stepSearch
		w.tags = nil
	}
	return w, nil
}

func (w *installWizard) performSearch() tea.Cmd {
	return func() tea.Msg {
		query := w.searchInput.Value()
		if query == "" {
			return nil
		}

		if w.source == "docker" {
			images, err := w.rurima.DockerSearch(query)
			if err != nil {
				w.err = err
				return nil
			}
			w.images = make([]interface{}, len(images))
			for i, img := range images {
				w.images[i] = img
			}
		} else {
			distros, err := w.rurima.LXCList()
			if err != nil {
				w.err = err
				return nil
			}
			// Filter by query
			var filtered []rootfs.LXCDistro
			for _, d := range distros {
				if strings.Contains(strings.ToLower(d.Name), strings.ToLower(query)) {
					filtered = append(filtered, d)
				}
			}
			w.images = make([]interface{}, len(filtered))
			for i, d := range filtered {
				w.images[i] = d
			}
		}

		w.selectedImg = 0
		return nil
	}
}

func (w *installWizard) loadTags() tea.Cmd {
	return func() tea.Msg {
		if w.selectedImg >= len(w.images) {
			return nil
		}

		if w.source == "docker" {
			img := w.images[w.selectedImg].(rootfs.DockerImage)
			tags, err := w.rurima.DockerTags(img.Name)
			if err != nil {
				w.err = err
				return nil
			}
			w.tags = make([]interface{}, len(tags))
			for i, tag := range tags {
				w.tags[i] = tag
			}
		} else {
			distro := w.images[w.selectedImg].(rootfs.LXCDistro)
			versions, err := w.rurima.LXCSearch(distro.Name)
			if err != nil {
				w.err = err
				return nil
			}
			w.tags = make([]interface{}, len(versions))
			for i, ver := range versions {
				w.tags[i] = ver
			}
		}

		w.selectedTag = 0
		w.step = stepTag
		return nil
	}
}

func (w *installWizard) startInstall() tea.Cmd {
	return func() tea.Msg {
		if w.selectedTag >= len(w.tags) {
			return installDoneMsg{err: fmt.Errorf("invalid tag selection")}
		}

		// Generate container name
		timestamp := time.Now().Unix()
		var name, imageName, tag string

		if w.source == "docker" {
			img := w.images[w.selectedImg].(rootfs.DockerImage)
			tagObj := w.tags[w.selectedTag].(rootfs.DockerTag)
			imageName = img.Name
			tag = tagObj.Tag
			name = fmt.Sprintf("%s-%s-%d",
				strings.ReplaceAll(imageName, "/", "_"), tag, timestamp)
		} else {
			distro := w.images[w.selectedImg].(rootfs.LXCDistro)
			version := w.tags[w.selectedTag].(rootfs.LXCVersion)
			imageName = distro.Name
			tag = version.Version
			name = fmt.Sprintf("%s-%s-%d", imageName, tag, timestamp)
		}

		savePath := filepath.Join(system.HomeDir, name)

		// Start installation
		w.step = stepInstall
		w.progressVal = 0

		progressChan := make(chan int, 100)
		go func() {
			for p := range progressChan {
				// Send progress updates
				_ = p // TODO: send to tea.Msg
			}
		}()

		var err error
		if w.source == "docker" {
			err = w.rurima.DockerPull(imageName, tag, savePath, progressChan)
		} else {
			err = w.rurima.LXCPull(imageName, tag, savePath, progressChan)
		}

		close(progressChan)

		if err != nil {
			return installDoneMsg{err: err}
		}

		// Run fixup
		if err := rootfs.RunFixup(savePath, w.backend); err != nil {
			return installDoneMsg{err: fmt.Errorf("fixup failed: %w", err)}
		}

		// Register container
		if err := w.manager.Register(name, savePath, w.backend); err != nil {
			return installDoneMsg{err: fmt.Errorf("failed to register: %w", err)}
		}

		return installDoneMsg{name: name}
	}
}

func (w *installWizard) View() string {
	var b strings.Builder

	title := TitleStyle.Render(fmt.Sprintf("Install Container (%d/5)", int(w.step)+1))
	b.WriteString(title)
	b.WriteString("\n\n")

	if w.err != nil {
		b.WriteString(ErrorStyle.Render(fmt.Sprintf("Error: %v", w.err)))
		b.WriteString("\n\n")
	}

	if w.message != "" {
		b.WriteString(InfoStyle.Render(w.message))
		b.WriteString("\n\n")
	}

	switch w.step {
	case stepSource:
		b.WriteString(w.viewSourceStep())
	case stepBackend:
		b.WriteString(w.viewBackendStep())
	case stepSearch:
		b.WriteString(w.viewSearchStep())
	case stepTag:
		b.WriteString(w.viewTagStep())
	case stepInstall:
		b.WriteString(w.viewInstallStep())
	case stepDone:
		b.WriteString(w.viewDoneStep())
	}

	return BorderStyle.Render(b.String())
}

func (w *installWizard) viewSourceStep() string {
	return SubtitleStyle.Render("Choose rootfs source:") + "\n\n" +
		"  [1] Docker Hub\n" +
		"  [2] LXC Mirror\n\n" +
		HelpStyle.Render("1/2 select • esc back")
}

func (w *installWizard) viewBackendStep() string {
	var b strings.Builder
	b.WriteString(SubtitleStyle.Render("Choose backend:"))
	b.WriteString("\n\n")
	b.WriteString("  [1] proot (no root required)\n")

	if system.HasRootAccess() {
		b.WriteString("  [2] ruri (requires root)\n")
	} else {
		b.WriteString("  [2] ruri ✗ Root not available\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("1/2 select • esc back"))
	return b.String()
}

func (w *installWizard) viewSearchStep() string {
	var b strings.Builder

	source := "Docker Hub"
	if w.source == "lxc" {
		source = "LXC"
	}

	b.WriteString(SubtitleStyle.Render(fmt.Sprintf("Search %s:", source)))
	b.WriteString("\n\n")
	b.WriteString(w.searchInput.View())
	b.WriteString("\n\n")

	if len(w.images) > 0 {
		b.WriteString(SubtitleStyle.Render("Results:"))
		b.WriteString("\n")

		for i, img := range w.images {
			cursor := "  "
			style := NormalStyle

			if i == w.selectedImg {
				cursor = CursorStyle.Render("▸ ")
				style = SelectedStyle
			}

			var line string
			if w.source == "docker" {
				dimg := img.(rootfs.DockerImage)
				line = fmt.Sprintf("%-30s %s", dimg.Name, dimg.Description)
			} else {
				distro := img.(rootfs.LXCDistro)
				line = distro.Name
			}

			b.WriteString(cursor + style.Render(line) + "\n")
		}

		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("enter search • ↑↓ select • tab next • esc back"))
	} else {
		b.WriteString(HelpStyle.Render("enter search • esc back"))
	}

	return b.String()
}

func (w *installWizard) viewTagStep() string {
	var b strings.Builder
	b.WriteString(SubtitleStyle.Render("Select tag/version:"))
	b.WriteString("\n\n")

	for i, tag := range w.tags {
		cursor := "  "
		style := NormalStyle

		if i == w.selectedTag {
			cursor = CursorStyle.Render("▸ ")
			style = SelectedStyle
		}

		var line string
		if w.source == "docker" {
			t := tag.(rootfs.DockerTag)
			line = t.Tag
		} else {
			v := tag.(rootfs.LXCVersion)
			line = fmt.Sprintf("%s (%s, %s)", v.Version, v.Variant, v.Arch)
		}

		b.WriteString(cursor + style.Render(line) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑↓ select • enter install • esc back"))
	return b.String()
}

func (w *installWizard) viewInstallStep() string {
	return SubtitleStyle.Render("Installing...") + "\n\n" +
		w.progress.ViewAs(w.progressVal) + "\n\n" +
		HelpStyle.Render("Please wait...")
}

func (w *installWizard) viewDoneStep() string {
	return SuccessStyle.Render("✓ Installation complete!") + "\n\n" +
		HelpStyle.Render("enter continue")
}
