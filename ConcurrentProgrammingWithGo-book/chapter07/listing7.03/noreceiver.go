/*
 * Archivo:     unbuffered_channel_deadlock.go
 *
 * Propósito:   Ilustrar el comportamiento de los canales sin búfer (unbuffered 
 *              channels) en Go y cómo se produce un interbloqueo (deadlock) 
 *              cuando se intenta enviar datos a través de un canal que nunca 
 *              es leído por el receptor.
 *
 * Compilación: go build unbuffered_channel_deadlock.go
 * Ejecución:   ./unbuffered_channel_deadlock
 *              (Alternativamente: go run unbuffered_channel_deadlock.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El programa imprimirá "Sending HELLO...", luego de 5 segundos 
 *              la goroutine receptora mostrará "Receiver slept for 5 seconds", 
 *              y finalmente el runtime de Go detectará un bloqueo y arrojará 
 *              el error fatal: "fatal error: all goroutines are asleep - deadlock!".
 *
 * Notas:       
 *    1. Los canales sin búfer en Go son estrictamente sincrónicos: una operación 
 *       de envío (`ch <- val`) bloquea a la goroutine emisora hasta que otra 
 *       goroutine realiza una operación de recepción (`<-ch`).
 *    2. En este ejemplo, la función `receiver` duerme durante 5 segundos e 
 *       imprime un mensaje, pero **nunca lee** del canal `messages` (`<-messages`).
 *    3. Como resultado, cuando la función `main` intenta ejecutar `msgChannel <- "HELLO"`, 
 *       se bloquea indefinidamente esperando un receptor. Dado que la única otra 
 *       goroutine activa termina su ejecución poco después, el runtime de Go 
 *       detecta que todas las goroutines están dormidas sin esperanza de despertar 
 *       y aborta el programa con un error de interbloqueo.
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
