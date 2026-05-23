package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
)

// TLSConfig is an internal structure representing TLS-related settings
// configured on the provider level.
type TLSConfig struct {
	*tls.Config
	IsEnabled bool
}

var (
	// ErrTLSParseCACert returned when parsing the root CA certificate failed.
	ErrTLSParseCACert = errors.New("unable to parse TLS root CA cert")

	// ErrTLSCertKeyBothOrNone returned when one of either client certificate or client key are specified, but the other is not.
	ErrTLSCertKeyBothOrNone = errors.New("TLS cert and key file paths are mutually inclusive " +
		"(if one is specified, the other must be too)")
)

// NewTLSConfig reads and parses necessary certs/keys and constructs new *TLSConfig.
func NewTLSConfig(
	isEnabled bool,
	skipVerify bool,
	caFile string,
	certFile string,
	keyFile string,
) (*TLSConfig, error) { // #nosec G402
	tlsConfig := &TLSConfig{
		Config: &tls.Config{
			InsecureSkipVerify: skipVerify,
		},
		IsEnabled: isEnabled,
	}

	if caFile != "" {
		certPool, err := tlsConfig.readCACert(caFile)
		if err != nil {
			return nil, err
		}

		tlsConfig.RootCAs = certPool
	}

	if certFile != "" || keyFile != "" {
		certificate, err := tlsConfig.readClientKeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates = []tls.Certificate{certificate}
	}

	return tlsConfig, nil
}

func (tlsConfig *TLSConfig) readCACert(caFile string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	pemCert, err := os.ReadFile(caFile) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("unable to read TLS root CA cert file: %w", err)
	}

	if !certPool.AppendCertsFromPEM(pemCert) {
		return nil, ErrTLSParseCACert
	}

	return certPool, nil
}

func (tlsConfig *TLSConfig) readClientKeyPair(certFile, keyFile string) (tls.Certificate, error) {
	if certFile == "" || keyFile == "" {
		return tls.Certificate{}, ErrTLSCertKeyBothOrNone
	}

	pemCert, err := os.ReadFile(certFile) //nolint:gosec
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to read TLS client cert file: %w", err)
	}

	pemKey, err := os.ReadFile(keyFile) //nolint:gosec
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to read TLS client key file: %w", err)
	}

	certificate, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to parse TLS client X509 key pair: %w", err)
	}

	return certificate, nil
}
