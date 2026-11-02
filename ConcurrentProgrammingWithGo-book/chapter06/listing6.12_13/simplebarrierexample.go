/*
 * Archivo:     simplebarrierexample.go
 *
 * Propósito:   Ilustrar una barrera con dos goroutines: "Red" trabaja 4
 *              segundos y "Blue" 10 segundos por ciclo. Al final de cada ciclo,
 *              ambas esperan en la barrera, de modo que la más rápida espera a
 *              la más lenta antes de empezar el siguiente.
 *
 * Compilación: go build simplebarrierexample.go
 * Ejecución:   ./simplebarrierexample
 *              (Alternativamente: go run simplebarrierexample.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes con el tiempo transcurrido: "<t> Red is running", "<t>
 *              Red is waiting on barrier", etc. Los dos colores empiezan cada
 *              ciclo juntos, cada 10 segundos aproximadamente. El programa dura
 *              100 segundos.
 *
 * Notas:
 *    1. "Red" pasa unos 6 segundos de cada ciclo bloqueada en la barrera
 *       esperando a "Blue".
 *    2. La barrera se reutiliza en cada ciclo porque `Wait()` reinicia el
 *       contador al liberar a todas.
 *    3. Las goroutines tienen bucles infinitos; terminan cuando acaba `main`.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter6/listing6.10"
    "time"
)

func workAndWait(name string, timeToWork int, barrier *listing6_10.Barrier) {
    start := time.Now()
    for {
        fmt.Println(time.Since(start), name, "is running")
        time.Sleep(time.Duration(timeToWork) * time.Second)
        fmt.Println(time.Since(start), name, "is waiting on barrier")
        barrier.Wait()
    }
}

func main() {
    barrier := listing6_10.NewBarrier(2)
    go workAndWait("Red", 4, barrier)
    go workAndWait("Blue", 10, barrier)
    time.Sleep(100 * time.Second)
}
