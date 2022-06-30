#!/usr/bin/env bash

TARGET_PATH=./bin
RUN_FILE_NAME=reskd_darwin_amd64

run() {
  echo running... ${TARGET_PATH}/${RUN_FILE_NAME} $2
  ${TARGET_PATH}/${RUN_FILE_NAME} $2
}

start() {
  echo "starting $2"
  nohup ${TARGET_PATH}/${RUN_FILE_NAME} $2 >/tmp/resk.log 2>&1 &
  echo "started ${RUN_FILE_NAME} $2"
}

stop() {
  # 一条命令解决问题
  # kill -9 $(ps -ef | grep reskd | grep -v "grep" | awk '{print $2}')
  PIDS=$(ps -ef | grep "${RUN_FILE_NAME}" | grep -v 'grep' | awk '{print $2}')
  echo "$PIDS"
  for pid in $PIDS; do
    echo kill $pid
    kill -15 $pid
  done
}

restart() {
  stop && start $@
}

rerun() {
  stop && run $@
}

#./run.sh run 1 3
action="$1"
if [ "${action}" == "" ]; then
  action="run"
fi

case "${action}" in
start)
  start "$@"
  ;;
stop)
  stop "$@"
  ;;
restart)
  reatart "$@"
  ;;
run)
  run "$@"
  ;;
rerun)
  rerun "$@"
  ;;
*)
  echo "Usage: $0 {start|stop|restart|rerun|run} {dev|test|prod|...}"
  echo "      eg: ${TARGET_PATH}/${RUN_FILE_NAME} run dev"
  exit 1
  ;;
esac

exit 0
