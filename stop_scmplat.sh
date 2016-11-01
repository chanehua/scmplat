#!/bin/sh

PROCESS_NAME="$HOME/scmplat"

CUR_USER=`whoami`
ps -ef | grep ${PROCESS_NAME} | grep ${CUR_USER} | grep -v grep | awk '{print $2}' | while read pid
do
        kill -9 ${pid} 2>&1 >/dev/null
        echo "进程名称:${PROCESS_NAME},PID:${pid} 成功停止"
done
