/* Também podemos usar a instrução select quando precisarmos escrever mensagens
em canais, não apenas quando estivermos lendo mensagens de canais.

Imagine que temos que encontrar 100 números primos aleatórios. Na vida real,
poderíamos escolher um número aleatório de um saco com um grande conjunto de
números e depois manter esse número apenas se for primo */

package main

import (
	"fmt"
	"math"
	"math/rand"
)

// lê os primos enviados por primesOnly e grava os numeros aleatorios no canal
func main() {
	numbersChannel := make(chan int)
	primes := primesOnly(numbersChannel)
	for i := 0; i < 100; {
		select {
		//grava um numero random no canal numbersChannel
		case numbersChannel <- rand.Intn(1000000000) + 1:

		//lê as msgs do canal primes
		case p := <-primes:
			fmt.Println("Found prime:", p)
			i++
		}
	}
}

// lê os numeros aleatorios enviados por main e grava os primos no canal
func primesOnly(inputs <-chan int) <-chan int {
	primes := make(chan int)
	go func() {
		//lê as msgs do canal numbersChannel
		for c := range inputs {
			isPrime := c != 1
			for i := 2; i <= int(math.Sqrt(float64(c))); i++ {
				if c%i == 0 {
					isPrime = false
					break
				}
			}
			if isPrime {
				//grava o resultado no canal primes
				primes <- c
			}
		}
	}()

	//retorna em paralelo
	return primes
}
