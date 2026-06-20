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

// Platform-specific paths
var (
	// PrefixDir is the prefix directory
	PrefixDir = initPrefixDir()

	// BinDir is the binary directory
	BinDir = filepath.Join(PrefixDir, "bin")

	// ShareDir is the shared data directory for daijin
	ShareDir = initShareDir()

	// VarDir is the variable data directory for daijin
	VarDir = initVarDir()

	// ContainersDir is where container configs are stored
	ContainersDir = filepath.Join(VarDir, "containers")

	// EtcDir is the config directory
	EtcDir = initEtcDir()

	// ProcDir is where dummy proc files are stored
	ProcDir = filepath.Join(ShareDir, "proc")

	// HomeDir is the home directory
	HomeDir = initHomeDir()
)

// initPrefixDir returns platform-specific prefix directory
func initPrefixDir() string {
	if IsAndroid {
		return getEnvOr("PREFIX", "/data/data/com.termux/files/usr")
	}
	// Linux: check if running as root or user
	if os.Geteuid() == 0 {
		return "/usr"
	}
	return filepath.Join(os.Getenv("HOME"), ".local")
}

// initShareDir returns platform-specific share directory
func initShareDir() string {
	if IsAndroid {
		return filepath.Join(PrefixDir, "share", "daijin")
	}
	// Linux FHS
	if os.Geteuid() == 0 {
		return "/usr/share/daijin"
	}
	return filepath.Join(os.Getenv("HOME"), ".local/share/daijin")
}

// initVarDir returns platform-specific var directory
func initVarDir() string {
	if IsAndroid {
		return filepath.Join(PrefixDir, "var", "daijin")
	}
	// Linux FHS
	if os.Geteuid() == 0 {
		return "/var/lib/daijin"
	}
	return filepath.Join(os.Getenv("HOME"), ".local/share/daijin/var")
}

// initEtcDir returns platform-specific etc directory
func initEtcDir() string {
	if IsAndroid {
		return filepath.Join(PrefixDir, "etc")
	}
	// Linux FHS
	if os.Geteuid() == 0 {
		return "/etc"
	}
	return filepath.Join(os.Getenv("HOME"), ".config/daijin")
}

// initHomeDir returns platform-specific home directory
func initHomeDir() string {
	if IsAndroid {
		return getEnvOr("HOME", "/data/data/com.termux/files/home")
	}
	// Linux: standard HOME
	return os.Getenv("HOME")
}

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
