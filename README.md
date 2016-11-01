# scmplat
采用beego框架golang开发，MVC架构的可视化后台进程创建，并发发布， 并发执行远程主机可执行脚本

db初始化：
create database SCMPLAT;grant all privileges on SCMPLAT.* to 'SCMPLAT'@'localhost' identified by '123456';grant all privileges on SCMPLAT.* to 'SCMPLAT'@'%' identified by '123456';grant all privileges on SCMPLAT.* to 'AIDBA'@'localhost' identified by '123456';grant all privileges on SCMPLAT.* to 'AIDBA'@'%' identified by '123456';flush privileges;
