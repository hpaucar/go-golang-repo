/*
 * Archivo:     closing.go
 *
 * Propósito:   Solución del ejercicio 7.2 del libro *Learn Concurrent
 *              Programming with Go* con un canal de `[]int`. Igual que en la
 *              versión 7.2a, el receptor no comprueba si el canal se cerró, y
 *              después de `close` recibe el valor cero del tipo, que para un
 *              slice es `nil`.
 *
 * Compilación: go build closing.go
 * Ejecución:   ./closing
 *              (Alternativamente: go run closing.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Tres recepciones de `[0 0 0 0 0 0 0 0 0 0]` con marca de tiempo
 *              y luego, durante unos 3 segundos, líneas "Received: []" (slice
 *              `nil`).
 *
 * Notas:
 *    1. El valor cero depende del tipo: "" para `string`, 0 para `int`, `nil`
 *       para slices, mapas y punteros.
 *    2. Un slice `nil` se imprime como `[]`, igual que uno vacío, lo que puede
 *       ocultar el error.
 */

package main

import (
	"fmt"
	"time"
)

func receiver(messages <-chan []int) {
	for {
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	msgChannel := make(chan []int)
	go receiver(msgChannel)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		msgChannel <- make([]int, 10)
		time.Sleep(1 * time.Second)
	}
	close(msgChannel)
	time.Sleep(3 * time.Second)
}
