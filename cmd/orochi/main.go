package main

import (
	"flag"
	"log"

	"github.com/pankona/orochi"
)

func main() {
	p := flag.Int("port", 3000, "specify port number of http server")
	flag.Parse()

	o := &orochi.Server{
		PortList: []int{3000, 3001, 3002},
	}

	log.Fatal(o.Serve(*p))
}
