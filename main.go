package main

import (
	"flag"
	"github.com/kpango/glg"
	"github.com/mangatmodi/consensus/lib"
	"github.com/mangatmodi/consensus/share"
)

func main() {
	c := flag.String("config", "", "-config=host1:port1,host2:port2")
	port := flag.Int("port", int(share.DEFAULT_PORT), "-port=200")
	flag.Parse()

	share.Port = uint16(*port)
	p, err := share.Parse(*c)
	
	heartBeats := lib.NewHeartbeatHandler()
	heartBeats.SendHeartBeats(p)
	server := share.NewTcpServer([]share.Handler{heartBeats})
	err = server.Start(int(share.Port))
	if err != nil {
		glg.Fatal(err)
	}
}
