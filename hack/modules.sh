#!/usr/bin/env bash
#
# Manage the repository's Go modules without vendoring their dependencies.

set -e

SCRIPTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPTDIR/.." && pwd)"

module_paths=(api client . man)

tidy() (
    for module_path in "${module_paths[@]}"; do
        (
            set -x
            cd "$PROJECT_DIR/$module_path"
            go mod tidy
        )
    done
)

replace() (
	set -x
	go mod edit -replace=github.com/moby/moby/api=./api -replace=github.com/moby/moby/client=./client
	go mod edit -modfile client/go.mod -replace=github.com/moby/moby/api=../api
)

dropreplace() (
	set -x
	ref=$1
	if [ -z "$ref" ]; then
		# module@master not always results in the actual latest commit on the
		# master branch.
		# Use the actual branch from the upstream or origin remote.
		for r in "upstream" "origin"; do
			if git remote get-url "$r" > /dev/null 2>&1; then
				ref="$r/master"
				break
			fi
		done
		if [ -z "$ref" ]; then
			echo "No valid master ref found" >&2
			exit 1
		fi
		echo "Using $ref" >&2
	fi

	ref=$(git rev-parse "$ref")
	go mod edit -dropreplace=github.com/moby/moby/api -dropreplace=github.com/moby/moby/client
	go mod edit -modfile client/go.mod -dropreplace=github.com/moby/moby/api

	go mod edit -modfile client/go.mod -require="github.com/moby/moby/api@$ref"
	(
		cd client
		go mod tidy
	)

	go mod edit \
		-require="github.com/moby/moby/api@$ref" \
		-require="github.com/moby/moby/client@$ref"
	go mod tidy
)

help() {
    printf '%s:\n' "$(basename "$0")"
    echo "  tidy: synchronize go.mod and go.sum for all modules"
    echo "  replace: use local api and client modules"
    echo "  dropreplace [<gitref>]: remove local module replaces and pin a Git revision"
    echo "  help: show this help"
}

case "${1:-tidy}" in
    tidy) tidy ;;
    replace) replace ;;
    dropreplace) dropreplace "${2:-}" ;;
    help) help ;;
    *) help; exit 1 ;;
esac
