#!/usr/bin/env bash

set -u
: "${GITHUB_TOKEN:?GITHUB_TOKEN not set in the environment.}"

TAG_VERSION="$1"
COMMIT_MESSAGE="$2"

git ls-remote --exit-code --tags origin "$TAG_VERSION"
if [[ "$?" == "0" ]]; then
    echo "ERROR! Remote tag $TAG_VERSION already exists."
    echo "Exiting."
    exit 1
fi

git tag -a "$TAG_VERSION" -m "$COMMIT_MESSAGE"
git push origin "$TAG_VERSION"

goreleaser release

# can delete tags with
#   git tag -d "$TAG_VERSION"
#   git push --delete origin "$TAG_VERSION"
