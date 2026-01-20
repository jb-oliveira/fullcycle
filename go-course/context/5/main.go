package main

import (
	"context"
	"fmt"
	"time"
)

// longRunningTask simulates a task that executes multiple steps (iterations).
// It checks the context at each step to see if it should stop.
func longRunningTask(ctx context.Context, taskName string, maxIterations int) (string, error) {
	fmt.Printf("[%s] Task started. Will try %d iterations.\n", taskName, maxIterations)

	for i := 1; i <= maxIterations; i++ {
		// 1. The Critical Checkpoint:
		select {
		case <-ctx.Done():
			// The ctx.Done() channel was closed, which means the timeout expired
			// or the context was canceled.
			fmt.Printf("[%s] ❌ Interrupted! Context canceled after %d iterations.\n", taskName, i-1)
			return "", ctx.Err() // Returns the context error (deadline exceeded)
		default:
			// No cancellation occurred, continue the work.
		}

		// Simulation of part of the work (each step takes 100ms)
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("[%s] Step %d/%d completed.\n", taskName, i, maxIterations)
	}

	// If the loop finishes without the context being canceled
	fmt.Printf("[%s] ✅ Completed successfully after %d iterations.\n", taskName, maxIterations)
	return fmt.Sprintf("Task '%s' completed.", taskName), nil
}

func main() {

	// --- Scenario 1: Success (Long enough timeout) ---
	fmt.Println("--- Scenario 1: Success (1 second timeout) ---")

	// We define 1 second timeout. The task completes (5 steps * 100ms = 500ms).
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1()

	result1, err1 := longRunningTask(ctx1, "Task-A", 5)

	if err1 != nil {
		fmt.Printf("Task-A Result: Error: %v\n", err1)
	} else {
		fmt.Printf("Task-A Result: %s\n", result1)
	}

	fmt.Println("\n" + "--- Scenario 2: Timeout Expired (Short timeout) ---")

	// --- Scenario 2: Timeout Expired ---

	// We define only 250ms timeout.
	// The task would need 1 second (10 steps * 100ms).
	ctx2, cancel2 := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel2()

	result2, err2 := longRunningTask(ctx2, "Task-B", 10)

	if err2 != nil {
		fmt.Printf("Task-B Result: Error: %v\n", err2)
	} else {
		fmt.Printf("Task-B Result: %s\n", result2)
	}
}
