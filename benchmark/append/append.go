package main

import (
	"flag"
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/nathanieltornow/PMLog/benchmark"
	log_client "github.com/nathanieltornow/PMLog/client"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var (
	configPath  = flag.String("config", "", "")
	resultC     chan *benchmarkResult
	record      = strings.Repeat("r", 4000)
	threadsFlag = flag.Int("threads", 0, "")
	wait        = flag.Bool("wait", false, "")
)

type benchmarkResult struct {
	latency time.Duration
}

type overallResult struct {
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

	f, err := os.OpenFile(fmt.Sprintf("results_a%v-r%v.csv", config.Appends, config.Reads),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()

	clients := make([]*log_client.Client, 0)
	for i := 0; i < config.Clients; i++ {
		client, err := log_client.NewClient(uint32(i), config.Endpoints)
		if err != nil {
			logrus.Fatalln(err)
		}
		clients = append(clients, client)
	}

	for t := 0; t < config.Times; t++ {
		resultC = make(chan *benchmarkResult, threads*len(clients))
		for _, client := range clients {
			for i := 0; i < threads; i++ {
				go executeBenchmark(client, config.Runtime, config.Appends, config.Reads)
			}
		}
		latencies := make([]float64, 0)
		for i := 0; i < threads*len(clients); i++ {
			res := <-resultC
			latencies = append(latencies, res.latency.Seconds())
		}

		latency, _ := stats.Mean(latencies)

		fmt.Printf("-----\nAppend:\nLatency: %v\nThroughput (ops/s): %v\n",
			latency, threads*len(clients)*config.Appends)

		for _, client := range clients {
			client.Trim(0, 0)
		}
	}

	//if _, err := f.WriteString(fmt.Sprintf("%v, %v, %v, %v, %v, %v, %v, %v, %v\n", throughput, appendLatency, appendMedian, append99Perc, append95Perc, readLatency, readMedian, read99Perc, read95Perc)); err != nil {
	//	logrus.Fatalln(err)
	//}

}

func executeBenchmark(client *log_client.Client, runtime time.Duration, numAppends, numReads int) {

	if *wait {
		<-time.After(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
	}

	stop := time.After(5 * time.Second)
	loadLoop(client, stop)

	stop = time.After(runtime)

	start := time.Now()
	for i := 0; i < numAppends; i++ {
		_ = client.Append(record, 0)
	}
	resultC <- &benchmarkResult{
		time.Since(start),
	}
	stop = time.After(5 * time.Second)
	loadLoop(client, stop)
}

func loadLoop(client *log_client.Client, stop <-chan time.Time) {
	for {
		select {
		case <-stop:
			return
		default:
			_ = client.Append(record, 0)
		}
	}
}
