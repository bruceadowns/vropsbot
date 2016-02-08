FROM golang
MAINTAINER bdowns@vmware.com

ADD docker/ssh_config.txt /root/.ssh/config
ADD docker/ssh_known_hosts.txt /root/.ssh/known_hosts
ADD docker/ssh_id_rsa.txt /root/.ssh/id_rsa
ADD docker/ssh_id_rsa_pub.txt /root/.ssh/id_rsa.pub

RUN chmod 700 /root/.ssh/
RUN chmod 600 /root/.ssh/id_rsa

ADD config.json /root/
ENV VROPSBOT_CONFIG_JSON /root/config.json

RUN go get github.com/bruceadowns/vropsbot

ENTRYPOINT /go/bin/vropsbot
