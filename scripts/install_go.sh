#!/bin/bash

set -e

wget --continue --quiet https://golang.org/dl/go1.16.4.linux-amd64.tar.gz

sudo tar -C /usr/local -xzf go1.16.4.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin

go get github.com/mattn/goreman

sudo sh -c  "echo 'export PATH=\$PATH:/usr/local/go/bin' >> /etc/profile"
