package server

import "github.com/miekg/dns"
import "math/rand"
import "log"

// RandomUpstream resolves using a random NS in the set.
type RandomUpstream struct {
	upstream []string
}

// ServeDNS resolution.
func (h *RandomUpstream) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	ns := h.upstream[rand.Intn(len(h.upstream))]
	ns = defaultPort(ns)

	for _, q := range r.Question {
		log.Printf("[%v] <== %s %s %v (ns %s)\n", r.Id,
			dns.ClassToString[q.Qclass],
			dns.TypeToString[q.Qtype],
			q.Name,
			ns)
	}

	client := &dns.Client{
		Net: w.RemoteAddr().Network(),
	}

	res, rtt, err := client.Exchange(r, ns)
	if err != nil {
		msg := new(dns.Msg)
		msg.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(msg)
		return
	}

	log.Printf("[%v] ==> %s:", r.Id, rtt)
	for _, a := range res.Answer {
		log.Printf("[%v] ----> %s\n", r.Id, a)
	}

	err = w.WriteMsg(res)
	if err != nil {
		log.Printf("[error] [%v] failed to respond – %s", r.Id, err)
	}
}
