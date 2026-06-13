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

// Backend defines the interface for container backends (proot, ruri)
type Backend interface {
	// Start starts the container with optional command
	Start(config *Config, command []string) error

	// Stop stops a running container
	Stop(config *Config) error

	// Remove removes the container and its data
	Remove(config *Config) error

	// IsAvailable checks if the backend is available on the system
	IsAvailable() bool

	// RequiresRoot returns true if the backend requires root privileges
	RequiresRoot() bool
}
