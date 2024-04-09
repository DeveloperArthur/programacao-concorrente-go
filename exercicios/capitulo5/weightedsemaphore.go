/* Um semáforo ponderado é uma variação de um semáforo que permite adquirir e liberar
mais de uma licença ao mesmo tempo. As assinaturas de função para um semáforo
ponderado são as seguintes:
func (rw *WeightedSemaphore) Acquire(permits int)
func (rw *WeightedSemaphore) Release(permits int)
Use essas assinaturas de função para implementar um semáforo
ponderado com funcionalidade semelhante a um semáforo de contagem
Deve permitir que você adquira ou libere mais de uma licença */

package main

import (
	"sync"
)

type WeightedSemaphore struct {
	permits int
	cond    *sync.Cond
}

func NewSemaphore(n int) *WeightedSemaphore {
	return &WeightedSemaphore{
		permits: n,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (rw *WeightedSemaphore) Acquire(permits int) {
	rw.cond.L.Lock()
	for rw.permits-permits < 0 {
		rw.cond.Wait()
	}
	rw.permits -= permits
	rw.cond.L.Unlock()
}

func (rw *WeightedSemaphore) Release(permits int) {
	rw.cond.L.Lock()
	rw.permits += permits
	rw.cond.Signal()
	rw.cond.L.Unlock()
}
