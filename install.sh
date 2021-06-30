#!/usr/bin/env bash

set -o errexit  \
    -o nounset  \
    -o pipefail

function check_for_go () {
    if which go >/dev/null 2>/dev/null; then
        echo "Found go compiler." >&2
    elif [[ -f "/etc/profile.d/modules.sh" ]]; then
        echo "No go compiler found, trying a module setup..." >&2
        source /etc/profile.d/modules.sh
        module purge
        module load gcc-libs
        module load compilers/go/1.16.5
    else
        echo "Could not get a go compiler, exiting..." >&2
        exit 1
    fi
}

echo "Checking go environment..." >&2
check_for_go

if [[ -n "${TRAVIS:-}" ]]; then
  INSTALL_PATH="$(mktemp -d)"
fi

install_path="${INSTALL_PATH:-/shared/ucl/apps/cluster-bin}"

echo "Changing into \"$(dirname -- "$0")\"..." >&2
cd "$(dirname -- "$0")"

echo "Recreating go.mod and go.sum..." >&2
go mod init github.com/UCL-RITS/CragPolicyExceptionCmd
go mod tidy

echo "Building..." >&2
./build.sh

echo "Installing to: $install_path" >&2
cp -vf bin/* "$install_path"/

echo "Done." >&2

