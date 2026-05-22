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
	Enable       bool
	SkipVerify   bool
	RootCertPath string
	CertPath     string
	KeyPath      string
}

var (
	// ErrTLSParseCACert returned when parsing the root CA certificate failed.
	ErrTLSParseCACert = errors.New("unable to parse TLS root CA cert")

	// ErrTLSCertKeyBothOrNone returned when one of either client certificate or client key are specified, but the other is not.
	ErrTLSCertKeyBothOrNone = errors.New("TLS cert and key file paths are mutually inclusive " +
		"(if one is specified, the other must be too)")
)

// GetDialerConfig reads and parses necessary certs/keys and returns them in form of std lib's *tls.Config.
func (tlsConfig *TLSConfig) GetDialerConfig() (*tls.Config, error) { // #nosec G402
	dialerConfig := &tls.Config{
		InsecureSkipVerify: tlsConfig.SkipVerify,
	}

	if tlsConfig.RootCertPath != "" {
		certPool, err := tlsConfig.readCACert()
		if err != nil {
			return nil, err
		}

		dialerConfig.RootCAs = certPool
	}

	if tlsConfig.CertPath != "" || tlsConfig.KeyPath != "" {
		certificate, err := tlsConfig.readClientKeyPair()
		if err != nil {
			return nil, err
		}

		dialerConfig.Certificates = []tls.Certificate{certificate}
	}

	return dialerConfig, nil
}

func (tlsConfig *TLSConfig) readCACert() (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	pemCert, err := os.ReadFile(tlsConfig.RootCertPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read TLS root CA cert file: %w", err)
	}

	if !certPool.AppendCertsFromPEM(pemCert) {
		return nil, ErrTLSParseCACert
	}

	return certPool, nil
}

func (tlsConfig *TLSConfig) readClientKeyPair() (tls.Certificate, error) {
	if tlsConfig.CertPath == "" || tlsConfig.KeyPath == "" {
		return tls.Certificate{}, ErrTLSCertKeyBothOrNone
	}

	pemCert, err := os.ReadFile(tlsConfig.CertPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to read TLS client cert file: %w", err)
	}

	pemKey, err := os.ReadFile(tlsConfig.KeyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to read TLS client key file: %w", err)
	}

	certificate, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to parse TLS client X509 key pair: %w", err)
	}

	return certificate, nil
}
