package sequencer

import (
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type colorService struct {
	ctr uint64

	oReqCh chan *sequencerpb.OrderRequest

	// tokens to all order-requests
	cache   map[uint64]*sequencerpb.OrderRequest
	cacheMu sync.Mutex
}

func newColorService(color, originColor uint32, interval time.Duration,
	outOReqCh chan *sequencerpb.OrderRequest) *colorService {

	cs := new(colorService)

	cs.oReqCh = make(chan *sequencerpb.OrderRequest, 2048)
	cs.cache = make(map[uint64]*sequencerpb.OrderRequest)

	go cs.batchOrderRequests(color, originColor, interval, outOReqCh)
	return cs
}

func (cs *colorService) insertOrderRequest(oReq *sequencerpb.OrderRequest) {
	cs.oReqCh <- oReq
}

func (cs *colorService) getOrderResponses(oRsp *sequencerpb.OrderResponse) []*sequencerpb.OrderResponse {
	res := make([]*sequencerpb.OrderResponse, 0)
	i := uint64(0)
	for _, token := range oRsp.Tokens {
		cs.cacheMu.Lock()
		oReq, ok := cs.cache[token]
		delete(cs.cache, token)
		cs.cacheMu.Unlock()
		if !ok {
			logrus.Fatalln("oReq not found")
		}
		res = append(res, &sequencerpb.OrderResponse{Tokens: oReq.Tokens, Color: oReq.Color, OriginColor: oReq.OriginColor, Gsn: oRsp.Gsn + i, NumOfRecords: oReq.NumOfRecords})
		i += uint64(oReq.NumOfRecords)
	}
	return res
}

func (cs *colorService) batchOrderRequests(color, originColor uint32, interval time.Duration, outOReqCh chan *sequencerpb.OrderRequest) {
	newBatch := true
	var send <-chan time.Time

	var token uint64
	var numOfRecords uint32
	var currentTokens []uint64

	for {
		select {
		case oReq := <-cs.oReqCh:
			if newBatch {
				numOfRecords = 0
				send = time.After(interval)
				currentTokens = make([]uint64, 0)
				newBatch = false
			}
			cs.cacheMu.Lock()
			cs.cache[token] = oReq
			cs.cacheMu.Unlock()
			currentTokens = append(currentTokens, token)
			token++
			numOfRecords += oReq.NumOfRecords
		case <-send:
			newBatch = true
			outOReqCh <- &sequencerpb.OrderRequest{Color: color, OriginColor: originColor, Tokens: currentTokens, NumOfRecords: numOfRecords}
		}
	}
}
