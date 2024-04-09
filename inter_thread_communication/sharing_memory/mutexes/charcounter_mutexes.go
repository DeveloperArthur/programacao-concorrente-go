/* IMPORTANTE: de acordo com a lei de Amdahl, a proporção sequencial
para paralelo limitará a escalabilidade de desempenho do nosso código,
por isso é essencial reduzirmos o tempo gasto mantendo o bloqueio mutex
ou seja, quanto menos código deixarmos bloqueado melhor, precisamos
deixar bloqueado cirurgicamente só o que precisa ser bloqueado,
pois esses bloqueios afetam diretamente nosso desempenho */

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency, &mutex)
	}
	time.Sleep(10 * time.Second)
	mutex.Lock()
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock()
}

func countLetters(url string, frequency []int, mutex *sync.Mutex) {
	//Executa a parte lenta da função (o download) concorrentemente
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)

	//Bloqueia apenas a seção de processamento crítico e rápido da função, essa parte é feita sequencialmente
	mutex.Lock()
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	mutex.Unlock()
	fmt.Println("Completed:", url)
}

/* Se colocarmos o lock no inicio da nossa função countLetters() e o unlock
no final dela, estariamos transformamos nosso programa concorrente em um
programa sequencial. Acabaremos baixando e processando uma página da web
por vez, pois bloquearemos desnecessariamente toda a execução. Se prosseguirmos
e executarmos isso, o tempo gasto será o mesmo da versão não concorrente do
programa.

Em termos de desempenho, não faz sentido bloquear toda a execução. A etapa de
download do documento não compartilha nada com outras goroutines, portanto não há
risco de ocorrer uma race condition.

Ao decidir como e quando usar mutexes, é melhor focar em quais recursos
devemos proteger e descobrir onde as seções críticas começam e terminam.
Então precisamos pensar em como minimizar o número de chamadas Lock() e Unlock()

Dependendo da implementação do mutex, geralmente há um custo de desempenho
se chamarmos as operações Lock() e Unlock() com muita frequência, nós poderíamos
muito bem usar o mutex para proteger apenas a instrução de frequency[cIndex] += 1
No entanto, isso significa que chamariamos Lock() e Unlock() para
cada letra do documento baixado, pois está dentro do loop. Por isso temos um melhor desempenho
chamando Lock() antes do loop e Unlock() depois de sairmos do loop. */
