/*
Essa é a "dependência transportada por loop", ocorre quando uma etapa em uma
iteração depende de outra etapa em uma iteração diferente no mesmo loop.
No caso aqui, a leitura dos arquivos precisa ser sequencial, nao pode ser
aleatoria a ordem
*/
package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	sha := sha256.New()
	var prev, next chan int
	for _, file := range files {
		if !file.IsDir() {
			next = make(chan int) // cria canal proximo
			go func(filename string, prev, next chan int) {
				fpath := filepath.Join(dir, filename)
				hashOnFile := FHash2(fpath)
				// essa checagem é feita pq a primeira iteracao nao tem anterior
				if prev != nil {
					<-prev // o go trava aqui até receber mensagem
				}
				sha.Write(hashOnFile)
				next <- 0 // coloca 0 no canal proximo
			}(file.Name(), prev, next)
			prev = next // o canal anterior recebe o proximo
		}
	}
	<-next
	fmt.Printf("%x\n", sha.Sum(nil))
}

func FHash2(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()

	sha := sha256.New()
	io.Copy(sha, file)

	return sha.Sum(nil)
}

/* então acontece o seguinte:
na primeira interacao é criado o canal next
dps de hashear adiciona 0 no canal next
e termina atribuindo o valor do canal next em prev
a segunda iteracao vai aguardar prev receber um valor
mesmo estando todos executando ao mesmo tempo, através
do canal eles se sincronizam, conforme for terminando,
vai enviando mensagem pro canal, e o canal vai liberando
o proximo, eles hasheam todos ao mesmo tempo, mas na hora
de adicionar no sha256, eles fazem sequencial

não da pra usar mutex, porque precisa ser na ordem
o mutex por mais que enfilere, ele vai executar um
de cada vez em ordem aleatória */
