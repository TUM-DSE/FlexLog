#!/bin/bash

ORDER_IP=$1
COLOR=$2

export ORDER_IP=$ORDER_IP
export COLOR=$COLOR

$GOPATH/bin/goreman start