package client

// TLSConfig is an internal structure representing TLS-related settings
// configured on the provider level.
type TLSConfig struct {
	Enable       bool
	SkipVerify   bool
	RootCertPath string
	CertPath     string
	KeyPath      string
}
