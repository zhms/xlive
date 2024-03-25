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

cd ..
call ossutil rm oss://bblive/app/adminapi
call ossutil cp adminapi oss://bblive/app/

del adminapi

ssh root@47.238.161.17 "./admin.sh"
