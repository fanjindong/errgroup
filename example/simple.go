package main

import (
	"context"
	"fmt"
	"github.com/fanjindong/errgroup"
	"time"
)

func printValue(x int) {
	time.Sleep(100 * time.Millisecond) // Time spent simulating real logic
	fmt.Println(x)
}

func main() {
	eg := errgroup.NewContinue(context.Background())
	eg.Go(func(ctx context.Context) error {
		printValue(1)
		return nil
	})
	eg.Go(func(ctx context.Context) error {
		printValue(2)
		return nil
	})
	eg.Go(func(ctx context.Context) error {
		printValue(3)
		return nil
	})
	eg.Wait()
	// output:
	// 1
	// 2
	// 3
	// Perform time-consuming: 100ms
}
