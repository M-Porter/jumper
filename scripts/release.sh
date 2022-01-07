#!/usr/bin/env bash

set -u
: "${GITHUB_TOKEN:?GITHUB_TOKEN not set in the environment.}"

TAG_VERSION="$1"
COMMIT_MESSAGE="$2"

git tag -a "$TAG_VERSION" -m "$COMMIT_MESSAGE"
git push origin "$TAG_VERSION"

goreleaser release
