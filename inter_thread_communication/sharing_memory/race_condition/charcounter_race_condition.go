/* Vejamos agora um exemplo envolvendo mais de duas goroutines, onde as goroutines
atualizam as mesmas variáveis ao mesmo tempo. Para este exemplo, escreveremos
um programa para descobrir com que frequência as letras do alfabeto inglês aparecem
em textos comuns. O programa processará páginas da web baixando-as e contando
com que frequência cada letra do alfabeto aparece nas páginas. Quando o
programa for concluído, ele deverá nos fornecer uma tabela de frequência com a
contagem de quantas vezes cada caractere ocorre */

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency)
	}
	time.Sleep(10 * time.Second)
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
}

func countLetters(url string, frequency []int) {
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
			frequency[cIndex] += 1
		}
	}
	fmt.Println("Completed:", url)
}

/* executando a função main sem concorrencia (sem a palavra chave go na linha 24)
o programa leva 17 segundos para ser concluído, o resultado:

a-103445 b-23074 c-61005 d-51733 e-181360 f-33381 g-24966
h-47722 i-103262 j-3279 k-8839l-49958 m-40026 n-108275
o-106320 p-41404 q-3410 r-101118 s-101040 t-136812 u-35765
v-13666 w-18259 x-4743 y-18416 z-1404

colocando a palavra go na frente da chamada da função estamos criando 31 goroutines
cada goroutine baixará e processará seu documento concorrentemente, o problema
é que usar a função countLetters() de múltiplas goroutines produzirá resultados
errôneos devido a race condition, pois uma goroutine pode escrever um valor em
uma variável, mas quando ela o lê de volta, o valor pode ser diferente,
pois outra goroutine pode tê-lo alterado

o resultado executando concorrentemente:

a-102692 b-23033 c-60678 d-51516 e-178892 f-33317 g-24922
h-47607 i-102559 j-3279 k-8827 l-49818 m-39911 n-107403
o-105612 p-41266 q-3407 r-100217 s-100314 t-135517 u-35645
v-13646 w-18239 x-4742 y-18396 z-1404

existem algumas inconsistencias, e isso é por causa da race condition

race condition é quando temos vários threads (ou processos) compartilhando
um recurso e eles se sobrepõem, gerando resultados inesperados.

solução em:
- inter_thread_communication/sharing_memory/mutexes/charcounter_mutexes.go
- inter_thread_communication/sharing_memory/mutexes/charcounter_mutex_nonblocking.go
- inter_thread_communication/sharing_memory/waitgroups/charcounter_mutexes_with_waitgroup.go
*/
