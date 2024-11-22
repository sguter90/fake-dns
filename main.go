package main

import (
	"flag"
	"github.com/miekg/dns"
	"log"
	"strconv"
)

func main() {
	configPath := flag.String("c", "./config.json", "Path to config file")
	flag.Parse()

	config, err := NewConfigFromPath(*configPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	srv := &dns.Server{Addr: ":" + strconv.FormatUint(config.Port, 10), Net: "udp"}
	srv.Handler = &Handler{
		c: *config,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
