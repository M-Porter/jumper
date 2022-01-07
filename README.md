# Jumper

Quickly jump to your project directories.

## Installation

```
go install github.com/m-porter/jumper@main
```

## Usage

The most effective way to use jumper is by making a bash function which combines
`jumper` with `cd`.

```shell
j() {
    cd "$(jumper to "${@}")" || return
}
```
