Output docs/assets/scanner-adapter-spike.gif

Require node

Set Shell "bash"
Set FontSize 24
Set Width 1280
Set Height 720
Set Theme "Dracula"
Set TypingSpeed 18ms
Set PlaybackSpeed 1.2

Type "npm run experiment:scanner-spike:summary" Enter
Sleep 2s

Type "sed -n '1,16p' experiments/scanner-adapter-spike/out/summary.md" Enter
Sleep 3s
