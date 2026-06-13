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
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file in old format
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")

	content := `backend="proot"
container_dir="/tmp/test-container"
extra_args="--verbose"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Test loading
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify fields
	if config.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", config.Name)
	}
	if config.Backend != "proot" {
		t.Errorf("Expected backend 'proot', got '%s'", config.Backend)
	}
	if config.ContainerDir != "/tmp/test-container" {
		t.Errorf("Expected container_dir '/tmp/test-container', got '%s'", config.ContainerDir)
	}
	if config.ExtraArgs != "--verbose" {
		t.Errorf("Expected extra_args '--verbose', got '%s'", config.ExtraArgs)
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")

	config := &Config{
		Name:         "test",
		Backend:      "ruri",
		ContainerDir: "/tmp/test",
		ExtraArgs:    "--debug",
	}

	// Save
	if err := config.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load back
	loaded, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig after save failed: %v", err)
	}

	// Verify
	if loaded.Backend != config.Backend {
		t.Errorf("Backend mismatch: expected %s, got %s", config.Backend, loaded.Backend)
	}
	if loaded.ContainerDir != config.ContainerDir {
		t.Errorf("ContainerDir mismatch: expected %s, got %s", config.ContainerDir, loaded.ContainerDir)
	}
	if loaded.ExtraArgs != config.ExtraArgs {
		t.Errorf("ExtraArgs mismatch: expected %s, got %s", config.ExtraArgs, loaded.ExtraArgs)
	}
}

func TestBackendValidation(t *testing.T) {
	tests := []struct {
		name    string
		backend string
		wantErr bool
	}{
		{"valid proot", "proot", false},
		{"valid ruri", "ruri", false},
		{"invalid backend", "invalid", true},
		{"empty backend", "", true},
	}

	tmpDir := t.TempDir()
	containerDir := filepath.Join(tmpDir, "container")
	os.MkdirAll(containerDir, 0755)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Name:         "test",
				Backend:      tt.backend,
				ContainerDir: containerDir,
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
