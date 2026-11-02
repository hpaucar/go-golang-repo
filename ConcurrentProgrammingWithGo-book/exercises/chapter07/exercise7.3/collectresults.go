/*
 * Archivo:     collectresults.go
 *
 * Propósito:   Solución del ejercicio 7.3 del libro *Learn Concurrent
 *              Programming with Go*, que pide recoger los resultados de varias
 *              goroutines con canales. Diez goroutines calculan los divisores
 *              de números aleatorios y cada una envía su resultado por un canal
 *              propio.
 *
 * Compilación: go build collectresults.go
 * Ejecución:   ./collectresults
 *              (Alternativamente: go run collectresults.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Diez listas de divisores, una por línea, en el orden de los
 *              canales (0 a 9). Puede tardar varios segundos.
 *
 * Notas:
 *    1. Un canal por goroutine permite recoger los resultados en un orden fijo,
 *       sin importar cuál termina primero.
 *    2. `findFactors` prueba todos los números hasta `number` (hasta mil
 *       millones), así que el cálculo es costoso a propósito.
 *    3. El índice se pasa como argumento `n` a la función anónima para que cada
 *       goroutine use su propio canal.
 */

package main

import (
	"fmt"
	"math/rand"
)

func findFactors(number int) []int {
	result := make([]int, 0)
	for i := 1; i <= number; i++ {
		if number%i == 0 {
			result = append(result, i)
		}
	}
	return result
}

func main() {
	resultChs := make([]chan []int, 10)
	for i := 0; i < 10; i++ {
		resultChs[i] = make(chan []int)
		go func(n int) {
			resultChs[n] <- findFactors(rand.Intn(1000000000))
		}(i)
	}
	for i := 0; i < 10; i++ {
		fmt.Println(<-resultChs[i])
	}
}
