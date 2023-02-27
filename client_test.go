package moneybutton

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// newTestClient returns a client for mocking (using a custom HTTP interface)
func newTestClient(httpClient httpInterface) *Client {
	client := NewClient(nil, nil)
	client.httpClient = httpClient
	return client
}

// TestNewClient tests the method NewClient()
func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("valid new client", func(t *testing.T) {
		client := NewClient(nil, nil)
		assert.NotNil(t, client)
		assert.NotNil(t, client.httpClient)
		assert.NotNil(t, client.Options)
	})

	t.Run("custom http client", func(t *testing.T) {
		client := NewClient(nil, http.DefaultClient)
		assert.NotNil(t, client)
		assert.NotNil(t, client.httpClient)
		assert.Nil(t, client.Options)
	})
}

// ExampleNewClient example using NewClient()
func ExampleNewClient() {
	client := NewClient(nil, nil)

	fmt.Printf("created new client: %s", client.Options.UserAgent)
	// Output:created new client: go-moneybutton: v0.1.0
}

// BenchmarkNewClient benchmarks the method NewClient()
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewClient(nil, nil)
	}
}

// TestDefaultClientOptions tests setting DefaultClientOptions()
func TestDefaultClientOptions(t *testing.T) {
	t.Parallel()

	t.Run("default client options", func(t *testing.T) {
		options := DefaultClientOptions()
		assert.NotNil(t, options)
		assert.Equal(t, defaultUserAgent, options.UserAgent)
		assert.Equal(t, 2.0, options.BackOffExponentFactor)
		assert.Equal(t, 2*time.Millisecond, options.BackOffInitialTimeout)
		assert.Equal(t, 2*time.Millisecond, options.BackOffMaximumJitterInterval)
		assert.Equal(t, 10*time.Millisecond, options.BackOffMaxTimeout)
		assert.Equal(t, 20*time.Second, options.DialerKeepAlive)
		assert.Equal(t, 5*time.Second, options.DialerTimeout)
		assert.Equal(t, 2, options.RequestRetryCount)
		assert.Equal(t, 10*time.Second, options.RequestTimeout)
		assert.Equal(t, 3*time.Second, options.TransportExpectContinueTimeout)
		assert.Equal(t, 20*time.Second, options.TransportIdleTimeout)
		assert.Equal(t, 10, options.TransportMaxIdleConnections)
		assert.Equal(t, 5*time.Second, options.TransportTLSHandshakeTimeout)
	})

	t.Run("no retry", func(t *testing.T) {
		options := DefaultClientOptions()
		options.RequestRetryCount = 0
		client := NewClient(options, nil)
		assert.NotNil(t, client)
		assert.NotNil(t, client.Options)
	})
}

// ExampleDefaultClientOptions example using DefaultClientOptions()
func ExampleDefaultClientOptions() {
	options := DefaultClientOptions()
	options.UserAgent = "Custom UserAgent v1.0"
	client := NewClient(options, nil)

	fmt.Printf("created new client with user agent: %s", client.Options.UserAgent)
	// Output:created new client with user agent: Custom UserAgent v1.0
}

// BenchmarkDefaultClientOptions benchmarks the method DefaultClientOptions()
func BenchmarkDefaultClientOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultClientOptions()
	}
}
