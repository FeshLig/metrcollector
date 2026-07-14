package flags

import (
	"fmt"
	"net"
)

// TrustedSubnet stores trusted subnet in CIDR notation.
type TrustedSubnet struct {
	Raw    string
	Subnet *net.IPNet
}

// Set implements flag.Value.
func (t *TrustedSubnet) Set(value string) error {
	t.Raw = value

	if value == "" {
		t.Subnet = nil
		return nil
	}

	_, subnet, err := net.ParseCIDR(value)
	if err != nil {
		return fmt.Errorf("parse trusted subnet: %w", err)
	}

	t.Subnet = subnet
	return nil
}

// String implements flag.Value.
func (t TrustedSubnet) String() string {
	return t.Raw
}

func (t TrustedSubnet) Contains(ip net.IP) bool {
	if t.Subnet == nil {
		return true
	}

	return t.Subnet.Contains(ip)
}
