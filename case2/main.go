package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var counter = 0

// decorator functuion to run each case
func run(f func()) {
	pc := reflect.ValueOf(f).Pointer()
	fn := runtime.FuncForPC(pc)
	name := strings.Split(fn.Name(), ".")[1]

	counter += 1
	fmt.Println(counter, ": ", name, "\n======================")
	f()
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println()
}

// case 1: try non buffered chan
func try_non_buffered_chan() {
	// non-buffered chan
	ch := make(chan int)

	// must be read in another goroutine, otherwise deadlock occurs
	// ch <- 1 // <--- dead lock
	// val, ok := <-ch
	// fmt.Println("read value from non buffered chan: ", val, ok)

	go func() {
		val, ok := <-ch
		fmt.Println("read value from non buffered chan: ", val, ok)
	}()

	fmt.Println("inpu chan 1")
	ch <- 1

	close(ch)
}

// case 2: try read only chan
func try_read_only_chan() {
	ch := make(chan int)

	// convert to read only chan
	readOnly := (<-chan int)(ch)

	go func() {
		ch <- 42
	}()

	value, ok := <-readOnly
	fmt.Println("read value from read only chan: ", value, ok)

	// compile error: cannot send to read only chan
	//readOnly <- 100
}

// case 3: try write only chan
func try_write_only_chan() {
	ch := make(chan int)
	// convert to read only chan
	writeOnly := (chan<- int)(ch)

	go func() {
		fmt.Println("read value from chan: ", <-ch)
	}()

	writeOnly <- 42

	// compile error: cannot read from write only chan
	// value, ok := <-writeOnly
}

// case 4: try buffered chan
func try_write_to_closed_chan() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic, reason: ", r)
		}
	}()

	// non-buffered chan
	ch := make(chan int)

	close(ch)

	fmt.Println("chan closed")

	// try to write to closed channel
	fmt.Println("will be panic next line..")
	ch <- 1
	fmt.Println("recovered from panic..")
}

// case 5: try buffered chan
func try_buffered_chan() {
	ch := make(chan int, 10)
	quit := make(chan bool)
	wg := sync.WaitGroup{}
	workloads := map[int]int{
		0: 0, 1: 0, 2: 0,
	}

	// 3 workers to finish these assignments
	for i := 0; i < 3; i++ {
		go func(workerId int) {
			for {
				fmt.Println("worker: ", workerId, " waiting for assignment")
				select {
				case assignmentId, ok := <-ch:
					if ok {
						workloads[workerId] += 1
						time.Sleep(1 * time.Second)
						fmt.Println("assignment ", assignmentId, " finished by worker ", workerId)
						wg.Done()
					}
				case _, ok := <-quit:
					if !ok {
						fmt.Println("worker ", workerId, "quit!")
						return
					}
				}
			}
		}(i)
	}

	// 10 assignments
	wg.Add(10)
	for assignmentId := 0; assignmentId < 10; assignmentId++ {
		ch <- assignmentId
	}

	wg.Wait()

	// all assignment has been done
	close(quit)

	fmt.Println("workloads: ", workloads)
}

func main() {
	run(try_non_buffered_chan)

	run(try_write_to_closed_chan)

	run(try_read_only_chan)

	run(try_write_only_chan)

	run(try_buffered_chan)

	fmt.Println("End Processing!!")
}
