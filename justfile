# Set dotenv-load+export loads all env vars defined in `.env` and exports them to child processes.
set dotenv-load
set export

[private]
default:
    just --list

format:
    gofmt -w .
    goimports -w .

[doc("Release a new version of jumper using goreleaser.")]
release version:
    #!/usr/bin/env -S bash -x
    set -u
    [[ "{{version}}" == v* ]] || { echo "Error: version must start with 'v' (got '{{version}}')"; exit 1; }
    read -r -p "Release {{version}}? [y/N] " confirm
    [[ "$confirm" == [yY] ]] || exit 0
    : "${GITHUB_TOKEN:?GITHUB_TOKEN not set in the environment.}"
    git ls-remote --exit-code --tags origin "{{version}}"
    if [[ "$?" == "0" ]]; then
        echo "Warning! Remote tag {{version}} already exists."
    else
        git tag -a "{{version}}" -m "Release {{version}}"
        git push origin "{{version}}"
    fi
    goreleaser release --clean

[doc("Same as just relase but with op cli helper")]
op-release version:
    op run --env-file=.env -- just release {{version}}

[doc("Tests goreleaser release")]
test-release:
    goreleaser release --snapshot --clean

[doc("Tests goreleaser build")]
build-release:
    goreleaser build --clean --snapshot

build:
    go build

clean:
    go clean

run cmd:
    go run . {{cmd}}

test:
    go test ./...

[doc("Build the binary and run a VHS tape from examples/ to produce a screenshot/gif in temp/.")]
example name: build
    mkdir -p temp
    vhs examples/{{name}}.tape
