package client

import (
	"context"
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"google.golang.org/grpc"
)

type Client struct {
	stream sequencerpb.Sequencer_GetOrderClient
	oRspC  chan *sequencerpb.OrderResponse
	oReqC  chan *sequencerpb.OrderRequest
}

func NewClient(IP string, originColor uint32) (*Client, error) {
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	pbClient := sequencerpb.NewSequencerClient(conn)
	stream, err := pbClient.GetOrder(context.Background())
	if err != nil {
		return nil, err
	}
	err = stream.Send(&sequencerpb.OrderRequest{OriginColor: originColor})
	if err != nil {
		return nil, err
	}
	client := new(Client)
	client.stream = stream
	client.oReqC = make(chan *sequencerpb.OrderRequest, 1024)
	client.oRspC = make(chan *sequencerpb.OrderResponse, 1024)
	go client.sendOReqs()
	go client.receiveORsps()
	return client, nil
}

func (c *Client) MakeOrderRequest(oReq *sequencerpb.OrderRequest) {
	c.oReqC <- oReq
}

func (c *Client) GetNextOrderResponse() *sequencerpb.OrderResponse {
	return <-c.oRspC
}

func (c *Client) sendOReqs() {
	for oReq := range c.oReqC {

		err := c.stream.Send(oReq)
		if err != nil {
			return
		}
	}
}

func (c *Client) receiveORsps() {
	for {
		rsp, err := c.stream.Recv()
		if err != nil {
			return
		}
		c.oRspC <- rsp
	}
}
