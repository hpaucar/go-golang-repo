/*
 * Archivo:     stopreadingafter.go
 *
 * Propósito:   Solución del ejercicio 8.2 del libro *Learn Concurrent
 *              Programming with Go*, que pide leer números aleatorios de un
 *              canal y dejar de leer después de 5 segundos. Se usa `time.After`
 *              como un caso más del `select`.
 *
 * Compilación: go build stopreadingafter.go
 * Ejecución:   ./stopreadingafter
 *              (Alternativamente: go run stopreadingafter.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Unos 25 números aleatorios entre 0 y 9 (uno cada 200 ms) y, a
 *              los 5 segundos, "Stopping reading after timeout".
 *
 * Notas:
 *    1. `time.After` devuelve un canal que recibe un valor cuando pasa el
 *       tiempo indicado; se crea una sola vez fuera del bucle.
 *    2. La goroutine generadora queda bloqueada al salir `main`; en código real
 *       se usaría un canal de cancelación o `context`.
 */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateNumbers() chan int {
	output := make(chan int)
	go func() {
		for {
			output <- rand.Intn(10)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return output
}

func main() {
	numbers := generateNumbers()
	timeout := time.After(5 * time.Second)
	for {
		select {
		case n := <-numbers:
			fmt.Println(n)
		case <-timeout:
			fmt.Println("Stopping reading after timeout")
			return
		}
	}
}
