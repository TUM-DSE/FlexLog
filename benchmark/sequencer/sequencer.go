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
	defer func() {
		fmt.Println(operations)
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
	stop = time.After(duration)
	for {
		select {
		case <-stop:
			return
		default:
			client.MakeOrderRequest(&sequencerpb.OrderRequest{Color: color, OriginColor: originColor, NumOfRecords: 1, Tokens: []uint64{12}})
			_ = client.GetNextOrderResponse()
			operations++
		}
	}

}

//type syncClient struct {
//	stream sequencerpb.Sequencer_GetOrderClient
//	oRspC  chan *sequencerpb.OrderResponse
//	oReqC  chan *sequencerpb.OrderRequest
//
//	originColor uint32
//
//	mu                   sync.RWMutex
//	ctr                  uint64
//	waitingOrderRequests map[uint64]chan *sequencerpb.OrderResponse
//}

//func newSyncClient(IP string, originColor uint32) (*syncClient, error) {
//	conn, err := grpc.Dial(IP, grpc.WithInsecure())
//	if err != nil {
//		return nil, err
//	}
//	pbClient := sequencerpb.NewSequencerClient(conn)
//	stream, err := pbClient.GetOrder(context.Background())
//	if err != nil {
//		return nil, err
//	}
//	err = stream.Send(&sequencerpb.OrderRequest{OriginColor: originColor})
//	if err != nil {
//		return nil, err
//	}
//	sc := new(syncClient)
//	sc.stream = stream
//	sc.oReqC = make(chan *sequencerpb.OrderRequest, 1024)
//	sc.oRspC = make(chan *sequencerpb.OrderResponse, 1024)
//	sc.waitingOrderRequests = make(map[uint64]chan *sequencerpb.OrderResponse, 2048)
//	go sc.handleORsps()
//	go sc.sendOReqs()
//	go sc.receiveORsps()
//	sc.originColor = originColor
//	return sc, nil
//}
//
//func (sc *syncClient) getOrder(color uint32) *sequencerpb.OrderResponse {
//	ch := make(chan *sequencerpb.OrderResponse, 1)
//	sc.mu.Lock()
//	token := sc.ctr
//	sc.waitingOrderRequests[token] = ch
//	sc.ctr++
//	sc.mu.Unlock()
//	sc.oReqC <- &sequencerpb.OrderRequest{Color: color, OriginColor: sc.originColor, Tokens: []uint64{token}, NumOfRecords: 1}
//	defer func() {
//		sc.mu.Lock()
//		delete(sc.waitingOrderRequests, token)
//		sc.mu.Unlock()
//	}()
//	return <-ch
//}
//
//func (sc *syncClient) handleORsps() {
//	for oRsp := range sc.oRspC {
//		sc.mu.RLock()
//		ch, ok := sc.waitingOrderRequests[oRsp.Tokens[0]]
//		sc.mu.RUnlock()
//		if !ok {
//			logrus.Fatalln("not found")
//		}
//		ch <- oRsp
//	}
//}
//
//func (sc *syncClient) sendOReqs() {
//	for oReq := range sc.oReqC {
//		err := sc.stream.Send(oReq)
//		if err != nil {
//			return
//		}
//	}
//}
//
//func (sc *syncClient) receiveORsps() {
//	for {
//		rsp, err := sc.stream.Recv()
//		if err != nil {
//			return
//		}
//		sc.oRspC <- rsp
//	}
//}
