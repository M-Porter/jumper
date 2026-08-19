Output temp/to.gif

Set Shell "fish"
Set FontSize 16
Set Width 1000
Set Height 600
Set Padding 20
Set Theme "Catppuccin Latte"
Set TypingSpeed 50ms
Set FontFamily "SauceCodePro Nerd Font"

Env JUMPER_FIXTURES "true"

Type "./jumper to"
Enter
Sleep 500ms
Down
Sleep 100ms
Down
Sleep 100ms
Down
Sleep 100ms
Screenshot temp/to.png
Down
Sleep 100ms
Down
Sleep 100ms
Up
Sleep 100ms
Up
Sleep 250ms
Type "service"
Sleep 500ms
Backspace 7
Type "platformin"
Sleep 100ms
Backspace 2
Type "-infra"
Sleep 500ms
Enter
Sleep 1000ms
