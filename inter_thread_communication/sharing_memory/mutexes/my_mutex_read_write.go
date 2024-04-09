package main

import "sync"

type ReadWriteMutex struct {
	//Variável inteira para contar o número de goroutines de leitura atualmente na seção crítica
	readersCounter int
	//Mutex para sincronizar o acesso de leitores
	readersLock sync.Mutex
	//Mutex para bloquear qualquer acesso de escritores
	globalLock sync.Mutex
}

// A função desse método é bloquear o mutex de escrita
// Quando o mutex de leitura for utilizado
func (rw *ReadWriteMutex) ReadLock() {
	// Bloqueia as goroutines pra nao ter race condition na operação de soma abaixo
	rw.readersLock.Lock()

	//A primeira goroutine de leitura bloqueia o mutex de escrita
	rw.readersCounter++
	if rw.readersCounter == 1 {
		rw.globalLock.Lock()
	}

	// Libera as goroutines pois estavam bloqueadas no inicio desse método
	rw.readersLock.Unlock()
}

// A função desse método é desbloquear o mutex de escrita
// Quando o mutex de leitura não for mais utilizado
func (rw *ReadWriteMutex) ReadUnlock() {
	// Bloqueia as goroutines pra nao ter race condition na operação de subtração abaixo
	rw.readersLock.Lock()

	//A última goroutine de leitura desbloqueia o mutex de escrita
	rw.readersCounter--
	if rw.readersCounter == 0 {
		rw.globalLock.Unlock()
	}

	// Libera as goroutines pois estavam bloqueadas no inicio desse método
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	rw.globalLock.Lock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.globalLock.Unlock()
}

func (rw *ReadWriteMutex) TryLock() bool {
	return rw.globalLock.TryLock()
}

/* Mas essa implementação tem um problema: a goroutine de escrita
nao conseguiria pegar o mutex enquanto tivesse goroutines de leitura,
apenas elas ficavam usando o mutex, monopolizando, e quando acabasse,
ai sim que a goroutine de leitura adquiria o mutex, esse cenario se chama
escassez de gravação, não podemos atualizar nossas estruturas de dados
compartilhadas porque as partes leitoras da execução as acessam continuamente,
bloqueando o acesso ao gravador */
