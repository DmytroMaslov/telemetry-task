package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	conn "telemetry-task/lib/connection/grpc_conn"
	logUtil "telemetry-task/lib/logger"
	"telemetry-task/lib/services/sensor"
)

var (
	logger = logUtil.LoggerWithPrefix("MAIN")
)

func main() {
	logger.Info("Sensor starting...")

	sensorName := flag.String("name", "test", "name of sensor")
	connectionStr := flag.String("addr", "localhost:8080", "Address of server in format IP:PORT")
	flag.Parse()

	conn, err := conn.GetClientConnection(*connectionStr)
	if err != nil {
		log.Fatalf("failed to create client connection, err: %s", err.Error())
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			logger.Error("failed to close client connection", "err", err.Error())
		}
	}()

	sensor, err := sensor.NewSensor(conn, 100, *sensorName)
	if err != nil {
		log.Fatalf("failed to create sensor service, err: %s", err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())

	err = sensor.Start(ctx)
	if err != nil {
		log.Fatalf("failed to run sensor service, err: %s", err.Error())
	}
	logger.Info("Sensor running...")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-exit
	logger.Info("exit signal")
	cancel()
	err = sensor.Stop()
	if err != nil {
		log.Fatalf("Failed to stop sensor, err: %s", err)
	}
	logger.Info("Sensor stopped gracefully.")
}
