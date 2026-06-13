#!/bin/bash
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

set -e

# Colors
COLOR="\033[1;38;2;254;228;208m"
ENDCOLOR="\033[0m"

echo -e "${COLOR}Building Daijin v2.0 (Go Edition)...${ENDCOLOR}"

# Install dependencies
printf "${COLOR}Installing dependencies...${ENDCOLOR}\n"
pkg install ndk-multilib-native-static tsu coreutils tar git wget dpkg proot termux-tools clang ndk-sysroot ndk-multilib libcap-static binutils libseccomp-static golang make

# Update submodules
printf "${COLOR}Initializing submodules...${ENDCOLOR}\n"
git submodule update --init --recursive

# Create build directories
printf "${COLOR}Creating build directory...${ENDCOLOR}\n"
mkdir -p build/DEBIAN
mkdir -p build/data/data/com.termux/files/usr/bin
mkdir -p build/data/data/com.termux/files/usr/share/daijin/proc
mkdir -p build/data/data/com.termux/files/usr/var/daijin/containers

# Copy dpkg config file
cp -r dpkg-conf/* build/DEBIAN/
chmod -R 755 build/DEBIAN

# Compile rurima
printf "${COLOR}Compiling rurima...${ENDCOLOR}\n"
cd src/rurima
git submodule update --init
./configure -s
make
cp rurima ../../build/data/data/com.termux/files/usr/bin/
# Create ruri symlink message
echo 'echo -e "\033[33mruri is built-in in rurima now, please use \033[32mrurima r\033[33m instead\033[0m"' >../../build/data/data/com.termux/files/usr/bin/ruri
chmod 755 ../../build/data/data/com.termux/files/usr/bin/ruri

# Return to root dir
cd ../..

# Decompress dummy proc files for proot
printf "${COLOR}Extracting proc files...${ENDCOLOR}\n"
tar -xf src/share/proc.tar.xz -C build/data/data/com.termux/files/usr/share/daijin/proc/

# Compile Go daijin
printf "${COLOR}Compiling daijin (Go)...${ENDCOLOR}\n"
VERSION=$(git describe --tags --always --dirty)
GOOS=android GOARCH=arm64 CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o build/data/data/com.termux/files/usr/bin/daijin \
    ./cmd/daijin

# Fix permissions
chmod 755 build/data/data/com.termux/files/usr/bin/*

# Build deb package
cd build

# Set build info
size=$(du -s . | awk '{printf $1}')
sed -i "s/\[size\]/${size}/" DEBIAN/control
arch=$(dpkg --print-architecture)
sed -i "s/\[arch\]/${arch}/" DEBIAN/control

# Build deb
printf "${COLOR}Building .deb package...${ENDCOLOR}\n"
dpkg -b . ../daijin-${arch}.deb

# Clean up
cd ..
rm -rf build

# Done
printf "${COLOR}Build complete: daijin-${arch}.deb${ENDCOLOR}\n"
printf "${COLOR}Binary size: $(du -h daijin-${arch}.deb | cut -f1)${ENDCOLOR}\n"
echo ""
echo "Install with: dpkg -i daijin-${arch}.deb"
echo "Run with: daijin"
