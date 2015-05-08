package server

import "github.com/tj/sdns/config"
import "github.com/miekg/dns"
import "log"

// Server.
type Server struct {
	udp *dns.Server
	tcp *dns.Server
	mux *dns.ServeMux
	*config.Config
}

// New server.
func New(config *config.Config) *Server {
	return &Server{
		Config: config,
	}
}

// Start server.
func (s *Server) Start() error {
	s.mux = dns.NewServeMux()

	s.udp = &dns.Server{
		Addr:    s.Bind,
		Net:     "udp",
		Handler: s.mux,
		UDPSize: 65535,
	}

	s.tcp = &dns.Server{
		Addr:    s.Bind,
		Net:     "tcp",
		Handler: s.mux,
	}

	for _, domain := range s.Domains {
		s.mux.Handle(dns.Fqdn(domain.Name), &Domain{domain})
	}

	s.mux.Handle(".", &RandomUpstream{s.Upstream})

	go func() {
		if err := s.udp.ListenAndServe(); err != nil {
			log.Fatalf("[error] failed to bind udp server: %v", err)
		}
	}()

	go func() {
		if err := s.tcp.ListenAndServe(); err != nil {
			log.Fatalf("[error] failed to bind tcp server: %v", err)
		}
	}()

	return nil
}

// Stop the server.
func (s *Server) Stop() error {
	err := s.tcp.Shutdown()
	if err != nil {
		return err
	}

	return s.udp.Shutdown()
}
