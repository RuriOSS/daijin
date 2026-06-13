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
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/moe-hacker/daijin/internal/system"
)

// RurimaClient wraps rurima binary
type RurimaClient struct {
	binaryPath string
}

// NewRurimaClient creates a new rurima client
func NewRurimaClient() *RurimaClient {
	path := "rurima"
	if system.CommandExists("rurima") {
		return &RurimaClient{binaryPath: path}
	}
	return nil
}

// IsAvailable checks if rurima is installed
func (r *RurimaClient) IsAvailable() bool {
	return r != nil && system.CommandExists(r.binaryPath)
}

// DockerImage represents a Docker image search result
type DockerImage struct {
	Name        string
	Description string
}

// DockerTag represents a Docker image tag
type DockerTag struct {
	Tag  string
	Size string
}

// LXCDistro represents an LXC distribution
type LXCDistro struct {
	Name string
}

// LXCVersion represents an LXC version
type LXCVersion struct {
	Version string
	Variant string
	Arch    string
}

// DockerSearch searches for Docker images
func (r *RurimaClient) DockerSearch(query string) ([]DockerImage, error) {
	cmd := exec.Command(r.binaryPath, "docker", "search", "-i", query, "-q")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker search failed: %w", err)
	}

	var images []DockerImage
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		// Format: "image_name    Description: description text"
		if strings.Contains(line, "Description:") {
			continue // Skip header
		}

		parts := strings.Fields(line)
		if len(parts) > 0 {
			images = append(images, DockerImage{
				Name:        parts[0],
				Description: strings.Join(parts[1:], " "),
			})
		}
	}

	return images, nil
}

// DockerTags gets available tags for a Docker image
func (r *RurimaClient) DockerTags(image string) ([]DockerTag, error) {
	cmd := exec.Command(r.binaryPath, "docker", "tag", "-i", image, "-q")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker tag failed: %w", err)
	}

	var tags []DockerTag
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			tags = append(tags, DockerTag{
				Tag:  parts[1],
				Size: "", // rurima doesn't provide size in tag list
			})
		}
	}

	return tags, nil
}

// DockerPull pulls a Docker image
func (r *RurimaClient) DockerPull(image, tag, savePath string, progressChan chan<- int) error {
	cmd := exec.Command(r.binaryPath, "docker", "pull", "-q", "-i", image, "-t", tag, "-s", savePath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Parse progress output
	// Format: [///////////////...] 68.50%
	progressRe := regexp.MustCompile(`\[.*?\]\s+(\d+(?:\.\d+)?)%`)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := progressRe.FindStringSubmatch(line); len(matches) > 1 {
			if percent, err := strconv.ParseFloat(matches[1], 64); err == nil {
				if progressChan != nil {
					progressChan <- int(percent)
				}
			}
		}
	}

	return cmd.Wait()
}

// LXCList lists available LXC distributions
func (r *RurimaClient) LXCList() ([]LXCDistro, error) {
	cmd := exec.Command(r.binaryPath, "lxc", "list", "-q")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("lxc list failed: %w", err)
	}

	var distros []LXCDistro
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) > 0 {
			distro := parts[0]
			if !seen[distro] {
				distros = append(distros, LXCDistro{Name: distro})
				seen[distro] = true
			}
		}
	}

	return distros, nil
}

// LXCSearch searches for LXC versions of a distro
func (r *RurimaClient) LXCSearch(distro string) ([]LXCVersion, error) {
	cmd := exec.Command(r.binaryPath, "lxc", "search", "-q", "-o", distro)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("lxc search failed: %w", err)
	}

	var versions []LXCVersion
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			versions = append(versions, LXCVersion{
				Version: parts[2],
				Variant: parts[1],
				Arch:    parts[0],
			})
		}
	}

	return versions, nil
}

// LXCPull pulls an LXC image
func (r *RurimaClient) LXCPull(distro, version, savePath string, progressChan chan<- int) error {
	cmd := exec.Command(r.binaryPath, "lxc", "pull", "-o", distro, "-v", version, "-s", savePath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Parse progress (same format as docker)
	progressRe := regexp.MustCompile(`\[.*?\]\s+(\d+(?:\.\d+)?)%`)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := progressRe.FindStringSubmatch(line); len(matches) > 1 {
			if percent, err := strconv.ParseFloat(matches[1], 64); err == nil {
				if progressChan != nil {
					progressChan <- int(percent)
				}
			}
		}
	}

	return cmd.Wait()
}
