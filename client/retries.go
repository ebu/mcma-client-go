package mcmaclient

import (
	"net/http"
	"time"
)

type RetryOptions struct {
	ShouldRetry func(*http.Response, error) bool
	Intervals   []time.Duration
}

var DefaultShouldRetry = func(resp *http.Response, err error) bool {
	return err == nil && resp.StatusCode < 500 && resp.StatusCode != 429
}

var DefaultRetryIntervals = []time.Duration{
	250 * time.Millisecond,
	500 * time.Millisecond,
	1 * time.Second,
	2 * time.Second,
	5 * time.Second,
	15 * time.Second,
	30 * time.Second,
	45 * time.Second,
	1 * time.Minute,
}

var DefaultRetryOptions = RetryOptions{
	ShouldRetry: DefaultShouldRetry,
	Intervals:   DefaultRetryIntervals,
}

func ExecuteWithDefaultRetries(client *http.Client, req *http.Request) (bool, *http.Response, error) {
	return ExecuteWithRetries(client, req, DefaultRetryOptions)
}

func ExecuteWithRetries(client *http.Client, req *http.Request, opts RetryOptions) (bool, *http.Response, error) {
	res, err := client.Do(req)
	done := opts.ShouldRetry(res, err)
	if !done {
		for i := 0; i < len(opts.Intervals); i++ {
			time.Sleep(opts.Intervals[i])
			res, err = client.Do(req)
			done = opts.ShouldRetry(res, err)
		}
	}
	return done, res, err
}
