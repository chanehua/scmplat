# scmplat
SCM平台使用beego框架，采用MVC架构,golang作为后台开发可视化后台进程启停，监控脚本的生成，补丁包、平台包、后台进程的并发发布，并发执行环境的可执行脚本,一键并发安装、升级、卸载docker；

# scmplat部署
#### 1.安装mysql，需要将my.cnf文件中的AUTOCOMMIT=0设置为AUTOCOMMIT=1自动提交事务；
#### 2.下载工程：go get github.com/chanehua/scmplat；
#### 3.执行dbinit.sql初始化数据库；
#### 4.修改conf/app.conf文件中的：
mysql数据库信息：  
[db]  
dbhost = "127.0.0.1"  
dbport = "3306"  
dbuser = "SCMPLAT"  
passwd = "123456"  
dbname = "SCMPLAT"  
ldap服务器信息，主要修改用户名以及密码:  
[ldap]  
ldapdomain = "ldap.xxx.yyy.com"  
ldapport = "389"  
binduame = "xxx yyy"  
bindpwd = "123456"  
basedn = "ou=xx,ou=yy-users,dc=zz,dc=com"   
docker安装，升级，卸载配置(注意命令若是有顺序要求一定需要在此按顺序排列 ，以下配置是以centos为例):  
[系统_操作类型]  
[centos_install]  
install1 = "yum -y update"  
install2 = "yum -y install curl"  
install3 = "curl -fsSL https://get.docker.com/ | sh"  
install4 = "chkconfig docker on"  
install5 = "systemctl start docker"  
install6 = "systemctl status docker"  
[centos_upgrade]  
upgrade = "yum -y upgrade docker-engine"  
[centos_uninstall]  
uninstall1 = "yum list installed | grep docker | awk '{print $1}'"  

#### 5.执行start_scmplat.sh启动程序，stop_scmplat.sh停止程序
