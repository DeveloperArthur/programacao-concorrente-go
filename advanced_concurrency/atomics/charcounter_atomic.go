package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	wg := sync.WaitGroup{}
	wg.Add(31)
	var frequency = make([]int32, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go func() {
			countLetters(url, frequency)
			wg.Done()
		}()
	}
	wg.Wait()
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, atomic.LoadInt32(&frequency[i]))
	}
}

func countLetters(url string, frequency []int32) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			atomic.AddInt32(&frequency[cIndex], 1)
		}
	}
	fmt.Println("Completed:", url)
}

/*
Quando executamos nossa aplicação de charcounter sem nenhum bloqueio mutex
(inter_thread_communication/sharing_memory/race_condition/charcounter_race_condition.go),
ela produziu resultados inconsistentes. Usar as variáveis atômicas tem o mesmo efeito que
eliminar a race condition usando um mutex. Porém, desta vez, nossas goroutines não
estão bloqueando umas às outras.

NOTA: O uso da função LoadInt32() não é estritamente necessário na listagem
anterior porque todas as goroutines estão concluídas no momento em que lemos
os resultados. No entanto, é uma boa prática usar operações de carga atômica ao
trabalhar com atômicos para garantir que lemos o valor mais recente da memória
principal e não um valor desatualizado em cache.
*/
