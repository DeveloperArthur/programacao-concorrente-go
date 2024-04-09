/* Modifique nosso programa de frequência de letras sequenciais para produzir uma lista
de frequências de palavras em vez de frequências de letras. Você pode usar as
mesmas URLs para as páginas da web RFC usadas na listagem 3.3. Ao terminar, o
programa deverá gerar uma lista de palavras com a frequência com que cada
palavra aparece na página web. Aqui está um exemplo de saída:
go run wordfrequency.go
the -> 5
a -> 8
car -> 1
program -> 3 */

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func main() {
	var frequency = make(map[string]int)
	for i := 1000; i <= 1020; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		countWords(url, frequency)
	}
	time.Sleep(10 * time.Second)
	for k, v := range frequency {
		fmt.Println(k, "->", v)
	}
}

func countWords(url string, frequency map[string]int) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
	for _, word := range wordRegex.FindAllString(string(body), -1) {
		wordLower := strings.ToLower(word)
		frequency[wordLower] += 1
	}
	fmt.Println("Completed:", url)
}
