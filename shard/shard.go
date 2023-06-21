package shard

import (
	"context"
	"fmt"
	"github.com/nathanieltornow/PMLog/sequencer/client"
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/nathanieltornow/PMLog/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

const (
	batchingInterval = 1 * time.Microsecond
)

type Replica struct {
	shardpb.UnimplementedReplicaServer

	originColor uint32

	appRespCh chan *shardpb.AppendResponse

	oReqCh chan *sequencerpb.OrderRequest

	createLogFunc storage.CreateLogFunc

	colorServices   map[uint32]*colorService
	colorServicesMu sync.RWMutex

	clients   map[uint32]chan *shardpb.AppendResponse
	clientsMu sync.RWMutex

	orderClient *client.Client
}

func NewReplica(originColor uint32, createLogFunc storage.CreateLogFunc) (*Replica, error) {
	r := new(Replica)
	r.originColor = originColor
	r.createLogFunc = createLogFunc

	r.appRespCh = make(chan *shardpb.AppendResponse, 4096)
	r.oReqCh = make(chan *sequencerpb.OrderRequest, 4096)
	r.colorServices = make(map[uint32]*colorService)
	r.clients = make(map[uint32]chan *shardpb.AppendResponse)

	return r, nil
}

func (r *Replica) Start(IP string, orderIP string) error {
	orderClient, err := client.NewClient(orderIP, r.originColor)
	if err != nil {
		return err
	}
	r.orderClient = orderClient
	r.addColorService(0)

	go r.handleOrderRequests()
	go r.handleOrderResponses()
	go r.handleAppendResponses()

	lis, err := net.Listen("tcp", IP)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	shardpb.RegisterReplicaServer(server, r)

	logrus.Infoln("Starting replica...")
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start replica: %v", err)
	}
	return nil
}

func (r *Replica) Append(stream shardpb.Replica_AppendServer) error {
	first := true
	appRspCh := make(chan *shardpb.AppendResponse, 1024)
	go replyAppendResponses(appRspCh, stream)

	for {
		appReq, err := stream.Recv()
		if err != nil {
			close(appRspCh)
			return err
		}
		if first {
			r.clientsMu.Lock()
			r.clients[uint32(appReq.Token>>32)] = appRspCh
			r.clientsMu.Unlock()
			first = false
			time.Sleep(time.Second)
		}
		r.getColorService(appReq.Color).insertAppendRequest(appReq)
	}
}

func replyAppendResponses(ch chan *shardpb.AppendResponse, stream shardpb.Replica_AppendServer) {
	for appRsp := range ch {
		err := stream.Send(appRsp)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
}

func (r *Replica) Read(stream shardpb.Replica_ReadServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		record, _ := r.getColorService(req.Color).read(req.Gsn)
		err = stream.Send(&shardpb.ReadResponse{Token: req.Token, Gsn: req.Gsn, Record: record})
		if err != nil {
			return err
		}
	}
}

func (r *Replica) Trim(_ context.Context, trimReq *shardpb.TrimRequest) (*shardpb.Ok, error) {
	return r.getColorService(trimReq.Color).trim(trimReq.Gsn)
}
