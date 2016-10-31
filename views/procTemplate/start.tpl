#!/bin/ksh

# 作者:此脚本由SCM生成
# 创建时间:{{.CREATE_TIME}}
# 脚本目的:启动进程
# 修改原因:
# 修改时间:
# 修改作者:

# *************************************************************************
# JAVA_OPTIONS - java启动选项
# JAVA_VM      - jvm选项
# MEM_ARGS     - 内存参数
# *************************************************************************

#判断进程是否重复启动
./monitor_{{.DEPLOY_SHELL_NAME}}.sh | read PROCESS_ALIVE_STATUS
if [ "$PROCESS_ALIVE_STATUS" = "PROCESS_EXIST" ];
then
    echo "此进程已经启动了,不能重复启动"
    exit 0;
fi
#判断进程是否重复启动结束

{{if eq .STANDARD "true"}}
BASE_APP_HOME="${HOME}"
{{else}}
BASE_APP_HOME="${HOME}/{{.DEPLOY_HOST}}"
{{end}}
CURRENT_APP_HOME="${BASE_APP_HOME}/app/{{.DEPLOY_TYPE}}/{{.DEPLOY_NAME}}"
CURRENT_APP_JAR=""

. "${BASE_APP_HOME}/sbin/setEnv.sh"

CLASSPATH="${CURRENT_APP_HOME}/config:${BASE_APP_HOME}/configext:${CURRENT_APP_JAR}:${CLASSPATH}"
export CLASSPATH
echo "CLASSPATH=${CLASSPATH}"

MEM_ARGS="{{.DEPLOY_MEM_ARGS}}"
echo "\n"
echo "MEM_ARGS=${MEM_ARGS}"

echo "\n"
echo "JAVA_OPTIONS=${JAVA_OPTIONS}"

#启动命令行
${JAVA_HOME}/bin/java ${MEM_ARGS} -D{{.DEPLOY_APPFRAME_KEY}}={{.DEPLOY_APPFRAME_VALUE}} ${JAVA_OPTIONS} {{.DEPLOY_SHELL_PARAMS}} {{.DEPLOY_CONNECTION_VALUE}} {{.DEPLOY_STARTUP_CLASS}} {{.DEPLOY_PARAMETER_VALUE}} 2>&1 | ${BASE_APP_HOME}/sbin/cronolog -k 3 ${CURRENT_APP_HOME}/log/{{.DEPLOY_SHELL_NAME}}-%Y%m%d.log &
echo "[Audit trail]:Process {{.DEPLOY_APPFRAME_VALUE}} Start" >> ${CURRENT_APP_HOME}/log/{{.DEPLOY_SHELL_NAME}}-$(date "+%Y%m%d").log

echo "Process Start Successful"

