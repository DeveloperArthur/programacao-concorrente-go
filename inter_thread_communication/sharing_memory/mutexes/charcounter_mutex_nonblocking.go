/* Em algumas aplicações, podemos não querer bloquear a
goroutine, mas sim realizar algum outro trabalho
antes de tentar novamente bloquear o mutex e acessar a seção crítica

quando o mutex.lock é chamado, uma goroutine é executada de cada vez,
a goroutine adquire um passe e vai ser executada, o passe seria o mutex
as N goroutines ficam numa fila esperando para adquirir o mutex

no nosso código ao invés de imprimir a frequencia só no final
a goroutine main vai tentar de 100 em 100 milisegundos  ver se tem um
mutex disponivel, se tiver ela executa um código, as demais
goroutines na fila continuam lá esperando, e depois a main libera
o mutex para as outras goroutines da fila.

o trylock retorna true quando chamada é feira num espaço de tempo entre a liberação
do mutex e a disponibilidade para a próxima goroutine na fila de espera */

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const allLetters2 = "abcdefghijklmnopqrstuvwxyz"

func main() {
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)
	for i := 2000; i <= 2200; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters2(url, frequency, &mutex)
	}

	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
		if mutex.TryLock() {
			for i, c := range allLetters2 {
				fmt.Printf("%c-%d ", c, frequency[i])
			}
			mutex.Unlock()
		} else {
			fmt.Println("Mutex already being used")
		}
	}
}

func countLetters2(url string, frequency []int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)

	mutex.Lock()
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters2, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	mutex.Unlock()
	fmt.Println("Completed:", url)
}
