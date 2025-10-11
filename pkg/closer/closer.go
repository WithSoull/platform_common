package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

// shutdownTimeout default, can be made configurable
const shutdownTimeout = 5 * time.Second

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// Closer manages the graceful shutdown process of the application
type Closer struct {
	mu     sync.Mutex                    // Protection against race conditions when adding functions
	once   sync.Once                     // Guarantees single execution of CloseAll
	done   chan struct{}                 // Channel for completion notification
	funcs  []func(context.Context) error // Registered shutdown functions
	logger Logger                        // Logger instance being used
}

// Global instance for use throughout the application
var globalCloser = NewWithLogger(&logger.NoopLogger{})

// AddNamed adds a shutdown function with a dependency name for logging to the global closer
func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

// Add adds shutdown functions to the global closer
func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

// CloseAll initiates the shutdown process for all registered functions in the global closer
func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

// SetLogger allows setting a custom logger for the global closer
func SetLogger(l Logger) {
	globalCloser.SetLogger(l)
}

// Configure configures the global closer to handle system signals
func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

// New creates a new Closer instance with the default logger log.Default()
func New(signals ...os.Signal) *Closer {
	return NewWithLogger(logger.Logger(), signals...)
}

// NewWithLogger creates a new Closer instance with a specified logger.
// If signals are provided, Closer will start listening for them and call CloseAll upon receipt.
func NewWithLogger(logger Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

// SetLogger sets the logger for the Closer
func (c *Closer) SetLogger(l Logger) {
	c.logger = l
}

// handleSignals processes system signals and calls CloseAll with a fresh shutdown context
func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.logger.Info(context.Background(), "system signal received, starting graceful shutdown...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(context.Background(), "error closing resources: %v", zap.Error(err))
		}

	case <-c.done:
		// CloseAll was already called manually, just exit
	}
}

// AddNamed adds a shutdown function with a dependency name for logging
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(ctx, fmt.Sprintf("closing %s...", name))

		err := f(ctx)

		duration := time.Since(start)
		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("error closing %s: %v (took %s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("%s closed successfully in %s", name, duration))
		}
		return err
	})
}

// Add adds one or more shutdown functions
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

// CloseAll calls all registered shutdown functions.
// Returns the first error encountered, if any.
func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info(ctx, "no functions to close.")
			return
		}

		c.logger.Info(ctx, "starting graceful shutdown process...")

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		// Execute in reverse order of addition
		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()

				// Panic protection
				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recovered in closer")
						c.logger.Error(ctx, "panic in shutdown function", zap.Any("error", r))
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		// Close error channel when all functions complete
		go func() {
			wg.Wait()
			close(errCh)
		}()

		// Read errors or context cancellation
		for {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "context cancelled during shutdown", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info(ctx, "all resources closed successfully")
					return
				}
				c.logger.Error(ctx, "error during shutdown", zap.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}
