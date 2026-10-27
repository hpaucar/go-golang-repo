/*
 * Archivo:     custom_channel_test.go
 *
 * Propósito:   Demostrar el uso práctico del canal personalizado basado en 
 *              semáforos (`listing7_14.Channel`) que implementamos en el ejemplo 
 *              anterior. El programa simula un patrón productor-consumidor 
 *              utilizando capacidad acotada y un valor centinela (`-1`) para 
 *              finalizar la goroutine receptora de forma coordinada.
 *
 * Compilación: go build custom_channel_test.go
 * Ejecución:   (Requiere la implementación del paquete `listing7_14` y sus 
 *              dependencias del libro de J. Cutajar).
 *
 * Entrada:     Ninguna.
 * Salida:      Imprime en consola los mensajes de envío secuencial de los 
 *              enteros del 1 al 6 más el valor centinela, junto con su respectiva 
 *              recepción diferida por el hilo receptor con intervalos de 1 segundo.
 *
 * Notas:       
 *    1. `listing7_14.NewChannel[int](10)` inicializa un canal genérico seguro 
 *       para hilos con una capacidad de búfer de hasta 10 elementos.
 *    2. La función `main` actúa como emisor enviando 6 elementos seguidos del 
 *       valor centinela `-1`, mientras que la función `receiver` extrae los 
 *       elementos utilizando el método `Receive()` con un retraso de 1 segundo 
 *       por cada lectura.
 *    3. Se emplea un `sync.WaitGroup` para asegurar que el hilo principal espere 
 *       a que el receptor procese todos los datos (incluyendo el centinela) y 
 *       llame a `wGroup.Done()` antes de finalizar la ejecución del programa.
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