package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/miekg/dns"
)

func main() {
	configPath := flag.String("c", "./config.json", "Path to config file")
	flag.Parse()

	config, err := NewConfigFromPath(*configPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	listenAddress := config.Ip + ":" + strconv.FormatUint(config.Port, 10)
	handler := &Handler{c: *config}

	errCh := make(chan error, 2)

	for _, net := range []string{"udp", "tcp"} {
		srv := &dns.Server{Addr: listenAddress, Net: net, Handler: handler}
		go func(net string) {
			if err := srv.ListenAndServe(); err != nil {
				errCh <- fmt.Errorf("failed to set %s listener: %w", net, err)
			}
		}(net)
	}

	log.Fatal(<-errCh)
}
