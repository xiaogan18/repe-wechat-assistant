# build in windows
go build -o D:\golang\src\gitee.com\repe\bin\repe.exe ./main
xcopy web\tmpl bin\tmpl /S /D /Y
copy main\repe.yaml bin\repe.yaml /D /Y

#upload

#dba
## 表编码
-- 修改表编码方式
alter table user_info convert to character set utf8mb4 collate utf8mb4_general_ci; 
alter table room_info convert to character set utf8mb4 collate utf8mb4_general_ci; 