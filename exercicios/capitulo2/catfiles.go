/*Escreva um programa que aceite uma lista de nomes
de arquivos de texto como argumentos. Para cada nome de arquivo, o programa
deve gerar uma nova goroutine que enviará o conteúdo desse arquivo para o
console. Você pode usar a função time.Sleep() para aguardar a conclusão das
goroutines filhas (até que você saiba como fazer isso melhor). Chame o programa
catfiles.go. Veja como você pode executar este programa Go:
go run catfiles.go txtfile1 txtfile2 txtfile3 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		go readFile(os.Args[i])
	}
	time.Sleep(1 * time.Second)
}

func readFile(file string) {
	conteudo, err := ioutil.ReadFile("./files/" + file + ".txt")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	fmt.Println(time.Now().String() + string(conteudo))
}
