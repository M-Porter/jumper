# Jumper

Seamlessly jump between projects on your machine.

## Usage

The most effective way to use jumper is by making a bash function which combines
`jumper` with `cd`.

```shell
j() {
    cd "$(jumper to)"
}
```
