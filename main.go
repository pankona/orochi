package main

import (
	"flag"
	"log"
)

var portList = []int{3000, 3001, 3002}

func main() {
	p := flag.Int("port", 3000, "specify port number of http server")
	flag.Parse()

	o := &Orochi{}
	log.Fatal(o.Serve(*p))
}
