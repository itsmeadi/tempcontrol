package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsmeadi/tempcontrol/iotControl"
	"go.uber.org/zap"
)

func main() {

	urlString := flag.String("url", "tcp://172.17.0.1:1883", "")
	clientIDString := flag.String("clientID", "cid-1", "")
	desiredTemp := flag.Float64("desiredTemp", 22, "")
	pmcEnable := flag.Bool("motion", false, "Enable Motion Sensor")
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger.Info("Initializing...", zap.String("url", *urlString))

	tempControl := iotControl.NewRoomControl(*clientIDString, *urlString, logger, *desiredTemp, *pmcEnable)
	tempControl.InitSubscribers(context.TODO())
	logger.Info("Ready to receive messages")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	logger.Info("Shutting Down...")
}
