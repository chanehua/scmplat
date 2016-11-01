# scmplat
SCM平台使用beego框架，采用MVC架构,golang作为后台开发可视化后台进程启停，监控脚本的生成，补丁包、平台包、后台进程的并发发布，并发执行环境的可执行脚本；

# scmplat部署
#### 1.安装mysql，需要将my.cnf文件中的AUTOCOMMIT=0设置为AUTOCOMMIT=1自动提交事务；
#### 2.下载工程；
#### 3.执行dbinit.sql初始化数据库；
#### 4.修改conf/app.conf文件中的：
mysql数据库信息：  
[db]  
dbhost = "127.0.0.1"  
dbport = "3306"  
dbuser = "SCMPLAT"  
passwd = "123456"  
dbname = "SCMPLAT"  
ldap服务器信息，主要修改用户名以及密码：  
[ldap]  
ldapdomain = "ldap.xxx.yyy.com"  
ldapport = "389"  
binduame = "xxx yyy"  
bindpwd = "123456"  
basedn = "ou=xx,ou=yy-users,dc=zz,dc=com"   
#### 5.执行start_scmplat.sh启动程序，stop_scmplat.sh停止程序