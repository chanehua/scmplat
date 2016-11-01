#!/bin/sh

# 作者:此脚本由SCM组生成 
# 创建时间:{{.CREATE_TIME}}
# 脚本目的:启动进程
# 修改原因:
# 修改时间:
# 修改作者:

{{if eq .STANDARD "true"}}
BASE_APP_HOME="${HOME}"
{{else}}
BASE_APP_HOME="${HOME}/{{.DEPLOY_HOST}}"
{{end}}
$BASE_APP_HOME/sbin/monitor_process.sh {{.DEPLOY_STARTUP_CLASS}} {{.DEPLOY_APPFRAME_VALUE}}