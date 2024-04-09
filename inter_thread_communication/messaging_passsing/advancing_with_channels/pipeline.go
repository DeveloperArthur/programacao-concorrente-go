package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	quit := make(chan int)
	//garante que um canal seja fechado quando a função que o criou terminar sua execução
	defer close(quit)

	//todos os métodos retornam canais
	urls := generateUrls(quit)
	pages := downloadPages(quit, urls)
	words := extractWords(quit, pages)

	for word := range words {
		fmt.Println(word)
	}
}

func generateUrls(quit <-chan int) <-chan string {
	//cria o canal
	urls := make(chan string)
	//todas as goroutines vao gravar no mesmo canal ao mesmo tempo
	go func() {
		//garante que um canal seja fechado quando a função que o criou terminar sua execução
		defer close(urls)
		for i := 100; i < 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)

			//urls <- url ERROR deadlocks
			/* colocar pra gravar url no canal aqui, em cima do select resulta em deadlocks
			porque o select vai bloquear a goroutine, ja que não há um case true que permita
			a continuação do loop, o select é assim, se nenhum dos casos estiver pronto para
			comunicação (ou seja, se nenhum dos canais (que estiverem no select) estiverem
			pronto para enviar ou receber), a goroutine ficará bloqueada no select até que
			pelo menos um dos casos esteja pronto, vai ficar bloqueado pra sempre e resulta
			em deadlocks */

			select {
			//a cada iteração será gravado a url no canal
			case urls <- url:
			case <-quit:
				return
			}
		}
	}()
	//retorna o canal no momento em que cria (go func executa assincrono)
	return urls
}

func downloadPages(quit <-chan int, urls <-chan string) <-chan string {
	//cria o canal
	pages := make(chan string)
	//todas as goroutines vao gravar no mesmo canal ao mesmo tempo
	go func() {
		//garante que um canal seja fechado quando a função que o criou terminar sua execução
		defer close(pages)
		moreData, url := true, ""
		for moreData {
			select {
			//lendo canal de urls
			case url, moreData = <-urls:
				if moreData {
					resp, _ := http.Get(url)
					if resp.StatusCode != 200 {
						panic("Server's error: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					//grava o conteudo da pagina no canal de pages
					pages <- string(body)
					resp.Body.Close()
				}
			case <-quit:
				return
			}
		}
	}()
	//retorna o canal no momento em que cria (go func executa assincrono)
	return pages
}

func extractWords(quit <-chan int, pages <-chan string) <-chan string {
	//cria o canal
	words := make(chan string)
	go func() {
		//garante que um canal seja fechado quando a função que o criou terminar sua execução
		defer close(words)
		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
		moreData, pg := true, ""
		for moreData {
			select {
			//lendo canal de pages
			case pg, moreData = <-pages:
				if moreData {
					for _, word := range wordRegex.FindAllString(pg, -1) {
						//grava palavra no canal word
						words <- strings.ToLower(word)
					}
				}
			case <-quit:
				return
			}
		}
	}()
	//retorna o canal no momento em que cria (go func executa assincrono)
	return words
}

/*generateUrls gera uma url, grava no canal de urls
ai downloadPages (que esta escutando o canal de urls) recebe a url
baixa a pagina, grava pagina no canal de pages
ai extractWords (que esta escutando o canal de pages) recebe a pagina,
extrai as palavras, grava as palavras no canal de words
e o main (que esta escutando o canal de words) recebe as palavras
e printa, tudo isso rodando ao mesmo tempo

Este pattern "pipeline" nos dá a capacidade de conectar
facilmente as execuções. Cada execução é representada por uma função que inicia
uma goroutine aceitando canais de entrada como argumentos e retornando
os canais de saída como valores de retorno.

O problema é que as páginas da web são baixadas
sequencialmente, uma após a outra, tornando a execução bastante lenta.
Idealmente, queremos acelerar isso e fazer os downloads concorrentemente.
É aqui que o próximo pattern (fan-in e fan-out) se torna útil, veja no arquivo */
