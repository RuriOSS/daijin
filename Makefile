# Copyright 2024 moe-hacker
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
.PHONY: all build build-android build-android-arm test run clean format deb

# Version from git
VERSION := $(shell git describe --tags --always --dirty)

all: build

# Build for current platform (development)
build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o daijin ./cmd/daijin

# Build for Android ARM64 (production - most Android devices)
build-android:
	GOOS=android GOARCH=arm64 CGO_ENABLED=0 go build \
		-ldflags="-s -w -X main.version=$(VERSION)" \
		-o daijin-android-arm64 ./cmd/daijin

# Build for Linux x86_64 (for testing in emulators or x86 environments)
build-linux-x64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags="-s -w -X main.version=$(VERSION)" \
		-o daijin-linux-x64 ./cmd/daijin

# Build for Linux ARM64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
		-ldflags="-s -w -X main.version=$(VERSION)" \
		-o daijin-linux-arm64 ./cmd/daijin

test:
	go test -v ./...

run:
	go run ./cmd/daijin

clean:
	rm -f daijin daijin-android-*
	rm -rf build/

format:
	gofmt -s -w .
	go mod tidy

deb:
	./build.sh
