/*
 * Archivo:     noreceiver.go
 *
 * Propósito:   Ilustrar qué ocurre al enviar por un canal sin búfer
 *              (unbuffered channel) cuando nadie lee de él: la goroutine
 *              `receiver` solo duerme 5 segundos y nunca lee el canal.
 *
 * Compilación: go build noreceiver.go
 * Ejecución:   ./noreceiver
 *              (Alternativamente: go run noreceiver.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Imprime "Sending HELLO...", después de 5 segundos "Receiver
 *              slept for 5 seconds" y luego el programa aborta con "fatal
 *              error: all goroutines are asleep - deadlock!".
 *
 * Notas:
 *    1. Un envío (`ch <- val`) en un canal sin búfer bloquea a la goroutine
 *       emisora hasta que otra goroutine haga la recepción (`<-ch`).
 *    2. `receiver` nunca ejecuta `<-messages`, así que `main` se queda
 *       bloqueada en `msgChannel <- "HELLO"`.
 *    3. Cuando `receiver` termina, todas las goroutines están bloqueadas y el
 *       runtime de Go detecta el interbloqueo (deadlock) y detiene el
 *       programa.
 */

package main

import (
    "fmt"
    "time"
)

func main() {
    msgChannel := make(chan string)
    go receiver(msgChannel)
    
    fmt.Println("Sending HELLO...")
    msgChannel <- "HELLO"
    
    fmt.Println("Sending THERE...")
    msgChannel <- "THERE"
    
    fmt.Println("Sending STOP...")
    msgChannel <- "STOP"
}

func receiver(messages chan string) {
    time.Sleep(5 * time.Second)
    fmt.Println("Receiver slept for 5 seconds")
}
