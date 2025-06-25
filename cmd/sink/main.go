package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	logUtil "telemetry-task/lib/logger"
	"telemetry-task/lib/services/sink"
	pb "telemetry-task/protocol/telemetry"
	"time"

	"google.golang.org/grpc"
)

var (
	logger = logUtil.LoggerWithPrefix("MAIN")
)

// TODO: add mTLS https://liambeeton.com/programming/secure-grpc-over-mtls-using-go
// https://dev.to/techschoolguru/how-to-secure-grpc-connection-with-ssl-tls-in-go-4ph

func main() {
	logger.Info("Starting sink server...")

	configPath := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	var cfg *sink.Config
	if *configPath == "" {
		cfg = sink.DefaultConfig()
	} else {
		var err error
		cfg, err = sink.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err.Error())
		}
	}

	sink, err := sink.NewSink(cfg)
	if err != nil {
		log.Fatalf("failed to create sink service: %v", err.Error())
	}

	if err = sink.Start(); err != nil {
		log.Fatalf("failed to start sink service: %v", err.Error())
	}
	lis, err := net.Listen("tcp", cfg.BindAddress)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTelemetryServiceServer(grpcServer, sink)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-exit
		timer := time.AfterFunc(10*time.Second, func() {
			log.Println("Server couldn't stop gracefully in time. Doing force stop.")
			grpcServer.Stop()
		})
		defer timer.Stop()
		grpcServer.GracefulStop()
		sink.Stop()
		logger.Info("Sink stopped gracefully.")
	}()

	logger.Info("Sink serve...")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}
