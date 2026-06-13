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

echo -e "${COLOR}Building Daijin v2.0...${ENDCOLOR}"

# Install dependencies
printf "${COLOR}Installing dependencies...${ENDCOLOR}\n"
pkg install ndk-multilib-native-static tsu coreutils p7zip gettext tar unzip zip git wget dpkg curl nano proot axel termux-tools util-linux pv gawk clang ndk-sysroot ndk-multilib libc-client-static libcap-static binutils libseccomp-static golang make

# Update submodule
printf "${COLOR}Initializing submodules...${ENDCOLOR}\n"
git submodule update --init

# Create build dir
printf "${COLOR}Creating build directory...${ENDCOLOR}\n"
mkdir -p build/DEBIAN
mkdir -p build/data/data/com.termux/files/usr/bin
mkdir -p build/data/data/com.termux/files/usr/share/daijin/proc
mkdir -p build/data/data/com.termux/files/usr/etc

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
echo "echo -e \"\033[33mruri is built-in in rurima now, please use \033[32mrurima r\033[33m instead\033[0m\"" >../../build/data/data/com.termux/files/usr/bin/ruri
chmod 777 ../../build/data/data/com.termux/files/usr/bin/ruri

# Return to root dir
cd ../..

# Copy rurima config file
cp src/rurima.conf build/data/data/com.termux/files/usr/etc/rurima.conf

# Decompress dummy files of procfs
tar -xf src/share/proc.tar.xz -C build/data/data/com.termux/files/usr/share/daijin/proc/

# Compile Go daijin
printf "${COLOR}Compiling daijin (Go)...${ENDCOLOR}\n"
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty)" \
    -o build/data/data/com.termux/files/usr/bin/daijin \
    ./cmd/daijin

# Fix permission
chmod 777 build/data/data/com.termux/files/usr/bin/*

cd build

# Set build info
size=$(du -s . | awk '{printf $1}')
sed -i "s/\[size\]/${size}/" DEBIAN/control
arch=$(dpkg --print-architecture)
sed -i "s/\[arch\]/${arch}/" DEBIAN/control

# Build deb
dpkg -b . ../daijin-${arch}.deb

# Clean
cd ..
rm -rf build

# Done
printf "${COLOR}Build complete: daijin-${arch}.deb${ENDCOLOR}\n"
printf "${COLOR}Binary size: $(du -h daijin-${arch}.deb | cut -f1)${ENDCOLOR}\n"
