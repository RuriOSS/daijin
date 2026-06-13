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
	"path/filepath"
	"strings"

	"github.com/moe-hacker/daijin/internal/system"
)

// Manager handles container operations
type Manager struct {
	configDir string
	backends  map[string]Backend
}

// NewManager creates a new container manager
func NewManager(configDir string) *Manager {
	return &Manager{
		configDir: configDir,
		backends: map[string]Backend{
			BackendProot: NewProotBackend(),
			BackendRuri:  NewRuriBackend(),
		},
	}
}

// List returns all registered containers
func (m *Manager) List() ([]*Config, error) {
	entries, err := os.ReadDir(m.configDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Config dir doesn't exist yet, return empty list
			return []*Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var configs []*Config
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".conf") {
			continue
		}

		configPath := filepath.Join(m.configDir, entry.Name())
		config, err := LoadConfig(configPath)
		if err != nil {
			// Skip invalid configs but don't fail
			continue
		}

		configs = append(configs, config)
	}

	return configs, nil
}

// Load loads a specific container config by name
func (m *Manager) Load(name string) (*Config, error) {
	configPath := filepath.Join(m.configDir, name+".conf")
	return LoadConfig(configPath)
}

// Save saves a container config
func (m *Manager) Save(config *Config) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	configPath := filepath.Join(m.configDir, config.Name+".conf")
	config.ConfigPath = configPath

	return config.Save(configPath)
}

// Start starts a container
func (m *Manager) Start(name string, command []string) error {
	config, err := m.Load(name)
	if err != nil {
		return fmt.Errorf("failed to load container config: %w", err)
	}

	backend, ok := m.backends[config.Backend]
	if !ok {
		return fmt.Errorf("unknown backend: %s", config.Backend)
	}

	if !backend.IsAvailable() {
		return fmt.Errorf("backend %s is not available", config.Backend)
	}

	if backend.RequiresRoot() && !system.HasRootAccess() {
		return fmt.Errorf("backend %s requires root access", config.Backend)
	}

	return backend.Start(config, command)
}

// Remove deletes a container
func (m *Manager) Remove(name string) error {
	config, err := m.Load(name)
	if err != nil {
		return fmt.Errorf("failed to load container config: %w", err)
	}

	backend, ok := m.backends[config.Backend]
	if !ok {
		return fmt.Errorf("unknown backend: %s", config.Backend)
	}

	if err := backend.Remove(config); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	// Remove config file
	if err := os.Remove(config.ConfigPath); err != nil {
		return fmt.Errorf("failed to remove config file: %w", err)
	}

	return nil
}

// Register registers an existing container directory
func (m *Manager) Register(name, containerDir, backend string) error {
	// Check if name is already taken
	if _, err := m.Load(name); err == nil {
		return fmt.Errorf("container name already exists: %s", name)
	}

	// Validate backend
	if backend != BackendProot && backend != BackendRuri {
		return fmt.Errorf("invalid backend: %s", backend)
	}

	// Check if directory exists
	if _, err := os.Stat(containerDir); os.IsNotExist(err) {
		return fmt.Errorf("container directory does not exist: %s", containerDir)
	}

	// Get absolute path
	absPath, err := filepath.Abs(containerDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	config := &Config{
		Name:         name,
		Backend:      backend,
		ContainerDir: absPath,
	}

	return m.Save(config)
}

// Exists checks if a container exists
func (m *Manager) Exists(name string) bool {
	_, err := m.Load(name)
	return err == nil
}
