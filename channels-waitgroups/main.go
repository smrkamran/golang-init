package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func doWork(d time.Duration, resCh chan string) {
	fmt.Println("Doing work...")
	time.Sleep(d)
	fmt.Println("Work done.")
	resCh <- fmt.Sprintf("Work %d", rand.Intn(100))
	wg.Done()
}

var wg *sync.WaitGroup

func main() {
	start := time.Now()
	resultCh := make(chan string)
	wg = &sync.WaitGroup{}
	wg.Add(2)

	go doWork(time.Second*2, resultCh)
	go doWork(time.Second*4, resultCh)

	go func() {
		for res := range resultCh {
			fmt.Println(res)
		}
		fmt.Printf("Work took %v seconds\n", time.Since(start))
	}()

	wg.Wait()
	close(resultCh)
	time.Sleep(time.Second)

}
