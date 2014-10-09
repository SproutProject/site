#!/bin/sh

service redis-server start
service nginx start

/usr/sbin/sshd -D
