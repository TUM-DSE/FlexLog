package shard

import (
	"github.com/sirupsen/logrus"
)

func (r *Replica) handleAppendResponses() {
	for appRsp := range r.appRespCh {
		r.clientsMu.RLock()
		ch, ok := r.clients[uint32(appRsp.Token>>32)]
		r.clientsMu.RUnlock()
		if !ok {
			logrus.Fatalln("couldn't find client")
		}
		ch <- appRsp
	}
}

func (r *Replica) handleOrderRequests() {
	for oReq := range r.oReqCh {
		r.orderClient.MakeOrderRequest(oReq)
	}
}

func (r *Replica) handleOrderResponses() {
	for {
		oRsp := r.orderClient.GetNextOrderResponse()
		r.getColorService(oRsp.Color).insertOrderResponse(oRsp)
	}
}

func (r *Replica) getColorService(color uint32) *colorService {
	r.colorServicesMu.RLock()
	cs, ok := r.colorServices[color]
	r.colorServicesMu.RUnlock()
	if !ok {
		return r.addColorService(color)
	}
	return cs
}

func (r *Replica) addColorService(color uint32) *colorService {
	r.colorServicesMu.Lock()
	defer r.colorServicesMu.Unlock()
	cs, ok := r.colorServices[color]
	if !ok {
		newLog, err := r.createLogFunc()
		if err != nil {
			logrus.Fatalln(err)
		}
		cs = newColorService(color, r.originColor, newLog, batchingInterval, r.oReqCh, r.appRespCh)
		r.colorServices[color] = cs
	}
	return cs
}
