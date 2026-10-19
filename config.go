package main

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/netip"
	"os"
	"strings"
)

type Config struct {
	Ip          string            `json:"ip,omitempty"`
	Port        uint64            `json:"port"`
	Nameservers map[string]string `json:"nameservers"`

	// AllowedNetworks restricts access to the given networks (CIDR notation) or single IPs.
	// If empty or omitted, requests from all sources are allowed.
	AllowedNetworks []string `json:"allowed_networks,omitempty"`

	allowedPrefixes []netip.Prefix
}

func NewConfigFromPath(path string) (c *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		err = errors.New("Error when opening file from path " + path + ": " + err.Error())
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		err = errors.New("Error when reading file from path " + path + ": " + err.Error())
		return
	}
	err = json.Unmarshal(content, &c)
	if err != nil {
		err = errors.New("Error when unmarshalling JSON c file from path " + path + ": " + err.Error())
		return
	}
	err = c.parseAllowedNetworks()
	return
}

func (c *Config) parseAllowedNetworks() error {
	for _, entry := range c.AllowedNetworks {
		entry = strings.TrimSpace(entry)
		var prefix netip.Prefix
		var err error
		if strings.Contains(entry, "/") {
			prefix, err = netip.ParsePrefix(entry)
		} else {
			var addr netip.Addr
			addr, err = netip.ParseAddr(entry)
			if err == nil {
				prefix = netip.PrefixFrom(addr, addr.BitLen())
			}
		}
		if err != nil {
			return errors.New("Invalid entry in allowed_networks \"" + entry + "\": " + err.Error())
		}
		c.allowedPrefixes = append(c.allowedPrefixes, prefix.Masked())
	}
	return nil
}

// IsAllowed reports whether the given remote address may access the server.
// Without configured networks every address is allowed.
func (c *Config) IsAllowed(remote net.Addr) bool {
	if len(c.allowedPrefixes) == 0 {
		return true
	}
	var ip net.IP
	switch a := remote.(type) {
	case *net.UDPAddr:
		ip = a.IP
	case *net.TCPAddr:
		ip = a.IP
	default:
		return false
	}
	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return false
	}
	addr = addr.Unmap()
	for _, prefix := range c.allowedPrefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}
