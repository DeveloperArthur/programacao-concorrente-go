/* é a ideia de uma barreira porque todos eles vão correr
talvez uns mais na frente, outros atras, mas todos vao parar
na barreira, esperar na barreira....
e só vao prosseguir quando todos tiverem chegado, ai a barreira
sai do caminho.
nesse cenario o red sempre termina primeiro mas fica esperando
o blue terminar. */

package main

import (
	"fmt"
	"sync"
	"time"
)

func workAndWait(name string, timeToWork int, barrier *Barrier) {
	start := time.Now()
	for {
		fmt.Println(time.Since(start), name, "is running")
		time.Sleep(time.Duration(timeToWork) * time.Second)
		fmt.Println(time.Since(start), name, "is waiting on barrier")
		barrier.Wait()
	}
}

func main() {
	barrier := NewBarrier(2)
	go workAndWait("Red", 4, barrier)
	go workAndWait("Blue", 10, barrier)
	time.Sleep(100 * time.Second)
}

type Barrier struct {
	size      int
	waitCount int
	cond      *sync.Cond
}

func NewBarrier(size int) *Barrier {
	condVar := sync.NewCond(&sync.Mutex{})
	return &Barrier{size, 0, condVar}
}

func (b *Barrier) Wait() {
	b.cond.L.Lock()
	b.waitCount += 1

	if b.waitCount == b.size {
		b.waitCount = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}

	b.cond.L.Unlock()
}
