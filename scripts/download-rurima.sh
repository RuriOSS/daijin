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
RED="\033[31m"
ENDCOLOR="\033[0m"

# Check if output path is provided
if [ -z "$1" ]; then
    echo -e "${RED}Usage: $0 <output-path>${ENDCOLOR}"
    exit 1
fi

OUTPUT_PATH="$1"

# Detect architecture and map to rurima release naming
ARCH=$(dpkg --print-architecture 2>/dev/null || uname -m)

case ${ARCH} in
    arm64|aarch64)
        RURIMA_ARCH="aarch64"
        ;;
    armhf)
        RURIMA_ARCH="armhf"
        ;;
    armv7*|armv7l)
        RURIMA_ARCH="armv7"
        ;;
    amd64|x86_64)
        RURIMA_ARCH="x86_64"
        ;;
    i386|i686)
        RURIMA_ARCH="i386"
        ;;
    riscv64)
        RURIMA_ARCH="riscv64"
        ;;
    s390x)
        RURIMA_ARCH="s390x"
        ;;
    ppc64el|ppc64le)
        RURIMA_ARCH="ppc64le"
        ;;
    loongarch64|loong64)
        RURIMA_ARCH="loongarch64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: ${ARCH}${ENDCOLOR}"
        echo -e "${RED}Supported: arm64, armhf, armv7, amd64, i386, riscv64, s390x, ppc64el, loongarch64${ENDCOLOR}"
        exit 1
        ;;
esac

echo -e "${COLOR}Downloading rurima for ${RURIMA_ARCH}...${ENDCOLOR}"

# Download URL
RURIMA_URL="https://github.com/RuriOSS/rurima/releases/latest/download/${RURIMA_ARCH}.tar"

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap "rm -rf ${TMP_DIR}" EXIT

# Download the tar file
echo -e "${COLOR}Fetching from ${RURIMA_URL}${ENDCOLOR}"
if command -v wget >/dev/null 2>&1; then
    wget -O "${TMP_DIR}/rurima.tar" "${RURIMA_URL}" || {
        echo -e "${RED}Failed to download rurima${ENDCOLOR}"
        exit 1
    }
elif command -v curl >/dev/null 2>&1; then
    curl -L -o "${TMP_DIR}/rurima.tar" "${RURIMA_URL}" || {
        echo -e "${RED}Failed to download rurima${ENDCOLOR}"
        exit 1
    }
else
    echo -e "${RED}Neither wget nor curl found. Please install one of them.${ENDCOLOR}"
    exit 1
fi

# Extract the binary
echo -e "${COLOR}Extracting rurima binary...${ENDCOLOR}"
tar -xf "${TMP_DIR}/rurima.tar" -C "${TMP_DIR}/"

# Verify the binary exists
if [ ! -f "${TMP_DIR}/rurima" ]; then
    echo -e "${RED}Failed to extract rurima binary${ENDCOLOR}"
    exit 1
fi

# Create output directory if needed
OUTPUT_DIR=$(dirname "${OUTPUT_PATH}")
mkdir -p "${OUTPUT_DIR}"

# Copy to output location
cp "${TMP_DIR}/rurima" "${OUTPUT_PATH}"
chmod 755 "${OUTPUT_PATH}"

echo -e "${COLOR}Successfully downloaded rurima to ${OUTPUT_PATH}${ENDCOLOR}"

# Verify it's executable
if [ ! -x "${OUTPUT_PATH}" ]; then
    echo -e "${RED}Warning: rurima binary is not executable${ENDCOLOR}"
    exit 1
fi

echo -e "${COLOR}Done!${ENDCOLOR}"
