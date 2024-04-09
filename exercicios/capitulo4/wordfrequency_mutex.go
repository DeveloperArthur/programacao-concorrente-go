/* No capítulo 3 você fez um exercicio chamado wordfrequency.go
um programa para gerar as frequências de palavras de páginas da web baixadas.
Se você usasse um mapa de memória compartilhada para armazenar as
frequências das palavras, o acesso ao mapa compartilhado precisaria ser
protegido. Você pode usar um mutex para garantir acesso exclusivo ao mapa? */

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {
	var frequency = make(map[string]int)
	mutex := sync.Mutex{}
	for i := 1000; i <= 1020; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countWords(url, &frequency, &mutex)
	}
	time.Sleep(10 * time.Second)
	mutex.Lock()
	for k, v := range frequency {
		fmt.Println(k, "->", v)
	}
	mutex.Unlock()
}

func countWords(url string, frequency *map[string]int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
	mutex.Lock()
	for _, word := range wordRegex.FindAllString(string(body), -1) {
		wordLower := strings.ToLower(word)
		//desreferenciando mapa
		(*frequency)[wordLower] += 1
	}
	mutex.Unlock()
	fmt.Println("Completed:", url)
}
