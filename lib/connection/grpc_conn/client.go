package grpcconn

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	logUtil "telemetry-task/lib/logger"

	"google.golang.org/grpc"
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

func GetClientConnection(addr string) (*grpc.ClientConn, error) {
	logger.Debug("create client for", "addr:", addr)

	if err := validateAddr(addr); err != nil {
		return nil, fmt.Errorf("invalid address: %s, err: %w", addr, err)
	}
	var opts []grpc.DialOption
	opts = append(opts,

		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy))

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
