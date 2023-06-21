package sequencer

//
//import (
//	"context"
//	"flag"
//	"fmt"
//	"github.com/nathanieltornow/PMLog/order/sequencer/sequencerpb"
//	"github.com/sirupsen/logrus"
//	"github.com/stretchr/testify/require"
//	"google.golang.org/grpc"
//	"os"
//	"strings"
//	"testing"
//	"time"
//)
//
//var (
//	endpointList []string
//	color        uint32
//)
//
//func TestMain(m *testing.M) {
//	endpoints := flag.String("endpoints", "", "")
//	colorFlag := flag.Int("color", 0, "")
//	flag.Parse()
//	if *endpoints == "" {
//		logrus.Fatalln("no endpoints given")
//	}
//	split := strings.Split(*endpoints, ",")
//	endpointList = split
//	color = uint32(*colorFlag)
//	ret := m.Run()
//	os.Exit(ret)
//}
//
//func TestSerial(t *testing.T) {
//	stream, err := getClientStream(endpointList[0])
//	require.NoError(t, err)
//	waitC := make(chan bool)
//	go func() {
//		for {
//			oRsp, err := stream.Recv()
//			if err != nil {
//				close(waitC)
//				return
//			}
//			fmt.Println(oRsp)
//		}
//	}()
//	for i := uint64(0); i < 10; i++ {
//		time.Sleep(time.Second)
//		err := stream.Send(&sequencerpb.OrderRequest{Lsn: i, NumOfRecords: 1, Color: color, OriginColor: 4000})
//		require.NoError(t, err)
//	}
//	time.Sleep(time.Second)
//	err = stream.CloseSend()
//	require.NoError(t, err)
//	<-waitC
//}
//
////func checkOReqORsp(t *testing.T, oReqC chan *pb.OrderRequest, oRspC chan *pb.OrderResponse) {
////	type OReqORspPair struct {
////		oReq *pb.OrderRequest
////		oRsp *pb.OrderResponse
////	}
////	lsnToPair := make(map[uint64]*OReqORspPair)
////
////
////}
//
//func getClientStream(IP string) (sequencerpb.Sequencer_GetOrderClient, error) {
//	conn, err := grpc.Dial(IP, grpc.WithInsecure())
//	if err != nil {
//		return nil, err
//	}
//	client := sequencerpb.NewSequencerClient(conn)
//	stream, err := client.GetOrder(context.Background())
//	if err != nil {
//		return nil, err
//	}
//	return stream, nil
//}
