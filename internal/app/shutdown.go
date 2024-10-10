package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type gracefulShutdown[T any] struct {
	operand   T
	close     func(ctx context.Context, operator T) error
	errLogger *log.Logger
}

func newGracefulShutdown[T any](
	operand T,
	closerFunc func(ctx context.Context, operand T) error,
	errorLogger *log.Logger,
) gracefulShutdown[T] {
	return gracefulShutdown[T]{
		operand:   operand,
		close:     closerFunc,
		errLogger: errorLogger,
	}
}

func (g gracefulShutdown[T]) listen(ctx context.Context) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	cancelSignal := <-ch
	err := g.close(ctx, g.operand)
	if err != nil {
		g.errLogger.Printf("graceful shutdown error: %v", err)
	}
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		g.errLogger.Printf("graceful finding process error: %v", err)
		return
	}
	_ = p.Signal(cancelSignal)
}
