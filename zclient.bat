@echo off
cd clientview
rd /s /q dist
call npm run build
xcopy /D /I /F /Y /E "dist/*"  "../xclientapi/www/"
rd /s /q dist
cd ..
cd xclientapi

del clientapi
go env -w GOOS=linux
go build -o clientapi -ldflags "-s -w" main.go
go env -w GOOS=windows

@REM go build -o adminapi.exe main.go

xcopy /D /I /F /Y "clientapi"  "../"

del clientapi
cd ..
call ossutil rm oss://bblive/app/clientapi
call ossutil cp clientapi oss://bblive/app/

del clientapi


ssh root@47.238.161.17 "./client.sh"
