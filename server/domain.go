package server

import "github.com/tj/sdns/config"
import "github.com/miekg/dns"
import "encoding/json"
import "os/exec"
import "time"
import "log"
import "net"
import "fmt"

// Command result.
type result struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   uint32 `json:"ttl"`
}

// String representation.
func (r *result) String() string {
	return fmt.Sprintf("type=%s value=%q ttl=%d", r.Type, r.Value, r.TTL)
}

// Domain resolver.
type Domain struct {
	*config.Domain
}

// ServeDNS resolution.
func (d *Domain) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	for _, q := range r.Question {
		log.Printf("[%v] <== %s %s %v\n", r.Id,
			dns.ClassToString[q.Qclass],
			dns.TypeToString[q.Qtype],
			q.Name)
	}

	res := new(dns.Msg)
	res.SetReply(r)
	res.Authoritative = true

	if r.Question[0].Qtype == dns.TypeSOA {
		res.Ns = append(res.Ns, d.soa())
	}

	start := time.Now()
	result, err := d.exec()
	if err != nil {
		log.Printf("[error] executing command: %s", err)
		msg := new(dns.Msg)
		msg.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(msg)
		return
	}

	if result.TTL == 0 {
		result.TTL = 300
	}

	ip := net.ParseIP(result.Value)
	if ip == nil {
		log.Printf("[error] failed to parse A record %q", result.Value)
		msg := new(dns.Msg)
		msg.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(msg)
		return
	}

	a := &dns.A{
		Hdr: dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    result.TTL,
		},
		A: ip,
	}

	res.Answer = append(res.Answer, a)

	log.Printf("[%v] ==> %s", r.Id, time.Since(start))
	log.Printf("[%v] ----> %s\n", r.Id, result)

	err = w.WriteMsg(res)
	if err != nil {
		log.Printf("[error] [%v] failed to respond – %s", r.Id, err)
	}
}

// SOA record.
func (d *Domain) soa() *dns.SOA {
	return &dns.SOA{
		Hdr: dns.RR_Header{
			Name:   d.Name,
			Rrtype: dns.TypeSOA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		Ns:      "ns." + d.Name,
		Mbox:    "hostmaster." + d.Name,
		Serial:  uint32(time.Now().Unix()),
		Refresh: 3600,
		Retry:   900,
		Expire:  172800,
		Minttl:  0,
	}
}

func (d *Domain) exec() (*result, error) {
	cmd := exec.Command("sh", "-c", d.Command)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	res := new(result)
	err = json.Unmarshal(out, res)
	return res, err
}
