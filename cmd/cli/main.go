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

	stdScanner := newScanner(os.Stdin, app.Logger)
	go stdScanner.run()

	for {
		select {
		case query := <-stdScanner.readStringChan():
			app.Logger.Info(app.Database.HandleQuery(query))
		case sig := <-stdScanner.readSignalsChan():
			stdScanner.close()
			app.Logger.Info("received signal, shutting down...", zap.String("signal", sig.String()))
			return
		}
	}
}

type scanner struct {
	input    *os.File
	to       chan string
	signals  chan os.Signal
	isClosed bool

	logger *zap.Logger
}

func newScanner(input *os.File, logger *zap.Logger) *scanner {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	return &scanner{
		input:   input,
		to:      make(chan string),
		signals: signals,
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

func (s *scanner) readStringChan() <-chan string {
	return s.to
}

func (s *scanner) readSignalsChan() <-chan os.Signal {
	return s.signals
}

func (s *scanner) close() {
	s.isClosed = true
	defer close(s.signals)
	defer close(s.to)

	err := s.input.Close()
	if err != nil {
		s.logger.Error("error closing input", zap.Error(err))
	}
}
