/*
 * Archivo:     forchannel.go
 *
 * Propósito:   Ilustrar cómo recorrer un canal con `for msg := range messages`.
 *              El bucle recibe valores hasta que el canal se cierra y entonces
 *              termina solo, sin comprobar el cierre a mano.
 *
 * Compilación: go build forchannel.go
 * Ejecución:   ./forchannel
 *              (Alternativamente: go run forchannel.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Líneas con hora (HH:MM:SS) con "Sending: i" y "Received: i"
 *              para 1, 2 y 3; después del cierre del canal se imprime
 *              "Receiver finished.".
 *
 * Notas:
 *    1. `for ... range` sobre un canal lee hasta que el canal está cerrado y
 *       vacío; equivale al patrón comma-ok de listing7.10.
 *    2. Debe cerrar el canal el emisor (`close(msgChannel)` en `main`); si
 *       nadie lo cierra, el bucle `range` espera para siempre.
 *    3. El `time.Sleep(3 * time.Second)` final da tiempo a que `receiver`
 *       imprima "Receiver finished." antes de que termine `main`.
 */

package main

import (
    "fmt"
    "time"
)

func receiver(messages <-chan int) {
    for msg := range messages {
        fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
        time.Sleep(1 * time.Second)
    }
    fmt.Println("Receiver finished.")
}

func main() {
    msgChannel := make(chan int)
    go receiver(msgChannel)
    for i := 1; i <= 3; i++ {
        fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
        msgChannel <- i
        time.Sleep(1 * time.Second)
    }
    close(msgChannel)
    time.Sleep(3 * time.Second)
}
