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
)

// CleanRootfs removes all contents inside a container directory
// but keeps the directory itself
func CleanRootfs(containerDir string) error {
	// Check if directory exists
	info, err := os.Stat(containerDir)
	if err != nil {
		return fmt.Errorf("failed to stat directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", containerDir)
	}

	// Read all entries
	entries, err := os.ReadDir(containerDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Remove each entry
	for _, entry := range entries {
		path := containerDir + "/" + entry.Name()
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}

	return nil
}

// CleanRootfsByName cleans a container's rootfs by name
func (m *Manager) CleanRootfs(name string) error {
	config, err := m.Load(name)
	if err != nil {
		return fmt.Errorf("failed to load container: %w", err)
	}

	return CleanRootfs(config.ContainerDir)
}
