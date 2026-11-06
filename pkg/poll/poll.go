// Package poll waits for a remote state by asking for it at a fixed interval.
//
// The v6 API has no webhook: a build is finished when reading it says so. That
// is all this package knows. It has no notion of builds, versions or HTTP, so
// it can serve any wait; the caller supplies how to read the state, when the
// state is final and which errors are worth asking again.
package poll

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	// DefaultInterval is the time between two reads. It is fixed, with no
	// backoff: builds settle in seconds and the API's rate limit is far above
	// what one CLI run asks for.
	DefaultInterval = 5 * time.Second
	// DefaultTimeout bounds one wait, that is, one call to Until.
	DefaultTimeout = 3 * time.Minute
)

// ErrTimeout is what Until wraps when the timeout passes before done reports
// true.
var ErrTimeout = errors.New("timed out waiting")

// Options tunes a wait. The zero value is valid and means the defaults.
type Options struct {
	// Interval is the time between two reads; zero means DefaultInterval.
	Interval time.Duration
	// Timeout bounds the whole wait; zero means DefaultTimeout.
	Timeout time.Duration
	// Retryable decides whether a failed read is asked again. It returns
	// true to retry, with the delay the server asked for, or zero to retry on
	// the next interval. Nil means no error is retried.
	Retryable func(error) (retry bool, after time.Duration)
	// After replaces time.After, so tests run without sleeping.
	After func(time.Duration) <-chan time.Time
}

func (o Options) withDefaults() Options {
	if o.Interval <= 0 {
		o.Interval = DefaultInterval
	}
	if o.Timeout <= 0 {
		o.Timeout = DefaultTimeout
	}
	if o.Retryable == nil {
		o.Retryable = func(error) (bool, time.Duration) { return false, 0 }
	}
	if o.After == nil {
		o.After = time.After
	}
	return o
}

// Until reads with fetch, right away and then every Interval, until done
// reports true, ctx ends or the timeout passes.
//
// It always returns the last value fetch read successfully, so the caller can
// say which state the wait ended in. The error is:
//
//   - the error done returned, when done ended the wait with one;
//   - the error of a read that Retryable does not accept;
//   - ctx.Err() when ctx ends first;
//   - ErrTimeout, wrapped with the error of the last read if it failed, when
//     the timeout passes first.
//
// A Retry-After longer than Interval is honoured: the next read waits for it,
// still bounded by the timeout.
func Until[T any](ctx context.Context, opt Options,
	fetch func(context.Context) (T, error),
	done func(T) (bool, error),
) (T, error) {
	opt = opt.withDefaults()
	deadline := opt.After(opt.Timeout)

	var last T
	var lastErr error
	for {
		v, err := fetch(ctx)
		if ctxErr := ctx.Err(); ctxErr != nil {
			return last, ctxErr
		}

		wait := opt.Interval
		if err != nil {
			retry, after := opt.Retryable(err)
			if !retry {
				return last, err
			}
			lastErr = err
			wait = max(wait, after)
		} else {
			last, lastErr = v, nil
			finished, doneErr := done(v)
			if finished || doneErr != nil {
				return last, doneErr
			}
		}

		select {
		case <-ctx.Done():
			return last, ctx.Err()
		case <-deadline:
			if lastErr != nil {
				return last, fmt.Errorf("%w after %s: %w", ErrTimeout, opt.Timeout, lastErr)
			}
			return last, fmt.Errorf("%w after %s", ErrTimeout, opt.Timeout)
		case <-opt.After(wait):
		}
	}
}
