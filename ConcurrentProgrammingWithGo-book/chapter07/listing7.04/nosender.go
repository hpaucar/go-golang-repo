/*
 * Archivo:     unbuffered_channel_receiver_block.go
 *
 * Propósito:   Ilustrar el comportamiento sincrónico y bloqueante de los canales 
 *              sin búfer (unbuffered channels) desde la perspectiva del receptor. 
 *              El hilo principal intenta leer de un canal, mientras que la 
 *              goroutine emisora duerme y termina su ejecución sin enviar ningún dato.
 *
 * Compilación: go build unbuffered_channel_receiver_block.go
 * Ejecución:   ./unbuffered_channel_receiver_block
 *              (Alternativamente: go run unbuffered_channel_receiver_block.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El programa imprimirá "Reading message from channel...", luego de 
 *              5 segundos la función `sender` mostrará "Sender slept for 5 seconds", 
 *              y finalmente el runtime de Go detectará un bloqueo y arrojará el 
 *              error fatal: "fatal error: all goroutines are asleep - deadlock!".
 *
 * Notas:       
 *    1. Las operaciones de lectura en un canal sin búfer (`<-ch`) son bloqueantes: 
 *       la goroutine receptora se suspende hasta que otra goroutine escriba un 
 *       valor en el canal.
 *    2. En este ejemplo, el hilo principal (`main`) se bloquea en la línea 
 *       `msg := <-msgChannel` esperando un mensaje.
 *    3. Sin embargo, la goroutine `sender` realiza un `time.Sleep` de 5 segundos 
 *       e imprime su mensaje, pero **nunca envía** nada a `msgChannel` (`messages <- val`) 
 *       antes de terminar.
 *    4. Al finalizar el emisor, la única goroutine restante (main) sigue bloqueada 
 *       esperando un dato que nunca llegará, provocando que el runtime de Go 
 *       interrumpa la ejecución por un interbloqueo (deadlock).
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