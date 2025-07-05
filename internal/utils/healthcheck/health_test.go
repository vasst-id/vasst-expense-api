package healthcheck

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
	// "github.com/julofinance/julo-go-library/redis"
)

func TestNew(t *testing.T) {

	t.Run("given component name is empty, when New(), then return panic", func(t *testing.T) {
		assert.Panics(t, func() {
			New(WithDB(&sql.DB{}, Config{}))
		})
	})

	t.Run("given checkFunc is nil, when New(), then return panic", func(t *testing.T) {
		assert.Panics(t, func() {
			New(WithComponent(Component{
				Name:      "grpc",
				CheckFunc: nil,
			}))
		})
	})

	t.Run("given component with same name already registered, when New(), then return panic", func(t *testing.T) {
		assert.Panics(t, func() {
			New(
				WithDB(&sql.DB{}, Config{Name: "registered"}),
				WithDB(&sql.DB{}, Config{Name: "registered"}),
			)
		})
	})

	t.Run("given timeout is not set, when New(), then component use default timeout", func(t *testing.T) {
		h := New(WithComponent(Component{
			Name:      "http",
			Timeout:   0,
			CheckFunc: func(ctx context.Context) error { return nil },
		}))

		assert.Equal(t, 3*time.Second, h.components["http"].Timeout)
	})

	t.Run("given WithLog option, when New(), then healthCheck should have logger", func(t *testing.T) {
		logger := logs.New(logs.Options{})
		h := New(WithLogger(logger))
		assert.Equal(t, logger, h.logger)
	})
}

func TestHealthService_Handler(t *testing.T) {
	logger := logs.New(logs.Options{})
	ctrl := gomock.NewController(t)
	// cache := redis.NewMockCache(ctrl)
	// cache.EXPECT().Ping().Return(errors.New("failed to connect to redis"))

	type want struct {
		status int
		health Health
	}

	testCases := []struct {
		name string
		args []Option
		want want
	}{
		{
			name: "given success, when Handler(), then return 200 OK",
			args: nil,
			want: want{
				status: http.StatusOK,
				health: Health{
					Status:     StatusOK,
					Failures:   nil,
					Components: nil,
				},
			},
		},
		{
			name: "given timeout, when Handler(), then return 503 Service Unavailable",
			args: []Option{WithDB(&sql.DB{}, Config{
				Name:    "DB",
				Timeout: 1 * time.Nanosecond,
			})},
			want: want{
				status: http.StatusServiceUnavailable,
				health: Health{
					Status:     StatusUnavailable,
					Failures:   map[string]string{"DB": "context deadline exceeded"},
					Components: []string{"DB"},
				},
			},
		},
		// {
		// 	name: "given error, when Handler(), then return 503 Service Unavailable",
		// 	args: []Option{
		// 		WithRedis(cache, Config{
		// 			Name:    "redis",
		// 			Timeout: 1 * time.Millisecond,
		// 		}),
		// 		WithLogger(logger),
		// 	},
		// 	want: want{
		// 		status: http.StatusServiceUnavailable,
		// 		health: Health{
		// 			Status:     StatusUnavailable,
		// 			Failures:   map[string]string{"redis": "failed to connect to redis"},
		// 			Components: []string{"redis"},
		// 		},
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := New(tc.args...)

			r := httptest.NewRequest(http.MethodGet, "/health-check", nil)
			w := httptest.NewRecorder()
			h.Handler().ServeHTTP(w, r)

			result := w.Result()
			body, _ := io.ReadAll(result.Body)
			var resp APIResponse
			_ = json.Unmarshal(body, &resp)

			var data Health
			b, _ := json.Marshal(resp.Data)
			_ = json.Unmarshal(b, &data)

			assert.Equal(t, tc.want.status, result.StatusCode)
			assert.Equal(t, tc.want.health, data)
		})
	}
}
