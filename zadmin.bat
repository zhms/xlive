@echo off
del adminapi
cd xadminview
call npm run build
xcopy /D /I /F /Y /E "dist/*"  "../xadminapi/www/"
rd /s /q dist
cd ..
cd xadminapi

del adminapi
go env -w GOOS=linux
go build -o adminapi -ldflags "-s -w" main.go
go env -w GOOS=windows


@REM go build -o adminapi.exe main.go

xcopy /D /I /F /Y "adminapi"  "../"

del adminapi