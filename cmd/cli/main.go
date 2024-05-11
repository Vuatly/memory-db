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
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	scanner := newFileScanner(os.Stdin, readChan, app.Logger)
	go scanner.startScanLoop()

	for {
		select {
		case query := <-readChan:
			app.Logger.Info(app.Database.HandleQuery(query))
		case sig := <-signals:
			scanner.close()

			close(readChan)
			close(signals)

			app.Logger.Info("received signal, shutting down...", zap.String("signal", sig.String()))
			return
		}
	}
}

type fileScanner struct {
	input    *os.File
	to       chan<- string
	isClosed bool

	logger *zap.Logger
}

func newFileScanner(input *os.File, toChan chan<- string, logger *zap.Logger) *fileScanner {
	return &fileScanner{
		input:  input,
		to:     toChan,
		logger: logger,
	}
}

func (s *fileScanner) startScanLoop() {
	scanner := bufio.NewScanner(s.input)
	for {
		if s.isClosed {
			return
		}

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			s.logger.Error("error reading input", zap.Error(err))
			continue
		}

		s.to <- scanner.Text()
	}
}

func (s *fileScanner) close() {
	s.isClosed = true

	err := s.input.Close()
	if err != nil {
		s.logger.Error("error closing input", zap.Error(err))
	}
}
