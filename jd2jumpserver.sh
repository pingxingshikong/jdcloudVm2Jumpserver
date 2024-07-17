#!/bin/bash

# 程序相关变量
PROGRAM_NAME="jdcloudVm2Jumpserver"
PROGRAM_PATH="./${PROGRAM_NAME}"
PID_FILE="./${PROGRAM_NAME}.pid"

# 启动程序
start() {
    echo "Entering start function"
    if [ -f "$PID_FILE" ]; then
        echo "Program is already running. PID=$(cat $PID_FILE)"
    else
        echo "Starting $PROGRAM_NAME..."
        nohup "$PROGRAM_PATH" > /dev/null 2>&1 &
        echo $! > "$PID_FILE"
        echo "$PROGRAM_NAME started with PID=$(cat $PID_FILE)"
    fi
}

# 停止程序
stop() {
    echo "Entering stop function"
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        echo "Stopping $PROGRAM_NAME with PID=$PID..."
        kill "$PID"
        rm "$PID_FILE"
        echo "$PROGRAM_NAME stopped."
    else
        echo "$PROGRAM_NAME is not running."
    fi
}

# 重新启动程序
restart() {
    echo "Entering restart function"
    stop
    start
}


# 检查程序状态
status() {
    echo "Entering status function"
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p $PID > /dev/null; then
            echo "$PROGRAM_NAME is running with PID=$PID."
        else
            echo "$PROGRAM_NAME is not running, but PID file exists."
        fi
    else
        echo "$PROGRAM_NAME is not running."
    fi
}

# 显示用法
usage() {
    echo "Usage: \$0 {start|stop|restart}"
    exit 1
}

# 检查参数并执行相应操作
echo "Argument passed to script: $@"
case "$@" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    restart)
        restart
        ;;
    *)
        usage
        ;;
esac
