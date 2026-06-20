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
	"runtime"
)

// Platform detection
var (
	IsAndroid = detectAndroid()
	IsLinux   = runtime.GOOS == "linux" && !IsAndroid
)

// detectAndroid checks if we're running on Android/Termux
func detectAndroid() bool {
	// Check for Termux-specific environment or paths
	if os.Getenv("TERMUX_VERSION") != "" {
		return true
	}
	if os.Getenv("PREFIX") != "" && os.Getenv("ANDROID_ROOT") != "" {
		return true
	}
	// Check if /data/data/com.termux exists
	if _, err := os.Stat("/data/data/com.termux"); err == nil {
		return true
	}
	return false
}
