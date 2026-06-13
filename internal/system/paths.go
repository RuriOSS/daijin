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
	"path/filepath"
)

// Termux paths
var (
	// PrefixDir is the Termux prefix directory
	PrefixDir = getEnvOr("PREFIX", "/data/data/com.termux/files/usr")

	// BinDir is the binary directory
	BinDir = filepath.Join(PrefixDir, "bin")

	// ShareDir is the shared data directory for daijin
	ShareDir = filepath.Join(PrefixDir, "share", "daijin")

	// VarDir is the variable data directory for daijin
	VarDir = filepath.Join(PrefixDir, "var", "daijin")

	// ContainersDir is where container configs are stored
	ContainersDir = filepath.Join(VarDir, "containers")

	// EtcDir is the config directory
	EtcDir = filepath.Join(PrefixDir, "etc")

	// ProcDir is where dummy proc files are stored
	ProcDir = filepath.Join(ShareDir, "proc")

	// HomeDir is the Termux home directory
	HomeDir = getEnvOr("HOME", "/data/data/com.termux/files/home")
)

// getEnvOr returns the environment variable value or a default
func getEnvOr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// EnsureDirs creates necessary directories if they don't exist
func EnsureDirs() error {
	dirs := []string{
		VarDir,
		ContainersDir,
		ShareDir,
		ProcDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
