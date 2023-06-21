# PMLog

**Requirements**
- [Go](https://golang.org/) (version >= 1.15). If you're using Linux, you can use `./scripts/install_go.sh && source /etc/profile`.
- [Make](https://www.gnu.org/software/make/)
- [pmdk](https://github.com/pmem/pmdk) (version >= 1.19)
- [libpmemobj-cpp](https://github.com/pmem/libpmemobj-cpp)


### Installation on each node

Clone the repository on each node:
```
git clone git@github.com:nathanieltornow/PMLog.git
cd PMLog
```

In Linux you might also want to manually install `Go`.
```
./scripts/install_go.sh
source /etc/profile
source ~/.profile
```

In NixOS we provie you with the `default.nix` file.


### Setup benchmark

1. Start the sequencer on one node
```shell
go run cmd/start_sequencer/start_sequencer.go -root -IP 0.0.0.0:7000
```

To create a tree-hierachy of sequencers, use the following commands:
```go run cmd/start_sequencer/start_sequencer.go -root -IP <root_ip>:7000```
```go run cmd/start_sequencer/start_sequencer.go -parIP <root_ip>:7000```
```go run cmd/start_sequencer/start_sequencer.go -parIP <middle_seq_ip>:7000```


2. Startup the shard
```shell
# for every replica
go run cmd/replica/start_replica.go -order <sequencer-IP>:7000 -IP 0.0.0.0:4000
```

3. Start a benchmark at the load generator
   - Modify benchmark/benchmark.config.yaml, that the endpoints are all pointing to the replicas
   - Modify `clients` benchmark/benchmark.config.yaml, for the number of clients
   - Modify `appends` and `reads` for the right ratio. (Just the ratio appends/reads is important)
   - Start the benchmark, scale the throughput with increasing x
   ```shell
   go run benchmark/shared_log/shared_log.go -config benchmark/benchmark.config.yaml -threads x
   ```

