package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsmeadi/tempcontrol/tc"
	"go.uber.org/zap"
)

func main() {

	urlString := flag.String("url", "tcp://localhost:1883", "")
	clientIDString := flag.String("clientID", "cid-1", "")
	desiredTemp := flag.Float64("desiredTemp", 22, "")
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger.Info("Initializing...")

	tempControl := tc.NewTempControl(*clientIDString, *urlString, logger, *desiredTemp)
	tempControl.InitSubscribers(context.TODO())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	logger.Info("Shutting Down...")
}
