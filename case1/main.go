package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	parent := context.Background()

	ctx, cancel := context.WithCancel(parent)
	runTimes := 5

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("finished!")
				return
			default:
				fmt.Println(runTimes)
				time.Sleep(1 * time.Second)
				runTimes -= 1
			}

			if runTimes == 0 {
				cancel()
				wg.Done()
			}
		}

	}(ctx)

	wg.Wait()
	fmt.Println("\nEnd Processing!!")
}
