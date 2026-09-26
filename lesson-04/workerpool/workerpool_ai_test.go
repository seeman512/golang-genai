package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestRunPool_AI_EnforcesTimeout verifies that RunPool passes a deadline to a
// job and returns a timeout result instead of waiting for the job indefinitely.
func TestRunPool_AI_EnforcesTimeout(t *testing.T) {
	const timeout = 50 * time.Millisecond

	jobs := make(chan Job, 1)
	jobs <- Job{
		ID: "timeout-job",
		Fetch: func(ctx context.Context) (int, error) {
			<-ctx.Done()
			return 0, ctx.Err()
		},
	}
	close(jobs)

	start := time.Now()
	results := RunPool(jobs, 1, timeout)

	var got Result
	select {
	case result, ok := <-results:
		if !ok {
			t.Fatal("RunPool closed the results channel without returning a result")
		}
		got = result
	case <-time.After(time.Second):
		t.Fatal("RunPool did not return a timeout result within one second")
	}

	if !errors.Is(got.Err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context deadline exceeded", got.Err)
	}
	if got.JobID != "timeout-job" {
		t.Fatalf("JobID = %q, want %q", got.JobID, "timeout-job")
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("RunPool took %v, want it to stop near the timeout", elapsed)
	}

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("RunPool returned more than one result for one job")
		}
	case <-time.After(time.Second):
		t.Fatal("RunPool did not close the results channel")
	}
}
