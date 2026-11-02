/*
 * Archivo:     waitgrpsemamain.go
 *
 * Propósito:   Probar el `WaitGrp` basado en semáforo del listado 6.3: `main`
 *              crea un grupo de tamaño 4, lanza cuatro goroutines que llaman a
 *              `Done()` y espera con `Wait()`.
 *
 * Compilación: go build waitgrpsemamain.go
 * Ejecución:   ./waitgrpsemamain
 *              (Alternativamente: go run waitgrpsemamain.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Cuatro líneas "N Done working" en orden no determinista y luego
 *              "All complete".
 *
 * Notas:
 *    1. "All complete" siempre aparece después de las cuatro tareas, porque
 *       `Wait()` necesita las cuatro llamadas a `Done()`.
 *    2. Si se lanzaran menos goroutines que el tamaño del grupo, `main`
 *       quedaría bloqueado (deadlock).
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter6/listing6.3"
)

func doWork(id int, wg *listing6_3.WaitGrp) {
    fmt.Println(id, "Done working ")
    wg.Done()
}

func main() {
    wg := listing6_3.NewWaitGrp(4)
    for i := 1; i <= 4; i++ {
        go doWork(i, wg)
    }
    wg.Wait()
    fmt.Println("All complete")
}
