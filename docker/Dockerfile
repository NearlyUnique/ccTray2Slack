FROM centos
ENV version=v0.0.9
RUN yum install -y wget
ADD run_ccTray2Slack.sh .
RUN chmod u+x run_ccTray2Slack.sh
CMD ./run_ccTray2Slack.sh
