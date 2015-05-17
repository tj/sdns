package server

import "github.com/tj/sdns/config"
import "github.com/miekg/dns"
import "github.com/tj/sdns"
import "encoding/json"
import "strings"
import "os/exec"
import "bytes"
import "time"
import "log"
import "fmt"

// Domain resolver.
type Domain struct {
	*config.Domain
}

// Strip suffix from the domain, for example "api-02.ec2." becomes "api-02".
func (d *Domain) strip(name string) string {
	return strings.Replace(name, "."+d.Name, "", 1)
}

// ServeDNS resolution.
func (d *Domain) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	for _, q := range r.Question {
		log.Printf("[info] [%v] <-- %s %s %v\n", r.Id,
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
	answers, err := d.resolve(&r.Question[0])
	if err != nil {
		log.Printf("[error] [%v] executing command: %s", r.Id, err)
		msg := new(dns.Msg)
		msg.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(msg)
		return
	}

	err = answers.Validate()
	if err != nil {
		log.Printf("[error] [%v] invalid answers: %s", r.Id, err)
		msg := new(dns.Msg)
		msg.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(msg)
		return
	}

	for _, answer := range answers {
		switch answer.Type {
		case "A":
			rr := &dns.A{
				Hdr: dns.RR_Header{
					Name:   r.Question[0].Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    answer.TTL,
				},
				A: answer.IP(),
			}
			res.Answer = append(res.Answer, rr)
		}
	}

	log.Printf("[info] [%v] --> %s %s", r.Id, answers, time.Since(start))

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

// Resolve query via command.
func (d *Domain) resolve(q *dns.Question) (sdns.Answers, error) {
	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	query := &sdns.Question{
		Name:  d.strip(q.Name),
		Type:  dns.TypeToString[q.Qtype],
		Class: dns.ClassToString[q.Qclass],
	}

	err := json.NewEncoder(stdin).Encode(query)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("sh", "-c", d.Command)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, stderr.String())
	}

	var answers sdns.Answers
	err = json.NewDecoder(stdout).Decode(&answers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %s", stdout.String())
	}

	return answers, nil
}
