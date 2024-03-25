@echo off
mysqldump -h127.0.0.1 -uroot -pBo5zcd*2ozsvtjss0 -P3306 --routines  x_live > live.sql
echo "backup live.sql done"
pause