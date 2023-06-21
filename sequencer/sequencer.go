package sequencer

import (
	"fmt"
	"github.com/nathanieltornow/PMLog/sequencer/client"
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

const (
	batchingInterval = 1 * time.Microsecond
)

type Sequencer struct {
	sequencerpb.UnimplementedSequencerServer

	id    uint32
	epoch uint32
	root  bool
	color uint32

	sn uint32

	parentClient *client.Client

	colorServices   map[uint32]*colorService
	colorServicesMu sync.RWMutex

	broadcastCh chan *sequencerpb.OrderResponse

	oRspCs   map[idColorTuple]chan *sequencerpb.OrderResponse
	oRspCsID uint32
	oRspCsMu sync.RWMutex

	oReqInCh chan *sequencerpb.OrderRequest
	oReqCh   chan *sequencerpb.OrderRequest
}

func NewSequencer(root bool, color uint32) *Sequencer {
	s := new(Sequencer)
	s.root = root
	s.color = color
	s.oRspCs = make(map[idColorTuple]chan *sequencerpb.OrderResponse)
	s.oReqCh = make(chan *sequencerpb.OrderRequest, 2048)
	s.oReqInCh = make(chan *sequencerpb.OrderRequest, 2048)
	s.colorServices = make(map[uint32]*colorService)
	s.broadcastCh = make(chan *sequencerpb.OrderResponse, 2048)
	return s
}

func (s *Sequencer) Start(IP string, parentIP string) error {
	if s.root {
		return s.startGRPCServer(IP)
	}

	cl, err := client.NewClient(parentIP, s.color)
	if err != nil {
		return fmt.Errorf("failed to connect to parent: %v", err)
	}
	s.parentClient = cl
	go s.handleOrderResponses()

	go s.forwardOrderRequests()

	return s.startGRPCServer(IP)
}

func (s *Sequencer) startGRPCServer(IP string) error {
	lis, err := net.Listen("tcp", IP)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	sequencerpb.RegisterSequencerServer(server, s)

	go s.handleOrderRequests()
	go s.broadcastOrderResponses()

	logrus.Infoln("starting sequencer on ", IP)
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start sequencer: %v", err)
	}
	return nil
}

func (s *Sequencer) GetOrder(stream sequencerpb.Sequencer_GetOrderServer) error {
	oRspC := make(chan *sequencerpb.OrderResponse, 256)

	first := true
	var ict idColorTuple

	go forwardOrderResponses(stream, oRspC)
	for {
		oReq, err := stream.Recv()
		if first {
			s.oRspCsMu.Lock()
			ict = idColorTuple{id: s.oRspCsID, color: oReq.OriginColor}
			s.oRspCs[ict] = oRspC
			s.oRspCsID++
			s.oRspCsMu.Unlock()
			first = false
			continue
		}
		if err != nil {
			s.oRspCsMu.Lock()
			delete(s.oRspCs, ict)
			s.oRspCsMu.Unlock()
			close(oRspC)
			return err
		}

		s.oReqInCh <- oReq
	}
}

type idColorTuple struct {
	id    uint32
	color uint32
}

func (s *Sequencer) getColorService(color uint32) *colorService {
	s.colorServicesMu.RLock()
	cs, ok := s.colorServices[color]
	s.colorServicesMu.RUnlock()
	if !ok {
		return s.addColorService(color)
	}
	return cs
}

func (s *Sequencer) addColorService(color uint32) *colorService {
	s.colorServicesMu.Lock()
	defer s.colorServicesMu.Unlock()
	cs, ok := s.colorServices[color]
	if !ok {
		cs = newColorService(color, s.color, batchingInterval, s.oReqCh)
		s.colorServices[color] = cs
	}
	return cs
}
