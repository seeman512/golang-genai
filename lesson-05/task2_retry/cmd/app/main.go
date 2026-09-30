// Command app demonstrates retrying an operation after temporary failures.
package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"lesson05/task2_retry"
)

func main() {
	op := retry.NewFlakyOperation(2, "operation succeeded")

	result, err := retry.Do(op, 3, 100*time.Millisecond)
	if err != nil {
		if errors.Is(err, retry.ErrTemporary) {
			fmt.Fprintln(os.Stderr, "temporary operation failure after retries:", err)
		} else {
			fmt.Fprintln(os.Stderr, "operation failed:", err)
		}
		return
	}

	fmt.Println(result)
}
