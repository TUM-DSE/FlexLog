package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/sequencer"
	"github.com/sirupsen/logrus"
)

var (
	IP    = flag.String("IP", ":7000", "The IP on which the sequencer listens")
	parIP = flag.String("parIP", "", "The IP of the parent-sequencer")
	root  = flag.Bool("root", false, "If the sequencer is the root-sequencer")
	color = flag.Int("color", 0, "The color which the sequencer represents")
)

func main() {
	flag.Parse()
	seq := sequencer.NewSequencer(*root, uint32(*color))
	err := seq.Start(*IP, *parIP)
	if err != nil {
		logrus.Fatalln(err)
	}
}
