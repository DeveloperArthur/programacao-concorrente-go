/* Altere o programa que você escreveu no segundo exercício para que, em vez
de passar uma lista de nomes de arquivos de texto, você passe um caminho
de diretório. O programa irá olhar dentro deste diretório e listar os arquivos.
Para cada arquivo, você pode gerar uma goroutine que procurará uma
correspondência de string (a mesma de antes). Chame o programa grepdir.go.
Veja como você pode executar este programa Go:
go run grepdir.go bubbles ../../commonfiles*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	palavra := os.Args[1]
	diretorio := os.Args[2]

	arquivos, _ := ioutil.ReadDir(diretorio)

	for _, arquivo := range arquivos {
		fmt.Println(time.Now().String() + " " + arquivo.Name())
		go searchStringInFile(palavra, arquivo.Name())
	}
	time.Sleep(1 * time.Second)
}

func searchStringInFile(palavra string, file string) {
	conteudo, err := ioutil.ReadFile("./files/" + file)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	if strings.Contains(string(conteudo), palavra) {
		fmt.Printf("%s o arquivo %s tem uma correspondendia \n", time.Now().String(), file)
	}
}
