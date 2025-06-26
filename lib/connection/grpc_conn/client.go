package grpc_conn

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	logUtil "telemetry-task/lib/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	logger = logUtil.LoggerWithPrefix("GRPC_CONN")
)

var (
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "protocol.TelemetryService"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".1s",
			  "MaxBackoff": ".1s",
			  "BackoffMultiplier": 2,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
)

func GetClientConnection(addr, certPath string) (*grpc.ClientConn, error) {
	logger.Debug("create client for", "addr:", addr)

	if err := validateAddr(addr); err != nil {
		return nil, fmt.Errorf("invalid address: %s, err: %w", addr, err)
	}
	var opts []grpc.DialOption

	tlsCredentials, err := loadClientTLSCredentials(certPath)
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}
	if tlsCredentials != nil {
		logger.Info("using TLS credentials for gRPC connection")
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		logger.Warn("using insecure credentials for gRPC connection")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	opts = append(opts, grpc.WithDefaultServiceConfig(retryPolicy))

	return grpc.NewClient(addr, opts...)
}

func validateAddr(addr string) error {
	if addr == "" {
		return errors.New("address cannot be empty")
	}
	el := strings.Split(addr, ":")
	if len(el) != 2 {
		return errors.New("address must be in the format 'host:port'")
	}
	if el[0] == "" || el[1] == "" {
		return errors.New("host and port cannot be empty")
	}
	if _, err := strconv.Atoi(el[1]); err != nil {
		return errors.New("port must be a valid number")
	}
	return nil
}

func loadClientTLSCredentials(certPath string) (credentials.TransportCredentials, error) {
	if certPath == "" {
		logger.Warn("client certificate path is empty, using insecure credentials")
		return nil, nil
	}
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		logger.Warn("client certificate not found, using insecure credentials")
		return nil, nil
	}

	pemServerCA, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("read server CA certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("add server CA's certificate")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}
	logger.Info("using TLS credentials for gRPC client")
	return credentials.NewTLS(config), nil
}
