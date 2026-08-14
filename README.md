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

### fish

Add the following to your fish config.
```fish
jumper init fish | source
```

This loads a `j` function similar to the one described in the bash & zsh section.

If the default function of `j` conflicts with an existing function you have, you
can change it with the `--function` flag. For example, `jumper init fish --function="jump"`
would change the created function name to `jump`.

### bash/zsh
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
  cd "$(realpath "$where")" || return
}
```

### Disabling nerd font glyphs

If you do not have nerd fonts installed, you can disable those glyphs by adding
`no_nerd_font: true` to the jumper config (`jumper config edit`).
