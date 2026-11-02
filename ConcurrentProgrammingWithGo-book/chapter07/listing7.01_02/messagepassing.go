/*
 * Archivo:     messagepassing.go
 *
 * Propósito:   Ilustrar el paso de mensajes (message passing) entre goroutines
 *              con un canal sin búfer. La goroutine `receiver` lee e imprime
 *              mensajes hasta recibir el valor centinela "STOP".
 *
 * Compilación: go build messagepassing.go
 * Ejecución:   ./messagepassing
 *              (Alternativamente: go run messagepassing.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Imprime "Sending HELLO...", "Sending THERE..." y "Sending
 *              STOP..." junto con los mensajes recibidos. Después de "STOP"
 *              el receptor termina, el envío final de "STOP2" no tiene
 *              receptor y el programa aborta con "fatal error: all goroutines
 *              are asleep - deadlock!".
 *
 * Notas:
 *    1. En un canal sin búfer, cada envío bloquea a `main` hasta que la
 *       goroutine receptora lee el valor; así se garantiza el orden de los
 *       mensajes.
 *    2. `receiver` sale de su bucle cuando lee "STOP" (valor centinela).
 *    3. El envío de "STOP2" se hace sin ningún receptor activo: `main` queda
 *       bloqueada para siempre y el runtime de Go detecta el interbloqueo
 *       (deadlock). El ejemplo muestra este error a propósito.
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
    
    msgChannel <- "STOP2"
}

func receiver(messages chan string) {
    msg := ""
    for msg != "STOP" {
        msg = <-messages
        fmt.Println("Received:", msg)
    }
}