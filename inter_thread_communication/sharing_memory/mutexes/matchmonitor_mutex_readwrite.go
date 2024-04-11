/* imagine 3 goroutines, 2 de leitura e 1 de gravação
enquanto a goroutine de gravação trabalha, as outras duas esperam
quando as goroutine de leitura trabalham, a goroutine de gravação espera
é isso que o mutex de leitura e escrita nos dá 

O RWMutex em si é responsável por controlar o acesso concorrente ao valor que
ele protege. Ele permite que várias threads façam operações de leitura ao mesmo 
tempo (read), mas garante que apenas uma thread por vez possa fazer uma operação 
de escrita (write), bloqueando as outras threads que tentarem fazer uma escrita 
simultaneamente. */

package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	mutex := sync.RWMutex{}
	var matchEvents = make([]string, 0, 10000)

	for j := 0; j < 1000; j++ {
		matchEvents = append(matchEvents, "Match event")
	}

	go matchRecorder2(&matchEvents, &mutex)

	start := time.Now()
	for j := 0; j < 5000; j++ {
		go clientHandler2(&matchEvents, &mutex, start)
	}
	time.Sleep(100 * time.Second)
}

func matchRecorder2(matchEvents *[]string, mutex *sync.RWMutex) {
	for i := 0; ; i++ {

		/* Quando há uma atualização do jogo, a goroutine no matchRecorder() adquire um bloqueio
		de gravação chamando a função Lock() no mutex. O bloqueio de gravação só será
		adquirido quando qualquer goroutine matchRecorder() ativa liberar seu bloqueio de
		leitura. Quando o bloqueio de gravação é adquirido, ele bloqueará qualquer outra
		goroutine de acessar a seção crítica em nossa função clientHandler() até que
		liberemos o bloqueio de gravação chamando UnLock(). */
		mutex.Lock()
		*matchEvents = append(*matchEvents,
			"Match event "+strconv.Itoa(i))
		mutex.Unlock()
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Appended match event")
	}
}

func clientHandler2(mEvents *[]string, mutex *sync.RWMutex, st time.Time) {
	for i := 0; i < 100; i++ {

		/* Os bloqueios de leitura usados aqui são necessários
		porque não queremos que a estrutura de dados da slice mude
		enquanto a percorremos. Por exemplo, modificar o ponteiro e
		o conteúdo da slice enquanto outra goroutine a lê pode nos levar
		a seguir uma referência de ponteiro inválida.. */
		mutex.RLock()
		allEvents := copyAllEvents2(mEvents)
		mutex.RUnlock()

		timeTaken := time.Since(st)
		fmt.Println(len(allEvents), "events copied in ", timeTaken)
	}
}

func copyAllEvents2(matchEvents *[]string) []string {
	allEvents := make([]string, 0, len(*matchEvents))
	for _, e := range *matchEvents {
		allEvents = append(allEvents, e)
	}
	return allEvents
}

/* Uma goroutine executando a seção de código crítica entre RLock() e RUnlock(),
em nossa função clientHandler() , bloqueia a goroutine de adquirir um mutex de
gravação em nossa função matchRecorder()

No entanto, isso não impede que outra goroutine também adquira um bloqueio de
leitor para uma seção crítica. Isto significa que podemos ter goroutines concorrentes
executando clientHandler() sem nenhuma goroutine de leitura bloqueando umas às outras*/
