package poll

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// clock is a fake time.After. Until asks for the deadline first; every later
// call is a wait between reads, which fires at once and is recorded. The
// deadline fires after expireAfter such waits, or never when expireAfter is
// negative. Nothing sleeps.
type clock struct {
	expireAfter int
	waits       []time.Duration
	deadline    chan time.Time
}

func newClock(expireAfter int) *clock { return &clock{expireAfter: expireAfter} }

func (c *clock) After(d time.Duration) <-chan time.Time {
	if c.deadline == nil {
		c.deadline = make(chan time.Time)
		return c.deadline
	}
	c.waits = append(c.waits, d)
	if c.expireAfter >= 0 && len(c.waits) > c.expireAfter {
		close(c.deadline)
		return make(chan time.Time) // never fires, so the deadline wins
	}
	ch := make(chan time.Time, 1)
	ch <- time.Time{}
	return ch
}

// script returns a fetch that answers with each step in turn.
type step struct {
	v   string
	err error
}

func script(steps ...step) (func(context.Context) (string, error), *int) {
	calls := 0
	return func(context.Context) (string, error) {
		s := steps[min(calls, len(steps)-1)]
		calls++
		return s.v, s.err
	}, &calls
}

func doneAt(want string) func(string) (bool, error) {
	return func(v string) (bool, error) { return v == want, nil }
}

var (
	errRetry = errors.New("try again")
	errFatal = errors.New("bad request")
)

type waitAsked struct{ after time.Duration }

func (waitAsked) Error() string { return "rate limited" }

func retryable(err error) (bool, time.Duration) {
	var w waitAsked
	if errors.As(err, &w) {
		return true, w.after
	}
	return errors.Is(err, errRetry), 0
}

const timeout = 50 * time.Second

func opts(c *clock) Options {
	return Options{Interval: time.Second, Timeout: timeout, Retryable: retryable, After: c.After}
}

func TestUntilReturnsWhenDone(t *testing.T) {
	c := newClock(-1)
	fetch, calls := script(step{v: "queued"}, step{v: "building"}, step{v: "ready"})

	got, err := Until(context.Background(), opts(c), fetch, doneAt("ready"))

	require.NoError(t, err)
	assert.Equal(t, "ready", got)
	assert.Equal(t, 3, *calls)
	assert.Equal(t, []time.Duration{time.Second, time.Second}, c.waits)
}

func TestUntilReadsRightAway(t *testing.T) {
	c := newClock(-1)
	fetch, calls := script(step{v: "ready"})

	_, err := Until(context.Background(), opts(c), fetch, doneAt("ready"))

	require.NoError(t, err)
	assert.Equal(t, 1, *calls)
	assert.Empty(t, c.waits)
}

func TestUntilTimesOutWithTheLastState(t *testing.T) {
	c := newClock(2)
	fetch, _ := script(step{v: "queued"}, step{v: "building"})

	got, err := Until(context.Background(), opts(c), fetch, doneAt("ready"))

	assert.ErrorIs(t, err, ErrTimeout)
	assert.Contains(t, err.Error(), "50s")
	assert.Equal(t, "building", got)
}

func TestUntilTimeoutKeepsTheLastRetriedError(t *testing.T) {
	c := newClock(1)
	fetch, _ := script(step{v: "queued"}, step{err: errRetry})

	got, err := Until(context.Background(), opts(c), fetch, doneAt("ready"))

	assert.ErrorIs(t, err, ErrTimeout)
	assert.ErrorIs(t, err, errRetry)
	assert.Equal(t, "queued", got, "a failed read does not erase the last state seen")
}

func TestUntilStopsWhenContextEnds(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := newClock(-1)
	calls := 0
	fetch := func(context.Context) (string, error) {
		calls++
		if calls == 2 {
			cancel()
		}
		return "building", nil
	}

	got, err := Until(ctx, opts(c), fetch, doneAt("ready"))

	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, "building", got)
	assert.Equal(t, 2, calls)
}

func TestUntilStopsWhenContextEndsWhileWaiting(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	fetch, calls := script(step{v: "building"})
	opt := Options{Interval: time.Second, Timeout: timeout, After: func(time.Duration) <-chan time.Time {
		cancel()
		return make(chan time.Time)
	}}

	_, err := Until(ctx, opt, fetch, doneAt("ready"))

	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 1, *calls)
}

func TestUntilRetries(t *testing.T) {
	c := newClock(-1)
	fetch, calls := script(step{err: errRetry}, step{err: errRetry}, step{v: "ready"})

	got, err := Until(context.Background(), opts(c), fetch, doneAt("ready"))

	require.NoError(t, err)
	assert.Equal(t, "ready", got)
	assert.Equal(t, 3, *calls)
}

func TestUntilHonoursRetryAfter(t *testing.T) {
	c := newClock(-1)
	fetch, _ := script(
		step{err: waitAsked{after: 3 * time.Second}},
		step{err: waitAsked{}},                              // 429 without Retry-After: next interval
		step{err: waitAsked{after: 100 * time.Millisecond}}, // shorter than the interval
		step{v: "ready"},
	)

	_, err := Until(context.Background(), opts(c), fetch, doneAt("ready"))

	require.NoError(t, err)
	assert.Equal(t, []time.Duration{3 * time.Second, time.Second, time.Second}, c.waits)
}

func TestUntilStopsOnAnErrorItWillNotRetry(t *testing.T) {
	c := newClock(-1)
	fetch, calls := script(step{v: "queued"}, step{err: errFatal})

	got, err := Until(context.Background(), opts(c), fetch, doneAt("ready"))

	assert.ErrorIs(t, err, errFatal)
	assert.Equal(t, "queued", got)
	assert.Equal(t, 2, *calls)
}

func TestUntilWithoutRetryableRetriesNothing(t *testing.T) {
	c := newClock(-1)
	fetch, calls := script(step{err: errRetry}, step{v: "ready"})
	opt := opts(c)
	opt.Retryable = nil

	_, err := Until(context.Background(), opt, fetch, doneAt("ready"))

	assert.ErrorIs(t, err, errRetry)
	assert.Equal(t, 1, *calls)
}

func TestUntilReturnsTheErrorDoneGives(t *testing.T) {
	c := newClock(-1)
	fetch, _ := script(step{v: "building"}, step{v: "error"})
	failed := errors.New("build failed")
	done := func(v string) (bool, error) {
		if v == "error" {
			return true, failed
		}
		return false, nil
	}

	got, err := Until(context.Background(), opts(c), fetch, done)

	assert.ErrorIs(t, err, failed)
	assert.Equal(t, "error", got)
}

func TestUntilKeepsWaitingOnAnUnknownState(t *testing.T) {
	c := newClock(-1)
	fetch, calls := script(step{v: "queued"}, step{v: "promoting"}, step{v: "ready"})

	_, err := Until(context.Background(), opts(c), fetch, doneAt("ready"))

	require.NoError(t, err)
	assert.Equal(t, 3, *calls)
}

func TestDefaults(t *testing.T) {
	o := Options{}.withDefaults()
	assert.Equal(t, 5*time.Second, o.Interval)
	assert.Equal(t, 3*time.Minute, o.Timeout)
	retry, _ := o.Retryable(errRetry)
	assert.False(t, retry)
	assert.NotNil(t, o.After)
}
