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
	"path/filepath"

	"github.com/moe-hacker/daijin/internal/system"
)

// ProotBackend implements the Backend interface for proot
type ProotBackend struct{}

// NewProotBackend creates a new proot backend
func NewProotBackend() Backend {
	return &ProotBackend{}
}

// IsAvailable checks if proot is available
func (p *ProotBackend) IsAvailable() bool {
	return system.CommandExists("proot")
}

// RequiresRoot returns false as proot doesn't need root
func (p *ProotBackend) RequiresRoot() bool {
	return false
}

// Start starts a container using proot
func (p *ProotBackend) Start(config *Config, command []string) error {
	// Unset LD_PRELOAD on Android only (required by proot)
	if system.IsAndroid {
		os.Unsetenv("LD_PRELOAD")
	}

	args := []string{
		"--link2symlink",
		"--kill-on-exit",
		"--sysvipc",
		"-L",
		"-0",
		"-r", config.ContainerDir,
		"-b", "/dev",
		"-b", "/sys",
		"-b", "/proc",
		"-w", "/root",
	}

	// Android-specific proot options
	if system.IsAndroid {
		args = append(args[:4], "--ashmem-memfd")
		args = append(args, args[4:]...)
	}

	// Mount proc dummy files (Android only - Linux has real /proc)
	if system.IsAndroid {
		procMounts := []string{
			"buddyinfo", "cgroups", "consoles", "crypto", "devices",
			"diskstats", "execdomains", "fb", "filesystems", "interrupts",
			"iomem", "ioports", "kallsyms", "key-users", "keys",
			"kpageflags", "loadavg", "locks", "misc", "modules",
			"pagetypeinfo", "partitions", "sched_debug", "softirqs",
			"stat", "timer_list", "uptime", "version", "vmallocinfo",
			"vmstat", "zoneinfo",
		}

		for _, mount := range procMounts {
			src := filepath.Join(system.ProcDir, mount)
			dst := filepath.Join("/proc", mount)
			args = append(args, fmt.Sprintf("--mount=%s:%s", src, dst))
		}
	}

	// Mount tmpdir
	tmpdir := os.Getenv("TMPDIR")
	if tmpdir == "" {
		if system.IsAndroid {
			tmpdir = filepath.Join(system.PrefixDir, "tmp")
		} else {
			tmpdir = "/tmp"
		}
	}
	args = append(args, fmt.Sprintf("--mount=%s:/tmp", tmpdir))

	// Mount container dir
	args = append(args, fmt.Sprintf("--mount=%s:/", config.ContainerDir))

	// Add extra args if any
	if config.ExtraArgs != "" {
		args = append(args, config.ExtraArgs)
	}

	// Add command
	if len(command) == 0 {
		// Default command
		suPath := filepath.Join(config.ContainerDir, "bin", "su")
		if _, err := os.Stat(suPath); err == nil {
			args = append(args, "/bin/su", "-", "root")
		} else {
			args = append(args, "/bin/sh")
		}
	} else {
		args = append(args, command...)
	}

	cmd := exec.Command("proot", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Stop stops a proot container (not applicable for proot)
func (p *ProotBackend) Stop(config *Config) error {
	return fmt.Errorf("stop is not supported for proot backend")
}

// Remove removes a proot container
func (p *ProotBackend) Remove(config *Config) error {
	// Ask for confirmation
	fmt.Printf("Remove container directory %s? [y/N] ", config.ContainerDir)
	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		return fmt.Errorf("removal cancelled")
	}

	// Remove the container directory
	return os.RemoveAll(config.ContainerDir)
}
