# Jumper

Quickly jump to your project directories.

## Installation

```
brew install m-porter/tap/jumper
```

OR

```
go install github.com/m-porter/jumper@main
```

## Usage

The most effective way to use jumper is by making a bash function which combines
`jumper` with `cd`.

```shell
j() {
  local f="$(mktemp)"
  jumper to --out="$f"
  cd "$(cat "$f")" || return
}
```
