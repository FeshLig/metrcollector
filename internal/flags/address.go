package flags

import (
	"errors"
	"strconv"
	"strings"
)

// NetAddress represents network address in host:port format.
type NetAddress struct {
	Host string
	Port int
}

// String returns string representation of network address.
func (a NetAddress) String() string {
	return a.Host + ":" + strconv.FormatInt(int64(a.Port), 10)
}

// Set parses network address from string.
func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}
