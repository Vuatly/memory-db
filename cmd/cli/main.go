package main

import (
	"bufio"
	"memory-db/cmd"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	app, err := cmd.NewApp()
	if err != nil {
		panic(err)
	}

	readChan := make(chan string)
	signalsChan := make(chan os.Signal)

	signal.Notify(signalsChan, syscall.SIGINT, syscall.SIGTERM)

	stdScanner := newScanner(os.Stdin, readChan, signalsChan, app.Logger)
	go stdScanner.run()

	for {
		select {
		case query := <-readChan:
			app.Logger.Info(app.Database.HandleQuery(query))
		case sig := <-signalsChan:
			stdScanner.close()

			close(readChan)
			close(signalsChan)

			app.Logger.Info("received signal, shutting down...", zap.String("signal", sig.String()))
			return
		}
	}
}

type scanner struct {
	input    *os.File
	to       chan<- string
	signals  chan<- os.Signal
	isClosed bool

	logger *zap.Logger
}

func newScanner(input *os.File, toChan chan<- string, signalsChan chan<- os.Signal, logger *zap.Logger) *scanner {
	return &scanner{
		input:   input,
		to:      toChan,
		signals: signalsChan,
		logger:  logger,
	}
}

func (s *scanner) run() {
	stdScanner := bufio.NewScanner(s.input)
	for {
		if s.isClosed {
			return
		}

		isScanned := stdScanner.Scan()
		err := stdScanner.Err()

		if err == nil && !isScanned {
			s.logger.Info("received EOF, exit...")
			s.signals <- syscall.SIGQUIT
			return
		}

		if err != nil {
			s.logger.Error("error reading input", zap.Error(err))
			continue
		}

		s.to <- stdScanner.Text()
	}
}

func (s *scanner) close() {
	s.isClosed = true

	err := s.input.Close()
	if err != nil {
		s.logger.Error("error closing input", zap.Error(err))
	}
}
