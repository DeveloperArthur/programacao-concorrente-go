/*Expanda o programa que você escreveu no primeiro exercício para que
em vez de imprimir o conteúdo dos arquivos de texto, ele procura uma
correspondência de string. A string a ser pesquisada é o primeiro argumento na
linha de comando. Quando você gera uma nova goroutine, em vez de imprimir
o conteúdo do arquivo, ela deve ler o arquivo e procurar uma correspondência. Se
a goroutine encontrar uma correspondência, ela deverá exibir uma mensagem
informando que o nome do arquivo contém uma correspondência. Chame o
programa grepfiles.go. Veja como você pode executar este programa Go
(“bolhas” é a string de pesquisa neste exemplo):
go run grepfiles.go bubbles txtfile1 txtfile2 txtfile3 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	for i := 2; i < len(os.Args); i++ {
		go searchStringIn(os.Args[1], os.Args[i])
	}
	time.Sleep(1 * time.Second)
}

func searchStringIn(palavra string, file string) {
	conteudo, err := ioutil.ReadFile("./files/" + file + ".txt")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	if strings.Contains(string(conteudo), palavra) {
		fmt.Printf("%s o arquivo %s tem uma correspondendia \n", time.Now().String(), file)
	}
}
