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
  local f
  local where
  f="$(mktemp)"
  jumper to "$1" --out="$f"
  where="$(cat "$f")"
  rm -f "$f"
  cd "$where" || return
}
```
