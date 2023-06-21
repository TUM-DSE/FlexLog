#!/bin/sh

sudo apt-get update
sudo apt-get install -y libpmem1 librpmem1 libpmemblk1 libpmemlog1 libpmemobj1 libpmempool1 libpmemobj-cpp-dev

sudo apt-get install -y libpmem-dev librpmem-dev libpmemblk-dev libpmemlog-dev libpmemobj-dev libpmempool-dev libpmempool-dev

sudo apt-get install -y libpmem1-debug librpmem1-debug libpmemblk1-debug libpmemlog1-debug libpmemobj1-debug libpmempool1-debug

sudo apt install -y libtbb-dev
sudo apt install -y pmdk-tools

read -p 'Log Size in disk: ' logsize

pmempool create --layout==logLayout obj PMLog --size=$logsize
