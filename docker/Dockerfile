FROM centos
ENV version=v0.0.9
RUN yum install -y wget
RUN wget https://github.com/christer79/ccTray2Slack/releases/download/${version}/ccTray2Slack -nv
RUN chmod u+x ccTray2Slack
CMD ./ccTray2Slack -config /root/config/config.json -username ***REMOVED*** -password ***REMOVED***
