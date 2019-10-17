package main

import (
	"flag"
	"fmt"
	"github.com/mangatmodi/consensus/share"
)

func main() {
	c := flag.String("config", "", "-config=host1:port1,host2:port2")
	flag.Parse()
	p, err := share.Parse(*c)
	fmt.Printf("%v\n", err)
	fmt.Printf("%v", p[0])
}
