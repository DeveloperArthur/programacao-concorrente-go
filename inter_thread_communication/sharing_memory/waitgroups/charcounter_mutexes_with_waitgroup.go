package main

import (
	"fmt"
	"sync"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	wg := sync.WaitGroup{}
	wg.Add(31)

	mutex := sync.Mutex{}
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)

		/* Observe que não precisamos modificar a função countLetters() para chamar Done();
		em vez disso, usamos uma função anônima executada em uma goroutine separada, chamando
		ambas as funções */
		go func() {
			countLetters(url, frequency, &mutex)
			wg.Done()
		}()
	}

	//time.Sleep(10 * time.Second) -> não esperamos mais 10 segundos
	//executamos quando o waitgroup finalizar:
	wg.Wait()

	mutex.Lock()
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock()
}

func countLetters(url string, frequency []int, mutex *sync.Mutex) {}
