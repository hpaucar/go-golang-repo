/*
 * Archivo:     messagepassingstop.go
 *
 * Propósito:   Solución del ejercicio 7.1 del libro *Learn Concurrent
 *              Programming with Go*, que pide que `main` espere a que el
 *              receptor procese el mensaje "STOP" antes de terminar. El
 *              receptor responde por el mismo canal sin búfer y `main` lo
 *              espera con `<-msgChannel`.
 *
 * Compilación: go build messagepassingstop.go
 * Ejecución:   ./messagepassingstop
 *              (Alternativamente: go run messagepassingstop.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los mensajes "Sending ..." y "Received: ..." para HELLO, THERE y
 *              STOP, siempre con los tres mensajes recibidos antes de salir.
 *
 * Notas:
 *    1. Un canal sin búfer se puede usar en ambas direcciones como señal de
 *       "terminé"; aquí el receptor envía una cadena vacía.
 *    2. Sin esta respuesta, `main` podría salir antes de que se imprima
 *       "Received: STOP".
 */

package main

import "fmt"

func main() {
	msgChannel := make(chan string)
	go receiver(msgChannel)
	fmt.Println("Sending HELLO...")
	msgChannel <- "HELLO"
	fmt.Println("Sending THERE...")
	msgChannel <- "THERE"
	fmt.Println("Sending STOP...")
	msgChannel <- "STOP"
	<-msgChannel
	close(msgChannel)
}

func receiver(messages chan string) {
	msg := ""
	for msg != "STOP" {
		msg = <-messages
		fmt.Println("Received:", msg)
	}
	messages <- ""
}
