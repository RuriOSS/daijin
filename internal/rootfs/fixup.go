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

package rootfs

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// FixupScript contains the embedded fixup.sh content
//
//go:embed fixup.sh
var FixupScript string

// RunFixup runs the fixup script inside a container
func RunFixup(containerDir string, backend string) error {
	// Write fixup script to container's /tmp
	tmpPath := filepath.Join(containerDir, "tmp", "fixup.sh")

	if err := os.MkdirAll(filepath.Dir(tmpPath), 0755); err != nil {
		return fmt.Errorf("failed to create tmp dir: %w", err)
	}

	if err := os.WriteFile(tmpPath, []byte(FixupScript), 0755); err != nil {
		return fmt.Errorf("failed to write fixup script: %w", err)
	}

	// Run the script based on backend
	var cmd *exec.Cmd

	switch backend {
	case "proot":
		cmd = exec.Command("proot", "-r", containerDir, "/tmp/fixup.sh")
	case "ruri":
		cmd = exec.Command("sudo", "LD_PRELOAD=", "rurima", "r", containerDir, "/tmp/fixup.sh")
	default:
		return fmt.Errorf("unknown backend: %s", backend)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
