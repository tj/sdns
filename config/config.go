package config

import "gopkg.in/yaml.v2"
import "io/ioutil"
import "io"
import "os"

// Config for bind address, domain, and upstream NS.
type Config struct {
	Bind     string
	Domains  []*Domain
	Upstream []string
}

// Domain represents a domain and its resolver command.
type Domain struct {
	Name    string
	Command string
}

// Read config from io.Reader.
func Read(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = yaml.Unmarshal(b, c)
	return c, err
}

// ReadFile reads config from `path`.
func ReadFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return Read(f)
}
