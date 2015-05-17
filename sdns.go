//
// SDNS is a recursive nameserver supporting pluggable resovers
// via arbitrary commands.
//
// Resolver commands accept a Question encoded as JSON via stdin,
// and write Answers to stdout encoded as JSON.
//
package sdns

import "encoding/json"
import "net"
import "fmt"
import "io"

// Question query.
type Question struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Class string `json:"class"`
}

// String representation.
func (q *Question) String() string {
	return fmt.Sprintf("name=%s type=%s class=%s", q.Name, q.Type, q.Class)
}

// Answer record(s).
type Answers []*Answer

// Validate the records.
func (a Answers) Validate() error {
	for i, rr := range a {
		err := rr.Validate()
		if err != nil {
			return fmt.Errorf("record[%v]: %s", i, err)
		}
	}
	return nil
}

// Answer record.
type Answer struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   uint32 `json:"ttl"`
}

// String representation.
func (a *Answer) String() string {
	return fmt.Sprintf("type=%s value=%s ttl=%d", a.Type, a.Value, a.TTL)
}

// Validate the record.
func (a *Answer) Validate() error {
	switch a.Type {
	case "A":
		ip := net.ParseIP(a.Value)
		if ip == nil {
			return fmt.Errorf("invalid A record")
		}
	}
	return nil
}

// IP address.
func (a *Answer) IP() net.IP {
	return net.ParseIP(a.Value)
}

// Read question from the given reader.
func Read(r io.Reader) (*Question, error) {
	q := new(Question)
	err := json.NewDecoder(r).Decode(q)
	return q, err
}

// Write answers to the given writer.
func Write(a Answers, w io.Writer) error {
	return json.NewEncoder(w).Encode(a)
}
