#!/bin/bash -x
[ -f ccTray2Slack ] || wget https://github.com/christer79/ccTray2Slack/releases/download/${VERSION}/ccTray2Slack -nv
chmod u+x ccTray2Slack
echo ./ccTray2Slack -config /root/config/ -username ${USERNAME} -password ${PASSWORD} start -log ${LOGFILE}
./ccTray2Slack -config /root/config/ -username ${USERNAME} -password ${PASSWORD} start -log ${LOGFILE} &
tail -F ${LOGFILE}
