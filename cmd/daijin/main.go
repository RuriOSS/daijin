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
	"github.com/moe-hacker/daijin/internal/rootfs"
	"github.com/moe-hacker/daijin/internal/system"
	"github.com/moe-hacker/daijin/internal/ui"
)

const version = "2.0.0-dev"

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

	// Ensure necessary directories exist
	if err := system.EnsureDirs(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directories: %v\n", err)
		os.Exit(1)
	}

	// Ensure rurima binary is available
	if err := rootfs.EnsureRurima(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to extract rurima: %v\n", err)
		os.Exit(1)
	}

	// Start the TUI
	model := ui.NewContainerListModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
