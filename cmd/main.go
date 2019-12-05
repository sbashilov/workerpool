package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const workerNum = 3

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	inputCh := make(chan func())
	wg := &sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		i := i
		wg.Add(1)
		go func() {
		LOOP:
			for {
				select {
				case val, ok := <-inputCh:
					fmt.Println(fmt.Sprintf("worker: %d, ok %v, work: %v", i, ok, val))
					time.Sleep(time.Second)
					val()
				case <-ctx.Done():
					for val := range inputCh {
						time.Sleep(time.Second)
						val()
					}
					wg.Done()
					break LOOP
				default:
				}
			}
		}()
	}
	for i := 0; i < 10; i++ {
		inputCh <- print(int64(i))
		if i == 5 {
			cancel()
			break
		}
	}
	close(inputCh)
	wg.Wait()
}

func print(i int64) func() {
	return func() {
		fmt.Println(fmt.Sprintf("%d - print", i))
	}
}
