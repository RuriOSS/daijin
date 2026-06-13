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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Backend types
const (
	BackendProot = "proot"
	BackendRuri  = "ruri"
)

// Config represents a container configuration
type Config struct {
	Name         string
	Backend      string
	ContainerDir string
	ExtraArgs    string
	ConfigPath   string
}

// LoadConfig loads a container config from a .conf file
// Supports the old bash format for backward compatibility:
//   backend="proot"
//   container_dir="/path/to/container"
//   extra_args="..."
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer file.Close()

	config := &Config{
		Name:       strings.TrimSuffix(filepath.Base(path), ".conf"),
		ConfigPath: path,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key="value" format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"'")

		switch key {
		case "backend":
			config.Backend = value
		case "container_dir":
			config.ContainerDir = value
		case "extra_args":
			config.ExtraArgs = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Validate required fields
	if config.Backend == "" {
		return nil, fmt.Errorf("backend not specified in config")
	}
	if config.ContainerDir == "" {
		return nil, fmt.Errorf("container_dir not specified in config")
	}

	return config, nil
}

// SaveConfig saves a container config to a .conf file
func (c *Config) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}
	defer file.Close()

	// Write in the old format for backward compatibility
	lines := []string{
		fmt.Sprintf("backend=\"%s\"", c.Backend),
		fmt.Sprintf("container_dir=\"%s\"", c.ContainerDir),
	}

	if c.ExtraArgs != "" {
		lines = append(lines, fmt.Sprintf("extra_args=\"%s\"", c.ExtraArgs))
	}

	for _, line := range lines {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	return nil
}

// Validate checks if the container config is valid
func (c *Config) Validate() error {
	if c.Backend != BackendProot && c.Backend != BackendRuri {
		return fmt.Errorf("invalid backend: %s (must be proot or ruri)", c.Backend)
	}

	if c.ContainerDir == "" {
		return fmt.Errorf("container_dir is required")
	}

	// Check if container directory exists
	if _, err := os.Stat(c.ContainerDir); os.IsNotExist(err) {
		return fmt.Errorf("container directory does not exist: %s", c.ContainerDir)
	}

	return nil
}
