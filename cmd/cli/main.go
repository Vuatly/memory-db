package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"memory-db/internal/network"
	"os"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, logErr := zap.NewProduction()
	if logErr != nil {
		panic(logErr)
	}

	address := flag.String("address", "0.0.0.0:5444", "database server address")
	idleTimeout := flag.Duration("idle_timeout", 15*time.Second, "connection idle timeout")
	maxMessageSize := flag.Int("max_message_size", 4096, "max message size in bytes")
	flag.Parse()

	client, tcpErr := network.NewTCPClient(
		*address,
		*maxMessageSize,
		*idleTimeout,
	)
	if tcpErr != nil {
		logger.Error("failed to create client", zap.Error(tcpErr))
		return
	}

	reader := bufio.NewReader(os.Stdin)
	logger.Info("connected to server, ready to accept commands")
	logger.Info("press ctrl (control for mac) + c to exit from command line")

	for {
		fmt.Print("[cli] > ")
		query, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				logger.Fatal("connection was closed", zap.Error(err))
			}

			logger.Error("failed to read query", zap.Error(err))
		}

		response, err := client.Send([]byte(query))
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				logger.Fatal("connection was closed", zap.Error(err))
			}

			logger.Error("failed to send query", zap.Error(err))
		}

		fmt.Println(string(response))
	}
}
