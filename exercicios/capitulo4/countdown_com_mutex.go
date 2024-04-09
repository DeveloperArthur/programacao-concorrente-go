/* O programa inter_thread_communication/sharing_memory/countdown.go não
usa nenhum mutex para proteger o acesso à sua variável compartilhada.
Esta é uma má prática.
Altere este programa para que o acesso à variável de segundos
compartilhada seja protegido por um mutex. Dica: pode ser necessário copiar
uma variável. */

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	count := 5
	mutex := sync.Mutex{}

	go countdown(&count, &mutex)

	remaining := count
	for remaining > 0 {
		time.Sleep(500 * time.Millisecond)
		mutex.Lock()
		fmt.Println(count)
		remaining = count
		mutex.Unlock()
	}
}

func countdown(seconds *int, mutex *sync.Mutex) {
	mutex.Lock()
	remaining := *seconds
	mutex.Unlock()

	for remaining > 0 {
		time.Sleep(1 * time.Second)
		mutex.Lock()
		*seconds -= 1
		remaining = *seconds
		mutex.Unlock()
	}
}

/* O uso da variável remaining ajuda a evitar que as funções lidem
diretamente com o mesmo endereço de memória (count) de forma não
sincronizada. Ao copiar o valor de count para remaining dentro do
mutex, garantimos que cada função (countdown e main) tenha sua
própria cópia local do valor de count durante sua execução.
Isso evita que a leitura e a escrita concorrentes no mesmo endereço
de memória causem problemas de concorrência, como race conditions */
