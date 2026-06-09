package flags

// CryptoKey represents path to RSA key file.
type CryptoKey string

// String returns key file path value.
func (c CryptoKey) String() string {
	return string(c)
}

// Set sets key file path value.
func (c *CryptoKey) Set(s string) error {
	*c = CryptoKey(s)
	return nil
}
