package healthcheck

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
)

// Component contains component to check
type Component struct {
	Name      string
	Timeout   time.Duration
	CheckFunc func(ctx context.Context) error
}

type healthService struct {
	mu sync.Mutex

	components map[string]Component
	logger     *logs.Logger
}

func New(opts ...Option) *healthService {
	h := &healthService{
		components: make(map[string]Component),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *healthService) register(c Component) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if c.Name == "" {
		panic("name is empty")
	}

	if c.CheckFunc == nil {
		panic("checkFunc is nil")
	}

	if _, ok := h.components[c.Name]; ok {
		panic("component already registered")
	}

	// set default timeout
	if c.Timeout == 0 {
		c.Timeout = 3 * time.Second
	}

	h.components[c.Name] = c
}

func (h *healthService) Handler() http.Handler {
	return http.HandlerFunc(h.HandlerFunc)
}

func (h *healthService) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	health := h.measure(r.Context())

	response := &APIResponse{
		Success: true,
		Error:   "",
		Data:    health,
	}

	code := http.StatusOK
	if health.Status != StatusOK {
		code = http.StatusServiceUnavailable
	}
	writeResponse(w, response, code)
}

func (h *healthService) componentNames() []string {
	var c []string
	for _, component := range h.components {
		c = append(c, component.Name)
	}
	return c
}

func (h *healthService) logError(err error, msg string) {
	if h.logger != nil {
		h.logger.Err(err).Msgf("healthCheck err: %s", msg)
	} else {
		log.Printf("healthCheck %s: %s \n", msg, err)
	}
}

func (h *healthService) measure(ctx context.Context) *Health {
	h.mu.Lock()
	defer h.mu.Unlock()

	failures := map[string]string{}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	for _, c := range h.components {

		wg.Add(1)

		go func(c Component) {
			defer func() {
				wg.Done()
			}()

			// buffered channel, for non blocking send
			errCh := make(chan error, 1)

			// set timeout to release resources of specific CheckFunc(ctx)
			ctx, cancel := context.WithTimeout(ctx, c.Timeout)
			defer cancel()

			// sender own goroutine, and close by sender
			go func() {
				errCh <- c.CheckFunc(ctx)
				defer close(errCh)
			}()

			select {
			case <-ctx.Done():
				// drain channel, ensure channel close/cleanup
				<-errCh

				mu.Lock()
				defer mu.Unlock()

				failures[c.Name] = ctx.Err().Error()
				h.logError(ctx.Err(), c.Name)

			case err := <-errCh:
				if err != nil {
					mu.Lock()
					defer mu.Unlock()

					failures[c.Name] = err.Error()
					h.logError(err, c.Name)
				}
			}
		}(c)
	}

	wg.Wait()

	status := StatusOK
	if len(failures) > 0 {
		status = StatusUnavailable
	}

	health := &Health{
		Status:     status,
		Failures:   failures,
		Components: h.componentNames(),
	}
	return health
}
