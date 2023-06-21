package client

import (
	"context"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
)

type Client struct {
	id  uint32
	ctr uint32

	numOfReplicas uint32

	appReqChs map[int]chan *shardpb.AppendRequest
	appRspCh  chan *shardpb.AppendResponse

	readReqCh chan *shardpb.ReadRequest
	readRspCh chan *shardpb.ReadResponse

	pbClients map[int]shardpb.ReplicaClient

	waitingAppends   map[uint64]chan uint64
	waitingAppendsMu sync.RWMutex

	waitingReads   map[uint64]chan string
	waitingReadsMu sync.RWMutex
}

func NewClient(id uint32, replicaIPs []string) (*Client, error) {
	c := new(Client)
	c.id = id
	c.numOfReplicas = uint32(len(replicaIPs))
	c.appReqChs = make(map[int]chan *shardpb.AppendRequest)
	c.appRspCh = make(chan *shardpb.AppendResponse, 2048)
	c.pbClients = make(map[int]shardpb.ReplicaClient)
	c.waitingAppends = make(map[uint64]chan uint64)
	c.waitingReads = make(map[uint64]chan string)
	c.readReqCh = make(chan *shardpb.ReadRequest, 2048)
	c.readRspCh = make(chan *shardpb.ReadResponse, 2048)

	for i, ip := range replicaIPs {
		conn, err := grpc.Dial(ip, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		pbClient := shardpb.NewReplicaClient(conn)
		c.pbClients[i] = pbClient

		stream, err := pbClient.Append(context.Background())
		if err != nil {
			return nil, err
		}
		ch := make(chan *shardpb.AppendRequest, 2048)
		c.appReqChs[i] = ch
		go sendAppendRequests(ch, stream)
		go c.receiveAppendResponses(stream)

		if i == (1 % len(replicaIPs)) {
			readStream, err := pbClient.Read(context.Background())
			if err != nil {
				return nil, err
			}
			go c.sendReadRequests(readStream)
			go c.receiveReadResponses(readStream)
		}
	}

	go c.handleReadResponses()
	go c.handleAppendResponses()
	return c, nil
}

func (c *Client) Append(record string, color uint32) uint64 {
	token := c.getNewToken()
	waitingGsn := make(chan uint64, 1)
	c.waitingAppendsMu.Lock()
	c.waitingAppends[token] = waitingGsn
	c.waitingAppendsMu.Unlock()

	responsible := int(c.id) % len(c.appReqChs)
	for i, ch := range c.appReqChs {
		ch <- &shardpb.AppendRequest{Token: token, Color: color, Record: record, Responsible: responsible == i}
	}
	defer func() {
		c.waitingAppendsMu.Lock()
		delete(c.waitingAppends, token)
		c.waitingAppendsMu.Unlock()
	}()
	return <-waitingGsn
}

func (c *Client) handleAppendResponses() {
	numOfAppends := make(map[uint64]uint32)
	for appRsp := range c.appRspCh {
		numOfAppends[appRsp.Token]++
		if numOfAppends[appRsp.Token] == c.numOfReplicas {
			c.waitingAppendsMu.RLock()
			c.waitingAppends[appRsp.Token] <- appRsp.Gsn
			c.waitingAppendsMu.RUnlock()
		}
	}
}

func sendAppendRequests(ch chan *shardpb.AppendRequest, stream shardpb.Replica_AppendClient) {
	for appReq := range ch {
		err := stream.Send(appReq)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
}

func (c *Client) receiveAppendResponses(stream shardpb.Replica_AppendClient) {
	for {
		appRsp, err := stream.Recv()
		if err != nil {
			logrus.Fatalln(err)
		}
		c.appRspCh <- appRsp
	}
}

func (c *Client) Read(gsn uint64, color uint32) string {

	token := c.getNewToken()
	waitingRecord := make(chan string, 1)
	c.waitingReadsMu.Lock()
	c.waitingReads[token] = waitingRecord
	c.waitingReadsMu.Unlock()
	c.readReqCh <- &shardpb.ReadRequest{Gsn: gsn, Color: color, Token: token}

	defer func() {
		c.waitingReadsMu.Lock()
		delete(c.waitingReads, token)
		c.waitingReadsMu.Unlock()
	}()
	return <-waitingRecord
}

func (c *Client) handleReadResponses() {
	for readRsp := range c.readRspCh {
		c.waitingReadsMu.RLock()
		ch, ok := c.waitingReads[readRsp.Token]
		c.waitingReadsMu.RUnlock()
		if !ok {
			logrus.Fatalln("failed to find waiting readrequest")
		}
		ch <- readRsp.Record
	}
}

func (c *Client) sendReadRequests(stream shardpb.Replica_ReadClient) {
	for readReq := range c.readReqCh {
		err := stream.Send(readReq)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
}

func (c *Client) receiveReadResponses(stream shardpb.Replica_ReadClient) {
	for {
		readRsp, err := stream.Recv()
		if err != nil {
			logrus.Fatalln(err)
		}
		c.readRspCh <- readRsp
	}
}

func (c *Client) Trim(gsn uint64, color uint32) {
	for _, pbCl := range c.pbClients {
		_, err := pbCl.Trim(context.Background(), &shardpb.TrimRequest{Color: color, Gsn: gsn})
		if err != nil {
			logrus.Fatalln(err)
		}
	}
}

func (c *Client) getNewToken() uint64 {
	ctr := atomic.AddUint32(&c.ctr, 1)
	return (uint64(c.id) << 32) + uint64(ctr)
}
