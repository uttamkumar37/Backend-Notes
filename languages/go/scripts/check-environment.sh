#!/usr/bin/env bash
# Inventory only. Compatible with the Bash 3.2 supplied by macOS.
# Do not add installers, VM startup, environment dumps, or credential checks here.
set -u

run_check() {
    printf '\n$'
    printf ' %s' "$@"
    printf '\n'
    if "$@"; then
        :
    else
        local check_status=$?
        printf 'CHECK FAILED (exit %s); continuing inventory.\n' "$check_status"
    fi
}

version_check() {
    local tool_name="$1"
    if command -v "$tool_name" >/dev/null 2>&1; then
        run_check "$@"
    else
        printf '\n%s: NOT FOUND ON PATH\n' "$tool_name"
    fi
}

printf 'Go lab: read-only macOS tool inventory\n'
run_check sw_vers
run_check uname -m
run_check xcode-select -p
version_check brew --version

if command -v go >/dev/null 2>&1; then
    # Use the installed SDK without automatically downloading another toolchain.
    run_check env GOTOOLCHAIN=local go version
    run_check env GOTOOLCHAIN=local go env GOOS GOARCH
else
    printf '\ngo: NOT FOUND ON PATH\n'
fi

version_check git --version
version_check podman --version
if command -v podman >/dev/null 2>&1; then
    run_check podman machine list
fi

printf '\nOptional tools: availability only; missing tools can wait for their phase.\n'
for tool_name in gh psql valkey-cli redis-cli curl jq make dlv golangci-lint kubectl oc kind helm code goland vault prometheus grafana bruno; do
    if tool_path=$(command -v "$tool_name" 2>/dev/null); then
        printf '%-16s %s\n' "$tool_name" "$tool_path"
    else
        printf '%-16s NOT FOUND ON PATH\n' "$tool_name"
    fi
done

printf '\nInventory finished. Review CHECK FAILED lines; exit 0 is not a readiness certificate.\n'
printf 'GUI apps may be installed even when their terminal launcher is absent.\n'
printf 'No installs, VM startup, configuration edits, or secret reads were requested.\n'
