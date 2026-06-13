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

package system

import (
	"os"
	"os/exec"
)

// IsRoot checks if the current user has root privileges
func IsRoot() bool {
	return os.Geteuid() == 0
}

// HasRootAccess checks if we can execute commands with sudo/su
func HasRootAccess() bool {
	// Try sudo first
	cmd := exec.Command("sudo", "-n", "true")
	if err := cmd.Run(); err == nil {
		return true
	}

	// Try su
	cmd = exec.Command("su", "-c", "true")
	if err := cmd.Run(); err == nil {
		return true
	}

	return false
}

// CommandExists checks if a command is available in PATH
func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
