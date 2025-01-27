package main

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"strconv"
	"time"
)

type Handler struct {
	c Config
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	fmt.Println(time.Now().Format(time.DateTime) + " " + r.Question[0].Name + " (" + dns.TypeToString[r.Question[0].Qtype] + ")")
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		ip := "127.0.0.1"
		if nsIp, ok := h.c.Nameservers[domain]; ok {
			ip = nsIp
		}
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   net.ParseIP(ip),
		})
	case dns.TypeNS:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		for nsName, nsIp := range h.c.Nameservers {
			msg.Answer = append(msg.Answer, &dns.NS{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 3600},
				Ns:  nsName,
			})

			msg.Extra = append(msg.Extra, &dns.A{
				Hdr: dns.RR_Header{Name: nsName, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
				A:   net.ParseIP(nsIp),
			})
		}
	case dns.TypeSOA:
		serial, _ := strconv.ParseUint(time.Now().UTC().Format("20060102"), 10, 64)
		msg.Authoritative = true
		domain := msg.Question[0].Name
		nameserver := ""
		nameIp := ""
		for nsName, nsIp := range h.c.Nameservers {
			nameserver = nsName
			nameIp = nsIp
			break
		}
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

		msg.Extra = append(msg.Extra, &dns.A{
			Hdr: dns.RR_Header{Name: nameserver, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   net.ParseIP(nameIp),
		})
	}
	w.WriteMsg(&msg)
	fmt.Println(time.Now().Format(time.DateTime) + " ")
	//out, _ := json.MarshalIndent(msg, "", "    ")
	//fmt.Println(string(out))
}
