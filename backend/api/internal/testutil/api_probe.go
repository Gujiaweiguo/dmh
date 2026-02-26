// API readiness probe for integration tests
package testutil

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	// APIReadinessTimeout is the maximum time to wait for API to be ready
	APIReadinessTimeout = 10 * time.Second
	// APIProbeInterval is the interval between API probes
	APIProbeInterval = 500 * time.Millisecond
)

// APIProbeResult contains the result of an API probe
type APIProbeResult struct {
	Available bool
	Status    int
	Error     error
	Latency   time.Duration
}

// ProbeAPI checks if the API is available
func ProbeAPI(baseURL string) APIProbeResult {
	start := time.Now()
	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(baseURL + "/api/v1/auth/login")
	latency := time.Since(start)

	if err != nil {
		return APIProbeResult{
			Available: false,
			Error:     err,
			Latency:   latency,
		}
	}
	defer resp.Body.Close()

	return APIProbeResult{
		Available: true,
		Status:    resp.StatusCode,
		Latency:   latency,
	}
}

// WaitForAPI waits for the API to become ready
func WaitForAPI(ctx context.Context, baseURL string) error {
	ticker := time.NewTicker(APIProbeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("API readiness check timed out: %v", ctx.Err())
		case <-ticker.C:
			result := ProbeAPI(baseURL)
			if result.Available {
				return nil
			}
		}
	}
}

// SkipIfNoAPI skips the test if API is not available within timeout
func SkipIfNoAPI(t *testing.T) *IntegrationConfig {
	t.Helper()
	config := GetIntegrationConfig()

	ctx, cancel := context.WithTimeout(context.Background(), APIReadinessTimeout)
	defer cancel()

	if err := WaitForAPI(ctx, config.BaseURL); err != nil {
		t.Skipf("SKIP_REASON:%s | API at %s not ready within %v: %v",
			SkipReasonAPIUnavailable, config.BaseURL, APIReadinessTimeout, err)
		return nil
	}

	return config
}

// SkipfWithReason skips the test with a labeled reason
func SkipfWithReason(t *testing.T, reason SkipReason, format string, args ...interface{}) {
	t.Helper()
	t.Skipf("SKIP_REASON:%s | %s", reason, fmt.Sprintf(format, args...))
}
