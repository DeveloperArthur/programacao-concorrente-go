/*Adapte o programa do terceiro exercício para continuar pesquisando
recursivamente em quaisquer subdiretórios. Se você fornecer um
arquivo à sua goroutine de pesquisa, ela deverá procurar uma correspondência
de string nesse arquivo, assim como nos exercícios anteriores. Caso contrário, se
você fornecer um diretório, ele deverá gerar recursivamente uma nova
goroutine para cada arquivo ou diretório encontrado dentro dele. Chame o
programa grepdirrec.go e execute-o executando este comando
go run grepdirrec.go bubbles ../../commonfiles :*/

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
	path := os.Args[2]

	fileInfo, _ := os.Stat(path)

	if fileInfo.IsDir() {
		fmt.Println("é um diretório")
		findFilesInDirectoryOrReadFiles(palavra, path)
	} else {
		fmt.Println("é um arquivo")
		go searchStringInFileByPath(palavra, path)
	}
	time.Sleep(1 * time.Second)

}

func findFilesInDirectoryOrReadFiles(palavra string, path string) {
	arquivos, _ := ioutil.ReadDir(path)
	for _, arquivo := range arquivos {
		if arquivo.IsDir() {
			go findFilesInDirectoryOrReadFiles(palavra, path+"/"+arquivo.Name())
		} else {
			go searchStringInFileByPath(palavra, path+"/"+arquivo.Name())
		}
	}
}

func searchStringInFileByPath(palavra string, path string) {
	conteudo, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	if strings.Contains(string(conteudo), palavra) {
		fmt.Printf("%s o arquivo %s tem uma correspondendia \n", time.Now().String(), path)
	}
}
