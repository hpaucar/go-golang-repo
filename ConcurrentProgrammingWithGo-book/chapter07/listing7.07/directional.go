/*
 * Archivo:     directional.go
 *
 * Propósito:   Ilustrar los canales direccionales (`<-chan int` solo para
 *              recibir y `chan<- int` solo para enviar) en un patrón
 *              productor-consumidor con un canal sin búfer. `sender` envía
 *              enteros sin fin cada segundo y `receiver` los imprime.
 *
 * Compilación: go build directional.go
 * Ejecución:   ./directional
 *              (Alternativamente: go run directional.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Líneas con hora (HH:MM:SS) que alternan "Sending: i" y
 *              "Received: i" durante unos 5 segundos; luego el programa
 *              termina.
 *
 * Notas:
 *    1. Los tipos direccionales (`<-chan int` y `chan<- int`) se verifican en
 *       tiempo de compilación: el emisor solo puede escribir y el receptor
 *       solo puede leer.
 *    2. Como el canal no tiene búfer, cada envío espera a que el receptor tome
 *       el valor (sincronización paso a paso).
 *    3. `time.Sleep(5 * time.Second)` en `main` decide la duración: al
 *       terminar `main` termina el programa y todas sus goroutines. En código
 *       real conviene detenerlas con un canal de cancelación o `context`.
 */

package main

import (
    "fmt"
    "time"
)

func receiver(messages <-chan int) {
    for {
        msg := <-messages
        fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
    }
}

func sender(messages chan<- int) {
    for i := 1; ; i++ {
        fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
        messages <- i
        time.Sleep(1 * time.Second)
    }
}

func main() {
    msgChannel := make(chan int)
    go receiver(msgChannel)
    go sender(msgChannel)
    time.Sleep(5 * time.Second)
}