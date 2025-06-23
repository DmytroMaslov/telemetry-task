package grpcconn

import (
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
	//TODO: add addr validation
	logger.Debug("create client for", "addr:", addr)
	var opts []grpc.DialOption
	opts = append(opts,

		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy))

	return grpc.NewClient(addr, opts...)
}
