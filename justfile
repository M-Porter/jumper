format:
    gofmt -w .
    goimports -w .

release version:
    #!/usr/bin/env -S bash -x
    set -u
    : "${GITHUB_TOKEN:?GITHUB_TOKEN not set in the environment.}"
    git ls-remote --exit-code --tags origin "{{version}}"
    if [[ "$?" == "0" ]]; then
        echo "Warning! Remote tag {{version}} already exists."
    else
        git tag -a "{{version}}" -m "Release {{version}}"
        git push origin "{{version}}"
    fi
    goreleaser release --clean

build:
    go build

clean:
    go clean

run cmd:
    go run . {{cmd}}
