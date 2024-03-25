@echo off
del adminapi
cd xzadminview
call npm run build
xcopy /D /I /F /Y /E "dist/*"  "../xzadminapi/www/"
rd /s /q dist
cd ..
cd xzadminapi

del adminapi
go env -w GOOS=linux
go build -o adminapi -ldflags "-s -w" main.go
go env -w GOOS=windows

xcopy /D /I /F /Y "adminapi"  "../"

del adminapi

cd ..
call ossutil rm oss://bblive/app/adminapi
call ossutil cp adminapi oss://bblive/app/

del adminapi

ssh root@47.238.161.17 "./admin.sh"

pause