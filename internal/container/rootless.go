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
)

// RootlessBackend implements rootless containers using user namespaces
type RootlessBackend struct{}

// NewRootlessBackend creates a new rootless backend
func NewRootlessBackend() *RootlessBackend {
	return &RootlessBackend{}
}

// IsAvailable checks if rootless containers are supported
func (r *RootlessBackend) IsAvailable() bool {
	// Check if unshare command exists
	if _, err := exec.LookPath("unshare"); err != nil {
		return false
	}

	// Check if user namespaces are enabled
	// Try to create a user namespace
	cmd := exec.Command("unshare", "--user", "--map-root-user", "true")
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

// Start implements the Backend interface
func (r *RootlessBackend) Start(config *Config, command []string) error {
	if len(command) == 0 {
		command = []string{"/bin/sh"}
	}

	// Build unshare command with user namespace mapping
	args := []string{
		"--user",           // Create user namespace
		"--map-root-user",  // Map current user to root inside namespace
		"--mount",          // Create mount namespace
		"--pid",            // Create PID namespace
		"--fork",           // Fork before exec
		"--",               // End of unshare options
		"chroot",           // Chroot into container
		config.ContainerDir, // Container root directory
	}
	args = append(args, command...)

	execCmd := exec.Command("unshare", args...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}

// Stop stops a running container
func (r *RootlessBackend) Stop(config *Config) error {
	// For rootless containers, we need to track PIDs
	// This is a simplified version
	return fmt.Errorf("stop not implemented for rootless backend yet")
}

// Remove removes a container
func (r *RootlessBackend) Remove(config *Config) error {
	// Remove the container directory
	// Rootless containers don't need special cleanup
	return nil
}

// RequiresRoot returns false since rootless doesn't need root
func (r *RootlessBackend) RequiresRoot() bool {
	return false
}
