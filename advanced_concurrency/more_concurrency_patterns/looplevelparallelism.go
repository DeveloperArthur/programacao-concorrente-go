/* Quando temos uma coleção de dados na qual precisamos executar uma tarefa, podemos
usar a concorrencia para executar várias tarefas em diferentes partes da coleção ao
mesmo tempo. Um programa serial pode ter um loop para executar a tarefa em cada
item da coleção, um após o outro.

O padrão de paralelismo em nível de loop transforma cada tarefa de iteração em uma tarefa
concorrente para que possa ser executada em paralelo.

Suponha que tenhamos que criar um programa para calcular o código hash de uma lista
de arquivos em um diretório específico. Na programação sequencial,
criaríamos uma função de hash de arquivo.

Então faríamos nosso programa coletar uma lista de arquivos do diretório
e iterá-los. Em cada iteração, chamaríamos nossa função hash e
imprimiríamos os resultados.

Em vez de processar cada arquivo no diretório um após o outro sequencialmente,
poderíamos usar o paralelismo em nível de loop e alimentar cada arquivo em uma goroutine
separada

O código abaixo lê todos os arquivos de um diretório especificado e depois
itera cada arquivo em um loop. Para cada iteração, ele inicia uma nova
goroutine para calcular o código hash do arquivo nessa iteração.

Esta listagem usa um grupo de espera para pausar a goroutine main() até que todas
as tarefas sejam concluídas. */

package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	wg := sync.WaitGroup{}
	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			//cria uma nova goroutine para cada arquivo
			go func(filename string) {
				fPath := filepath.Join(dir, filename)
				hash := FHash(fPath)
				fmt.Printf("%s - %x\n", filename, hash)
				wg.Done()
			}(file.Name())
		}
	}
	wg.Wait()
}

func FHash(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()

	sha := sha256.New()
	io.Copy(sha, file)

	return sha.Sum(nil)
}

/* Neste exemplo, podemos facilmente usar o padrão de paralelismo em nível de loop
porque não há dependência entre as tarefas. O resultado do cálculo do código hash
de um arquivo não afeta o cálculo do código hash do próximo arquivo. Se tivéssemos
processadores suficientes, poderíamos executar cada iteração em um processador
separado. Mas e se o cálculo de uma iteração depender de uma etapa calculada
em uma iteração anterior? looplevelparallelism2.go */
