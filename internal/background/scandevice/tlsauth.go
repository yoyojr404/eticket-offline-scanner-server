package scandevice

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func loadCert(certFile, keyFile, caFile string, asTlsServer bool) (*tls.Config, error) {
	// load own certificate and key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("tls.LoadX509KeyPair: %w", err)
	}

	// load CA certificates
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: failed to read CA cert %w", err)
	}
	// Create a CA certificate pool and add the CA cert
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA cert to pool")
	}

	if asTlsServer {
		return &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caCertPool,
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
		}, nil
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}
