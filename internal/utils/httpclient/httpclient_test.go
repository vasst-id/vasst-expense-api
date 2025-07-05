package httpclient

import (
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupHTTPClient(t *testing.T) Client {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	httpClientCfg := &Config{
		Timeout: 3000,
	}
	httpClient := New(httpClientCfg)
	return httpClient
}

func TestHttp(t *testing.T) {
	t.Parallel()

	httpClient := setupHTTPClient(t)

	t.Run("success doing Do for httpClient", func(t *testing.T) {
		req, err := http.NewRequest("GET", "https://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		assert.NotNil(t, resp)
		assert.Nil(t, err)

	})
}
