#!/bin/sh

NAME='filterWord/filterWord'

if [ -f './log/process.pid' ];then
    PID=`cat ./log/process.pid`
fi

if [ ! -n "$PID" ];then
    PID=`ps -A | grep "$NAME" | grep -v grep | awk '{print $1}'`
fi

if [ -n "$PID" ];then
    CNT=`ps -ef | grep "$PID" | grep "$NAME"| grep -v grep | wc -l`
else
    CNT=0
fi

if [ $CNT = 0 ];then
   ./filterWord -s restart
fi
