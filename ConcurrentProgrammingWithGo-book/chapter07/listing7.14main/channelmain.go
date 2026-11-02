/*
 * Archivo:     channelmain.go
 *
 * Propósito:   Usar el canal propio basado en semáforos de listing7.14
 *              (`listing7_14.Channel`) en un patrón productor-consumidor con
 *              capacidad 10 y el valor centinela `-1` para terminar.
 *
 * Compilación: go build channelmain.go
 * Ejecución:   ./channelmain
 *              (Alternativamente: go run channelmain.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Sending:  1" a "Sending:  6" casi al instante (caben en el
 *              búfer) y luego "Received: 1" ... "Received: -1", uno por
 *              segundo.
 *
 * Notas:
 *    1. `listing7_14.NewChannel[int](10)` crea un canal con búfer de 10
 *       elementos, por eso `main` no se bloquea al enviar.
 *    2. `receiver` lee con `Receive()` cada segundo hasta recibir `-1`;
 *       `sync.WaitGroup` hace que `main` lo espere.
 *    3. La ruta de importación usa `chapter7/listing7.14`, pero la carpeta
 *       real es `chapter07/`; tal como está, `go build` falla con "no
 *       required module provides package".
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter7/listing7.14"
    "sync"
    "time"
)

func receiver(messages *listing7_14.Channel[int], wGroup *sync.WaitGroup) {
    msg := 0
    for msg != -1 {
        time.Sleep(1 * time.Second)
        msg = messages.Receive()
        fmt.Println("Received:", msg)
    }
    wGroup.Done()
}

func main() {
    channel := listing7_14.NewChannel[int](10)
    wGroup := sync.WaitGroup{}
    wGroup.Add(1)
    
    go receiver(channel, &wGroup)
    
    for i := 1; i <= 6; i++ {
        fmt.Println("Sending: ", i)
        channel.Send(i)
    }
    
    channel.Send(-1)
    wGroup.Wait()
}