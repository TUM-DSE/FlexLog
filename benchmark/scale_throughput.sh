#!/bin/bash

set -e

for i in {1..25}
do
  go run shared_log/shared_log.go -config benchmark.config.yaml -threads ${i}
done
