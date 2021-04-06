cd %~dp0
rmdir /S /Q bin
mkdir bin
copy config.json bin/config.json
copy config.tmp.json bin/config.tmp.json
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build -o bin/autologin-linux-amd64
set GOARCH=386
set GOOS=linux
go build -o bin/autologin-linux-386
set GOARCH=amd64
set GOOS=windows
go build -o bin/autologin-windows-amd64.exe
set GOARCH=386
set GOOS=windows
go build -o bin/autologin-windows-386.exe
pause