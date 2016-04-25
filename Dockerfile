FROM centos
RUN yum install -y wget
ADD run_ccTray2Slack.sh /root/
ADD html /root/
ADD ccTray2Slack /root/
RUN chmod u+x /root/run_ccTray2Slack.sh
WORKDIR /root/
CMD ./run_ccTray2Slack.sh
