#!/usr/bin/env bash

set -eu
: "${GITHUB_TOKEN:?GITHUB_TOKEN not set in the environment.}"

TAG_VERSION="$1"
COMMIT_MESSAGE="$2"

git ls-remote --exit-code --tags origin "$TAG_VERSION"
if [[ "$?" == "0" ]]; then
    echo "Warning! Remote tag $TAG_VERSION already exists."
else
    git tag -a "$TAG_VERSION" -m "$COMMIT_MESSAGE"
    git push origin "$TAG_VERSION"
fi

goreleaser release --rm-dist

# can delete tags with
#   git tag -d "$TAG_VERSION"
#   git push --delete origin "$TAG_VERSION"
