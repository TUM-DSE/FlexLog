package main

import (
	"flag"
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/nathanieltornow/PMLog/benchmark"
	seq_client "github.com/nathanieltornow/PMLog/sequencer/client"
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	configPath  = flag.String("config", "", "")
	resultC     chan *benchmarkResult
	color       = flag.Int("color", 0, "")
	threadsFlag = flag.Int("threads", 0, "")
	wait        = flag.Bool("wait", false, "")
)

type benchmarkResult struct {
	operations int
}

func main() {

	flag.Parse()
	if *configPath == "" {
		logrus.Fatalln("no config file")
	}

	config, err := benchmark.GetBenchConfig(*configPath)
	if err != nil {
		logrus.Fatalln(err)
	}
	threads := config.Threads
	if *threadsFlag != 0 {
		threads = *threadsFlag
	}

	f, err := os.OpenFile("results_sequencer.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()

	throughputResults := make([]float64, 0)
	latencyResults := make([]float64, 0)

	clients := make([]*seq_client.Client, 0)
	for i := 0; i < threads; i++ {
		cl, err := seq_client.NewClient(config.Endpoints[0], uint32(i+1333))
		if err != nil {
			logrus.Fatalln(err)
		}
		clients = append(clients, cl)
	}

	for t := 0; t < config.Times; t++ {
		resultC = make(chan *benchmarkResult, config.Clients)
		for i, cl := range clients {
			go executeBenchmark(cl, uint32(*color), uint32(i+1333), config.Runtime)
		}

		overallOperations := 0
		overallLatency := time.Duration(0)

		for i := 0; i < threads; i++ {
			res := <-resultC
			overallOperations += res.operations
			overallLatency += config.Runtime
		}

		throughput := float64(overallOperations) / config.Runtime.Seconds()
		latency := time.Duration(overallLatency.Nanoseconds() / int64(overallOperations))

		throughputResults = append(throughputResults, throughput)
		latencyResults = append(latencyResults, latency.Seconds())

		fmt.Printf("-----Latency: %v, \nThroughput (ops/s): %v\n", latency, throughput)
	}

	overallThroughput, err := stats.Mean(throughputResults)
	overallLatency, err := stats.Mean(latencyResults)
	if err != nil {
		logrus.Fatalln(err)
	}

	if _, err := f.WriteString(fmt.Sprintf("%v, %v, %v\n", threads, overallThroughput, overallLatency)); err != nil {
		logrus.Fatalln(err)
	}

}

func executeBenchmark(client *seq_client.Client, color, originColor uint32, duration time.Duration) {

	operations := 0
	op := 0
	defer func() {
		fmt.Println(operations, op)
		resultC <- &benchmarkResult{operations: operations}
	}()

	if *wait {
		<-time.After(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
	}

	stop := time.After(2 * time.Second)
load:
	for {
		select {
		case <-stop:
			break load
		default:
			client.MakeOrderRequest(&sequencerpb.OrderRequest{Color: color, OriginColor: originColor, NumOfRecords: 1, Tokens: []uint64{1}})
			_ = client.GetNextOrderResponse()
		}

	}

	waitC := make(chan bool, 1)
	stop = time.After(duration)
	go func() {
		for {
			select {
			case <-stop:
				waitC <- true
				return
			default:
				client.MakeOrderRequest(&sequencerpb.OrderRequest{Color: color, OriginColor: originColor, NumOfRecords: 1, Tokens: []uint64{12}})
				op++
			}
		}
	}()

	for {
		select {
		case <-waitC:
			return
		default:
			_ = client.GetNextOrderResponse()
			operations++
		}
	}
}
