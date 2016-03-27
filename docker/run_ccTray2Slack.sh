#!/bin/bash
rm -f ccTray2Slack
wget https://github.com/christer79/ccTray2Slack/releases/download/${VERSION}/ccTray2Slack -nv
chmod u+x ccTray2Slack
./ccTray2Slack -config /root/config/config.json -username ${USERNAME} -password ${PASSWORD} start
