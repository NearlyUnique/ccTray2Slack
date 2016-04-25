#!/bin/bash -x
./ccTray2Slack --log ${LOGFILE} --config /root/config/ --username ${USERNAME} --password ${PASSWORD} start --port ${PORT} &
tail -F ${LOGFILE}
