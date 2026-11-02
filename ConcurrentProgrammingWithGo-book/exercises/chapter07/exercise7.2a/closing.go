/*
 * Archivo:     closing.go
 *
 * Propósito:   Solución del ejercicio 7.2 del libro *Learn Concurrent
 *              Programming with Go* con un canal de `string`. El ejercicio pide
 *              ver qué recibe un receptor que no comprueba si el canal se
 *              cerró: después de `close`, cada lectura devuelve el valor cero
 *              del tipo.
 *
 * Compilación: go build closing.go
 * Ejecución:   ./closing
 *              (Alternativamente: go run closing.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Tres envíos y recepciones de "test" con marca de tiempo y luego,
 *              durante unos 3 segundos, líneas "Received:" con una cadena
 *              vacía.
 *
 * Notas:
 *    1. Leer de un canal cerrado nunca bloquea: devuelve el valor cero de
 *       inmediato.
 *    2. Para detectar el cierre use `msg, ok := <-messages` o `for msg := range
 *       messages`.
 *    3. La versión con `[]int` está en `../exercise7.2b/closing.go`.
 */

package main

import (
	"fmt"
	"time"
)

func receiver(messages <-chan string) {
	for {
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	msgChannel := make(chan string)
	go receiver(msgChannel)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		msgChannel <- "test"
		time.Sleep(1 * time.Second)
	}
	close(msgChannel)
	time.Sleep(3 * time.Second)
}
