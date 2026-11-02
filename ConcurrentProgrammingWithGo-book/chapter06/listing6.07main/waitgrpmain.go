/*
 * Archivo:     waitgrpmain.go
 *
 * Propósito:   Probar el `WaitGrp` basado en `sync.Cond` del listado 6.7: en
 *              cada una de cuatro iteraciones, `main` llama a `wg.Add(2)` y
 *              lanza dos goroutines; después espera con `wg.Wait()` a las ocho.
 *
 * Compilación: go build waitgrpmain.go
 * Ejecución:   ./waitgrpmain
 *              (Alternativamente: go run waitgrpmain.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Ocho líneas "N Done working" (cada id de 1 a 4 aparece dos
 *              veces) en orden no determinista y luego "All complete".
 *
 * Notas:
 *    1. `Add()` se llama antes de lanzar las goroutines de cada iteración,
 *       igual que con `sync.WaitGroup`.
 *    2. Si `Add()` se llamara dentro de las goroutines, `Wait()` podría
 *       terminar antes de que todas se registren.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter6/listing6.7"
)

func doWork(id int, wg *listing6_7.WaitGrp) {
    fmt.Println(id, "Done working ")
    wg.Done()
}

func main() {
    wg := listing6_7.NewWaitGrp()
    for i := 1; i <= 4; i++ {
        wg.Add(2)
        go doWork(i, wg)
        go doWork(i, wg)
    }
    wg.Wait()
    fmt.Println("All complete")
}
