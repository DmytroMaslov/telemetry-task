package grpc_conn

import (
	"crypto/tls"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

func LoadServerTLSCredentials(certPath, keyPath string) (credentials.TransportCredentials, error) {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		logger.Warn("server certificate not found, using insecure credentials")
		return nil, nil
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		logger.Warn("server key not found, using insecure credentials")
		return nil, nil
	}
	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("load server TLS credentials: %w", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	logger.Info("using TLS credentials for gRPC server")
	return credentials.NewTLS(config), nil
}
