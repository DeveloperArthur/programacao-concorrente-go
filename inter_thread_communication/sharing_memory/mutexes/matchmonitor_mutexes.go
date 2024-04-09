/* E se tivéssemos um aplicativo servindo principalmente dados estáticos para muitos
clientes simultâneos? um aplicativo que fornece atualizações aos usuários
sobre um jogo de basquete.
Neste aplicativo, os usuários verificam atualizações sobre um jogo de basquete ao vivo em
seus dispositivos. O aplicativo Go, executado em nossos servidores, fornece essas
atualizações

as duas goroutines estao compartilhando o mesmo slice, uma goroutine adiciona
dados no slice de 200 em 200 milisegundos, e as outras 5000 goroutines
que sao 5000 usuarios, copiam todos os dados do slice 100 vezes, como se cada
usuario estivesse solicitando 100 vezes, e printa no console */

package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	mutex := sync.Mutex{}
	var matchEvents = make([]string, 0, 10000)

	// Preenche slice de eventos com muitos eventos, simulando um jogo em andamento
	for j := 0; j < 1000; j++ {
		matchEvents = append(matchEvents, "Match event")
	}

	go matchRecorder(&matchEvents, &mutex)

	//Starts a large number of client handler goroutines
	start := time.Now()
	for j := 0; j < 5000; j++ {
		go clientHandler(&matchEvents, &mutex, start)
	}
	time.Sleep(100 * time.Second)
}

// Adiciona uma string contendo um evento a cada 200 milissegundos
func matchRecorder(matchEvents *[]string, mutex *sync.Mutex) {
	for i := 0; ; i++ {
		/* Com apenas 1 goroutine no matchRecorder, não faz sentido usar Lock
		pois não há fila de goroutines, mas futuramente este código pode ser alterado
		para permitir diversas goroutines executando matchRecorder, por isso é
		uma boa prática colocar Lock em todo código que grava em dado compartilhado */
		mutex.Lock()
		*matchEvents = append(*matchEvents,
			"Match event "+strconv.Itoa(i))
		mutex.Unlock()
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Appended match event")
	}
}

func clientHandler(mEvents *[]string, mutex *sync.Mutex, st time.Time) {
	//A função clientHandler() possui um
	//loop que se repete 100 vezes para simular o mesmo usuário fazendo várias
	//solicitações
	for i := 0; i < 100; i++ {

		/* Necessário para garantir que as goroutines de leitura
		sempre leiam os dados atualizados, mas ler dessa forma
		não é a melhor maneira */
		mutex.Lock()

		//Esta função simula a construção
		//de uma resposta para enviar ao usuário. No mundo real, poderíamos enviar esta
		//resposta formatada em algo como JSON
		allEvents := copyAllEvents(mEvents)

		mutex.Unlock()

		timeTaken := time.Since(st)
		fmt.Println(len(allEvents), "events copied in ", timeTaken)
	}
}

// Copia todo o conteúdo do slice correspondente, simulando a construção de uma resposta para o cliente
func copyAllEvents(matchEvents *[]string) []string {
	allEvents := make([]string, 0, len(*matchEvents))
	for _, e := range *matchEvents {
		allEvents = append(allEvents, e)
	}
	return allEvents
}

/*
Gravar dados com 1 goroutine e ler dados com várias goroutines, há risco de race condition
porque se uma goroutine de leitura ler o array enquanto a goroutine de escrita está no meio
de adicionar um novo número, essa leitura pode resultar em um array incompleto ou em um estado
inconsistente. Da mesma forma, se várias goroutines de leitura tentarem ler o array
simultaneamente enquanto a goroutine de escrita está modificando-o, elas podem ler dados incompletos
ou parcialmente atualizados, por isso é importante usar Lock nas leituras, para garantir que
as goroutines de leitura sempre leiam os dados atualizados.

Mas quando usamos bloqueios nas leituras com mutex normais, toda vez que uma goroutine lê os dados
compartilhados do basquete, ela bloqueia todas as outras goroutines em serviço
até terminar. Mesmo que os manipuladores do cliente estejam apenas lendo a slice
compartilhada sem nenhuma modificação, ainda damos a cada um deles acesso
exclusivo à fatia. Observe que se várias goroutines estiverem apenas lendo dados
compartilhados sem atualizá-los, não há necessidade desse acesso exclusivo; a
leitura simultânea de dados compartilhados não causa nenhuma interferência,
nós podemos gravar bloqueando com mutex e ler de forma concorrente, não há risco de race condition
nesse cenário.

	NOTA: As race conditions só acontecem se alterarmos o estado compartilhado
	sem a sincronização adequada. Se não modificarmos os dados compartilhados,
	e apenas lê-los concorrentemente, não há risco de race condition

Seria melhor se todas as goroutines do manipulador de cliente tivessem acesso não
exclusivo à slice para que pudessem ler a lista ao mesmo tempo, se necessário. Isso
melhoraria o desempenho, pois permitiria que vários goroutines que estão apenas
lendo os dados compartilhados pudessem acessá-los ao mesmo tempo. Só
bloquearíamos o acesso aos dados compartilhados se houvesse necessidade de
atualizá-los. */
