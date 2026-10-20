#!/usr/bin/env bash
# Run from any directory. Requires Go 1.27+; race mode also needs a C toolchain.
set -euo pipefail

go_root=$(cd "$(dirname "$0")" && pwd)
mode=${1:-fast}
if [ "$#" -gt 1 ] || { [ "$mode" != fast ] && [ "$mode" != race ]; }; then
    printf 'Usage: bash check.sh [fast|race]\n' >&2
    exit 2
fi
export GOTOOLCHAIN=local GOWORK=off
test_args=(-mod=readonly -count=1)
if [ "$mode" = race ]; then
    test_args+=(-race)
fi

for module in practice project; do
    printf '\nChecking %s (%s)\n' "$module" "$mode"
    (
        cd "$go_root/$module"
        unformatted=$(find . -type f -name '*.go' -exec gofmt -l {} +)
        if [ -n "$unformatted" ]; then
            printf 'Run gofmt on:\n%s\n' "$unformatted" >&2
            exit 1
        fi
        go vet -mod=readonly ./...
        go test "${test_args[@]}" ./...
        # -o keeps the project executable out of the learning folder.
        build_dir=$(mktemp -d)
        trap 'rm -rf "$build_dir"' EXIT
        go build -mod=readonly -o "$build_dir/" ./...
        go mod verify
    )
done
printf '\nBoth Go modules passed.\n'
