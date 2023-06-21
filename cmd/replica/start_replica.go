package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/shard"
	"github.com/nathanieltornow/PMLog/storage/mem_log_go"
	"github.com/sirupsen/logrus"
	"log"
)

var (
	IP          = flag.String("IP", "", "")
	orderIP     = flag.String("order", "", "")
	originColor = flag.Int("color", 10, "")
)

func main() {
	flag.Parse()
	replica, err := shard.NewReplica(uint32(*originColor), mem_log_go.Create)
	if err != nil {
		log.Fatalln(err)
	}
	if err := replica.Start(*IP, *orderIP); err != nil {
		logrus.Fatalln(err)
	}
}
