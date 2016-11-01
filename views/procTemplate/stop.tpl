#!/bin/sh

# 作者:此脚本由SCM生成 
# 创建时间:{{.CREATE_TIME}}
# 脚本目的:启动进程
# 修改原因:
# 修改时间:
# 修改作者:

PROCESS_NAME="{{.DEPLOY_STARTUP_CLASS}}"
PROCESS_PARM="{{.DEPLOY_APPFRAME_VALUE}}"
{{if eq .STANDARD "true"}}
BASE_APP_HOME="${HOME}"
{{else}}
BASE_APP_HOME="${HOME}/{{.DEPLOY_HOST}}"
{{end}}
CURRENT_APP_HOME="${BASE_APP_HOME}/app/{{.DEPLOY_TYPE}}/{{.DEPLOY_NAME}}"

CUR_USER=`whoami`
ps -ef | grep ${PROCESS_NAME} | grep ${CUR_USER} | grep ${PROCESS_PARM} | grep java | grep -v grep | awk '{print $2}' | while read pid
do
        kill ${pid} 2>&1 >/dev/null
        echo "Process Name :${PROCESS_NAME},Parameter:${PROCESS_PARM},PID:${pid} Stop"
        echo "[Audit trail]:Process ${PROCESS_NAME} Stop" >> ${CURRENT_APP_HOME}/log/{{.DEPLOY_SHELL_NAME}}-$(date "+%Y%m%d").log
done

