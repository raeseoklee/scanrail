Output examples/headers-demo/out/headers-demo.gif
Set Shell zsh
Set FontSize 18
Set Width 1200
Set Height 720
Set TypingSpeed 40ms

Type "npm install -g scanrail"
Enter
Sleep 1s
Type "node examples/headers-demo/server.mjs > /tmp/scanrail-headers-demo.log 2>&1 &"
Enter
Sleep 1s
Type "scanrail init --non-interactive --project-name headers-demo --target http://127.0.0.1:18080 --force"
Enter
Sleep 1s
Type "scanrail run --only headers"
Enter
Sleep 2s
Type "ls .scanrail/reports"
Enter
Sleep 1s
