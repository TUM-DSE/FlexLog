package test

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	clientID = flag.Int("id", 0, "")
	shardIPs = flag.String("IPs", "", "")
	record   = strings.Repeat("p", 4)
	cl       *client.Client
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *shardIPs == "" {
		return
	}
	ipList := strings.Split(*shardIPs, ",")
	newClient, err := client.NewClient(uint32(*clientID), ipList)
	if err != nil {
		logrus.Fatalln("failed to start client", err)
	}
	cl = newClient
	ret := m.Run()
	os.Exit(ret)
}

func TestAppendRead(t *testing.T) {
	numOfAppends := 5
	gsns := make([]uint64, numOfAppends)
	for i := 0; i < numOfAppends; i++ {
		start := time.Now()
		gsn := cl.Append(record, 1)
		fmt.Println(time.Since(start))
		fmt.Println(gsn)
		gsns[i] = gsn
		time.Sleep(time.Second)
	}
	for _, gsn := range gsns {
		rec := cl.Read(gsn, 1)
		fmt.Println(rec, gsn)
		require.Equal(t, rec, record)
	}
}
