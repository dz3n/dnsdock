package main

import (
	"github.com/miekg/dns"
	"log"
)

type DNSServer struct {
	config *Config
	server *dns.Server
}

func NewDNSServer(c *Config) *DNSServer {
	s := &DNSServer{config: c}

	mux := dns.NewServeMux()
	mux.HandleFunc(".", s.forwardRequest)

	s.server = &dns.Server{Addr: c.dnsAddr, Net: "udp", Handler: mux}

	return s
}

func (s *DNSServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *DNSServer) Stop() {
	s.server.Shutdown()
}

func (s *DNSServer) forwardRequest(w dns.ResponseWriter, r *dns.Msg) {
	c := new(dns.Client)
	if in, _, err := c.Exchange(r, s.config.nameserver); err != nil {
		log.Print(err)
		w.WriteMsg(new(dns.Msg))
	} else {
		w.WriteMsg(in)
	}
}
