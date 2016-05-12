#!/bin/bash -x

CONSUL_ARGS=""
if [ x$CONSUL_URL != x"" ] ; then
  CONSUL_ARGS="--consul-url $CONSUL_URL --consul-address $(hostname --ip-address)"
fi
./ccTray2Slack --log ${LOGFILE} --config /root/config/ --username ${USERNAME} --password ${PASSWORD} start --port ${PORT} ${CONSUL_ARGS} &
tail -F ${LOGFILE}
