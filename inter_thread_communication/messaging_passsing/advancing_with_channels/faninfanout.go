/* No programa anterior (pipeline.go) se quisermos acelerar as coisas, podemos
realizar os downloads simultaneamente, balanceando a carga das URLs para
vários goroutines. Podemos criar um número fixo de goroutines, cada uma
lendo o mesmo canal de entrada de URL. Cada uma das goroutines receberá
uma URL separada da goroutine generateUrls() , e poderão realizar os
downloads simultaneamente. As páginas de texto baixadas podem então ser
escritas no próprio canal de saída de cada goroutine.

No Go, o fan-out, pattern de concorrência, ocorre quando
várias goroutines são lidas no mesmo canal.
Desta forma, podemos distribuir o trabalho entre um
conjunto de goroutines, nesse exemplo vamos criar 20
goroutines, ambas recebendo o mesmo canal retornado
pelo generateUrls() por parâmetro, conforme as mensagens
forem ficando disponíveis, as goroutines vão pegando
não é que vai ser gerado 1 goroutine pra cada URL
todas as goroutines vão olhar pro mesmo canal, semelhante
ao código presente em messaging_passing/channel/workers.go

No Go, o fan-in, pattern de concorrência, ocorre quando mesclamos
o conteúdo de vários canais em um, nesse exemplo,downloadPages
terá diversas goroutines executando 1 URL cada, e fan-in
é para mesclar, unificar as respostas de todas as goroutines para
o extractWords()

fan-in e fan-out: https://miro.medium.com/v2/resize:fit:1374/1*wWMhMypygxHMJ4_uBbeyqA.png

NOTA: Como o processamento concorrente não é determinístico,
algumas mensagens serão processadas mais rapidamente do que
outras, resultando em mensagens processadas em uma ordem imprevisível.
O pattern fan-out só faz sentido se não nos importarmos com a
ordem das mensagens recebidas. */

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const downloaders = 20

func main() {
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls2(quit)

	//criando um slice do tipo output channel do tipo string
	pages := make([]<-chan string, downloaders)

	//fan-out
	for i := 0; i < downloaders; i++ {
		//insere no slice o canal que downloadPages2() retorna
		pages[i] = downloadPages2(quit, urls)
	}

	paginasUnificadas := FanIn(quit, pages...)

	words := extractWords2(quit, paginasUnificadas)
	for word := range words {
		fmt.Println(word)
	}
}

func generateUrls2(quit <-chan int) <-chan string {
	urls := make(chan string)
	go func() {
		defer close(urls)
		for i := 100; i < 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			select {
			case urls <- url:
			case <-quit:
				return
			}
		}
	}()
	return urls
}

func downloadPages2(quit <-chan int, urls <-chan string) <-chan string {
	pages := make(chan string)
	go func() {
		defer close(pages)
		moreData, url := true, ""
		for moreData {
			select {
			case url, moreData = <-urls:
				if moreData {
					resp, _ := http.Get(url)
					if resp.StatusCode != 200 {
						panic("Server's error: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
					resp.Body.Close()
				}
			case <-quit:
				return
			}
		}
	}()
	return pages
}

func extractWords2(quit <-chan int, pages <-chan string) <-chan string {
	words := make(chan string)
	go func() {
		defer close(words)
		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
		moreData, pg := true, ""
		for moreData {
			select {
			case pg, moreData = <-pages:
				if moreData {
					for _, word := range wordRegex.FindAllString(pg, -1) {
						words <- strings.ToLower(word)
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return words
}

/*
a goroutine downloadPages2() retorna 1 canal, se temos 20 goroutines dessa
temos 20 canais, portanto esse método vai receber como parametro todos os canais
que as goroutines de downloadPages2() retornaram, itera por todos os canais
criando uma goroutine para cada iteração de canal, nessa goroutine faz um loop
pra adicionar mensagem por mensagem no canal de saída.
e tem um wait group pra esperar todas essas goroutines executarem antes de
retornar o canal de output
*/
func FanIn[K any](quit <-chan int, allChannels ...<-chan K) chan K {
	wg := sync.WaitGroup{}
	wg.Add(len(allChannels))
	output := make(chan K)
	for _, c := range allChannels {
		go func(channel <-chan K) {
			defer wg.Done()
			for i := range channel {
				select {
				case output <- i:
				case <-quit:
					return
				}
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}

/* Quando executamos essa nova implementação, ela roda muito mais rápido porque os
downloads estão sendo realizados concorrentemente. Como efeito colateral de
fazermos os downloads juntos, a ordem das palavras extraídas é diferente cada
vez que executamos o programa, o pipeline.go executa com cerca de 10 segundos
comparado a este programa que executa com certa de 3 segundos */
