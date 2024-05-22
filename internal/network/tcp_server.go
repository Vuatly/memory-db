package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"memory-db/internal/utils"
	"net"
	"time"

	"go.uber.org/zap"
)

type RequestHandler = func(ctx context.Context, request []byte) []byte

type TCPServer struct {
	address        string
	idleTimeout    time.Duration
	maxMessageSize int
	handler        RequestHandler
	logger         *zap.Logger
	semaphore      *utils.Semaphore
}

func NewTCPServer(
	address string,
	maxMessageSize int,
	maxConnections int,
	handler RequestHandler,
	idleTimeout time.Duration,
	logger *zap.Logger,
) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	if handler == nil {
		return nil, errors.New("handler required")
	}

	if maxConnections <= 0 {
		return nil, errors.New("maxConnections must be greater than zero")
	}

	return &TCPServer{
		address:        address,
		idleTimeout:    idleTimeout,
		logger:         logger,
		handler:        handler,
		maxMessageSize: maxMessageSize,
		semaphore:      utils.NewSemaphore(maxConnections),
	}, nil
}

func (s *TCPServer) Serve(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.address)
	defer func() {
		deferErr := listener.Close()
		if deferErr != nil {
			s.logger.Warn("failed to close listener", zap.Error(deferErr))
		}
	}()

	if err != nil {
		return fmt.Errorf("failed to initialize tcp listener: %w", err)
	}

	s.logger.Info("tcp server started", zap.String("address", s.address))

	for {
		conn, listenerErr := listener.Accept()
		if listenerErr != nil {
			if errors.Is(listenerErr, net.ErrClosed) {
				return nil
			}
			s.logger.Warn(fmt.Sprintf("failed to accept tcp connection: %e", listenerErr))
			continue
		}

		s.semaphore.Acquire()
		go func() {
			defer s.semaphore.Release()
			s.handleConnection(ctx, conn)
		}()
	}
}

func (s *TCPServer) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		if deferErr := conn.Close(); deferErr != nil {
			s.logger.Warn("failed to close connection", zap.Error(deferErr))
		}
	}()

	request := make([]byte, s.maxMessageSize)

	for {
		if err := conn.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
			s.logger.Error("failed to set read deadline", zap.Error(err))
		}

		count, err := conn.Read(request)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				s.logger.Warn("failed to read", zap.Error(err))
			}

			break
		}

		response := s.handler(ctx, request[:count])
		if _, err = conn.Write(response); err != nil {
			s.logger.Warn("failed to write", zap.Error(err))
			break
		}
	}
}
