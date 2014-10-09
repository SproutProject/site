FROM debian:testing

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y openssh-server redis-server nginx golang git node-less vim sudo

ADD ./sshd-config /etc/ssh/sshd_config
RUN useradd backend &&\
	useradd -m -s /bin/bash -p '$6$DlhXA0v/$.pmXpWRyed2hz4rbVKzg7eti5KwW7xLj9WCT007EMsqIZzYdo9Wo3UX0r6JreFHIISC41Q3BVHtsJbyWJo4bw1' docker &&\
	gpasswd -a docker sudo && chmod 770 /home/docker &&\
	mkdir /home/docker/.ssh && chmod 700 /home/docker/.ssh
ADD ./sshkey.pub /home/docker/.ssh/authorized_keys
RUN chmod 400 /home/docker/.ssh/authorized_keys && chown -R docker:docker /home/docker/.ssh

ENV GOPATH /srv/goroot
RUN mkdir /srv/goroot &&\
	go get github.com/zenazn/goji &&\
	go get github.com/garyburd/redigo/redis

ADD ./nginx-default /etc/nginx/sites-enabled/default
ADD ./startup.sh /startup.sh
RUN mkdir /var/run/sshd && chmod +x /startup.sh

CMD /startup.sh
