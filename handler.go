package main

import (
	"fmt"
	"maps"
	"net"
	"slices"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

type Handler struct {
	c Config
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if !h.c.IsAllowed(w.RemoteAddr()) {
		// Drop the request silently; for TCP this also closes the connection.
		fmt.Println(time.Now().Format(time.DateTime) + " dropped request from " + w.RemoteAddr().String())
		w.Close()
		return
	}
	fmt.Println(time.Now().Format(time.DateTime) + " " + r.Question[0].Name + " (" + dns.TypeToString[r.Question[0].Qtype] + ")")
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA, dns.TypeAAAA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		ips := []string{"127.0.0.1", "::1"}
		if nsIps, ok := h.c.Nameservers[domain]; ok {
			ips = nsIps
		}
		for _, ip := range ips {
			record := addressRecord(domain, ip)
			// only return A records for A requests and AAAA records for AAAA requests
			if record.Header().Rrtype == r.Question[0].Qtype {
				msg.Answer = append(msg.Answer, record)
			}
		}
	case dns.TypeNS:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		for nsName, nsIps := range h.c.Nameservers {
			msg.Answer = append(msg.Answer, &dns.NS{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 3600},
				Ns:  nsName,
			})

			for _, ip := range nsIps {
				msg.Extra = append(msg.Extra, addressRecord(nsName, ip))
			}
		}
	case dns.TypeSOA:
		serial, _ := strconv.ParseUint(time.Now().UTC().Format("20060102"), 10, 64)
		msg.Authoritative = true
		domain := msg.Question[0].Name
		// map iteration order is random, so use the alphabetically first nameserver
		nameserver := ""
		if names := slices.Sorted(maps.Keys(h.c.Nameservers)); len(names) > 0 {
			nameserver = names[0]
		}
		nameIps := h.c.Nameservers[nameserver]
		msg.Answer = append(msg.Answer, &dns.SOA{
			Hdr:     dns.RR_Header{Name: domain, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 3600},
			Ns:      nameserver,
			Mbox:    "fake.example.com.",
			Serial:  uint32(serial),
			Refresh: 14400,
			Retry:   3600,
			Expire:  604800,
			Minttl:  3600,
		})

		for _, ip := range nameIps {
			msg.Extra = append(msg.Extra, addressRecord(nameserver, ip))
		}
	}
	w.WriteMsg(&msg)
	fmt.Println(time.Now().Format(time.DateTime) + " ")
	//out, _ := json.MarshalIndent(msg, "", "    ")
	//fmt.Println(string(out))
}

// addressRecord returns an A record for an IPv4 address and an AAAA record for an IPv6 address.
func addressRecord(name string, ip string) dns.RR {
	parsedIp := net.ParseIP(ip)
	if parsedIp.To4() != nil {
		return &dns.A{
			Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   parsedIp,
		}
	}
	return &dns.AAAA{
		Hdr:  dns.RR_Header{Name: name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 3600},
		AAAA: parsedIp,
	}
}
