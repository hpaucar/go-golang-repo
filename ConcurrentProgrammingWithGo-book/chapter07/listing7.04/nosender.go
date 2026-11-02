/*
 * Archivo:     nosender.go
 *
 * Propósito:   Ilustrar el bloqueo del lado receptor en un canal sin búfer:
 *              `main` intenta leer del canal, pero la goroutine `sender` solo
 *              duerme 5 segundos y termina sin enviar ningún dato.
 *
 * Compilación: go build nosender.go
 * Ejecución:   ./nosender
 *              (Alternativamente: go run nosender.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Imprime "Reading message from channel...", después de 5
 *              segundos "Sender slept for 5 seconds" y luego el programa
 *              aborta con "fatal error: all goroutines are asleep - deadlock!".
 *              "Received:" nunca se imprime.
 *
 * Notas:
 *    1. Una lectura (`<-ch`) en un canal sin búfer bloquea a la goroutine
 *       hasta que otra goroutine escriba un valor en el canal.
 *    2. `main` se bloquea en `msg := <-msgChannel`, pero `sender` nunca
 *       ejecuta `messages <- ...`.
 *    3. Cuando `sender` termina, solo queda `main` bloqueada esperando un dato
 *       que nunca llegará, y el runtime de Go detecta el interbloqueo
 *       (deadlock).
 */

package main

import (
    "fmt"
    "time"
)

func main() {
    msgChannel := make(chan string)
    go sender(msgChannel)
    
    fmt.Println("Reading message from channel...")
    msg := <-msgChannel
    fmt.Println("Received:", msg)
}

func sender(messages chan string) {
    time.Sleep(5 * time.Second)
    fmt.Println("Sender slept for 5 seconds")
}