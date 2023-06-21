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
	readLatency time.Duration
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
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	client, err := log_client.NewClient(90, config.Endpoints)
	if err != nil {
		logrus.Fatalln(err)
	}

	overallLatencies := make([]float64, 0)

	for t := 0; t < config.Times; t++ {
		resultC = make(chan *benchmarkResult, threads)
		toRead := config.Reads / threads
		for i := 0; i < threads; i++ {
			go executeBenchmark(client, toRead)
		}
		readLatencies := make([]float64, 0)
		for i := 0; i < threads; i++ {
			res := <-resultC
			readLatencies = append(readLatencies, res.readLatency.Seconds())
		}

		maxLatency, _ := stats.Max(readLatencies)
		overallLatencies = append(overallLatencies, maxLatency)

		fmt.Printf("-----Latency: %v\n", maxLatency)

		client.Trim(0, 0)
	}

	latency, _ := stats.Mean(overallLatencies)

	if _, err := f.WriteString(fmt.Sprintf("%v\n", latency)); err != nil {
		logrus.Fatalln(err)
	}

}

func executeBenchmark(client *log_client.Client, numReads int) {

	overallLatency := time.Duration(0)

	defer func() {
		resultC <- &benchmarkResult{
			readLatency: overallLatency,
		}
	}()

	if *wait {
		<-time.After(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
	}

	stop := time.After(5 * time.Second)
	loadLoop(client, stop)

	gsn := client.Append(record, 0)
	start := time.Now()
	for i := 0; i < numReads; i++ {
		client.Read(gsn, 0)
	}
	overallLatency = time.Since(start)
}

func loadLoop(client *log_client.Client, stop <-chan time.Time) {
	for {
		select {
		case <-stop:
			return
		default:
			gsn := client.Append(record, 0)
			_ = client.Read(gsn, 0)
		}
	}
}
