#!/bin/sh

SERVER_NAME=scmplat

nohup $HOME/scmplat > $HOME/logs/$SERVER_NAME-$(date +%Y%m%d).log 2>&1 &
