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

package container

import (
	"github.com/moe-hacker/daijin/internal/system"
)

// Backend defines the interface for container backends (rootless, ruri, proot)
type Backend interface {
	// Start starts the container with optional command
	Start(config *Config, command []string) error

	// Stop stops a running container
	Stop(config *Config) error

	// Remove removes the container and its data
	Remove(config *Config) error

	// IsAvailable checks if the backend is available on the system
	IsAvailable() bool

	// RequiresRoot returns true if the backend requires root privileges
	RequiresRoot() bool
}

// DetectBestBackend automatically selects the best available backend
func DetectBestBackend() string {
	if system.IsLinux {
		// Linux: prefer rootless with user namespaces
		rootless := NewRootlessBackend()
		if rootless.IsAvailable() {
			return "rootless"
		}

		// Fallback to ruri if we have root
		ruri := NewRuriBackend()
		if ruri.IsAvailable() {
			return "ruri"
		}

		// Last resort: proot (not recommended on Linux)
		proot := NewProotBackend()
		if proot.IsAvailable() {
			return "proot"
		}
	} else if system.IsAndroid {
		// Android/Termux: prefer ruri if available
		ruri := NewRuriBackend()
		if ruri.IsAvailable() {
			return "ruri"
		}

		// Fallback to proot
		proot := NewProotBackend()
		if proot.IsAvailable() {
			return "proot"
		}
	}

	// Default fallback
	return "proot"
}
