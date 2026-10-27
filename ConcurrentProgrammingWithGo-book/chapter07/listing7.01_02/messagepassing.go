/*
 * Archivo:     channel_sentinel_string.go
 *
 * Propósito:   Ilustrar un patrón productor-consumidor clásico utilizando 
 *              un canal sin búfer y una cadena de texto como valor centinela 
 *              ("STOP") para controlar la finalización controlada de la goroutine 
 *              receptora.
 *
 * Compilación: go build channel_sentinel_string.go
 * Ejecución:   ./channel_sentinel_string
 *              (Alternativamente: go run channel_sentinel_string.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El hilo principal envía secuencialmente "HELLO", "THERE" y 
 *              "STOP", los cuales son recibidos e impresos por el receptor. 
 *              Al cumplirse la condición del centinela, el receptor sale de 
 *              su bucle y finaliza. El último envío ("STOP2") nunca llega a 
 *              leerse.
 *
 * Notas:       
 *    1. Dado que se trata de un canal sin búfer, cada operación de envío 
 *       bloquea al hilo principal hasta que la goroutine receptora lee el valor, 
 *       garantizando un orden estricto de paso de mensajes.
 *    2. La función `receiver` evalúa la condición `msg != "STOP"` antes de 
 *       cada lectura. Una vez que procesa el mensaje "STOP", el bucle termina 
 *       y la función concluye su ejecución.
 *    3. Cualquier valor enviado al canal *después* de que el receptor haya 
 *       salido de su bucle (como "STOP2" en este código) provocará que el hilo 
 *       principal se quede bloqueado indefinidamente, ya que no habrá ningún 
 *       receptor activo esperando dicho dato.
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