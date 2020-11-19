package main

import (
	"fmt"
	"sync"
)

var (
	counter = 0

	lock sync.Mutex

	synchronizedInt SynchronizedInt
)

type SynchronizedInt struct {
	value int
	lock  sync.Mutex
}

func (i *SynchronizedInt) Increase() {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.value++
}

func (i *SynchronizedInt) Decrease() {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.value--
}

func (i *SynchronizedInt) Value() int {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.value
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go updateCounter(&wg)
	}

	wg.Wait()
	fmt.Printf("final counter %d\n", counter)
	fmt.Printf("final synchronized counter %d\n", synchronizedInt.Value())
}

func updateCounter(wg *sync.WaitGroup) {
	lock.Lock()
	defer lock.Unlock()
	counter++

	synchronizedInt.Increase()
	wg.Done()
}
