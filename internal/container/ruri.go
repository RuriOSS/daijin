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
	"fmt"
	"os"
	"os/exec"

	"github.com/moe-hacker/daijin/internal/system"
)

// RuriBackend implements the Backend interface for ruri
type RuriBackend struct{}

// NewRuriBackend creates a new ruri backend
func NewRuriBackend() Backend {
	return &RuriBackend{}
}

// IsAvailable checks if rurima (which includes ruri) is available
func (r *RuriBackend) IsAvailable() bool {
	return system.CommandExists("rurima")
}

// RequiresRoot returns true as ruri needs root
func (r *RuriBackend) RequiresRoot() bool {
	return true
}

// Start starts a container using ruri (via rurima)
func (r *RuriBackend) Start(config *Config, command []string) error {
	// Mount /data as suid (Android only)
	if system.IsAndroid {
		_ = exec.Command("sudo", "mount", "-o", "remount,suid", "/data").Run()
	}

	// Use rurima r (ruri) to start the container
	args := []string{"rurima", "r"}

	// Unset LD_PRELOAD on Android only
	if system.IsAndroid {
		// On Android, we need to unset LD_PRELOAD
		args = []string{"LD_PRELOAD=", "rurima", "r"}
	}

	if config.ConfigPath != "" {
		args = append(args, "-c", config.ConfigPath)
	} else {
		args = append(args, config.ContainerDir)
	}

	if len(command) > 0 {
		args = append(args, command...)
	}

	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Stop stops a ruri container
func (r *RuriBackend) Stop(config *Config) error {
	return fmt.Errorf("stop is not yet implemented for ruri backend")
}

// Remove removes a ruri container
func (r *RuriBackend) Remove(config *Config) error {
	// Ask for confirmation
	fmt.Printf("Remove container directory %s? [y/N] ", config.ContainerDir)
	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		return fmt.Errorf("removal cancelled")
	}

	// Use sudo to remove
	cmd := exec.Command("sudo", "rm", "-rf", config.ContainerDir)
	return cmd.Run()
}
