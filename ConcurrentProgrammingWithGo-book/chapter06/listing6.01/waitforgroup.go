/*
 * Archivo:     waitforgroup.go
 *
 * Propósito:   Esperar a varias goroutines con `sync.WaitGroup`: `main`
 *              registra 4 tareas con `wg.Add(4)`, lanza cuatro goroutines que
 *              duermen un tiempo aleatorio de 0 a 4 segundos y espera con
 *              `wg.Wait()` a que todas llamen a `wg.Done()`.
 *
 * Compilación: go build waitforgroup.go
 * Ejecución:   ./waitforgroup
 *              (Alternativamente: go run waitforgroup.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Cuatro líneas "N Done working after S seconds" en el orden en
 *              que terminan y, al final, "All complete".
 *
 * Notas:
 *    1. `Add()` debe llamarse antes de lanzar las goroutines para que `Wait()`
 *       no termine antes de tiempo.
 *    2. A diferencia de `time.Sleep`, `Wait()` espera exactamente lo necesario.
 *    3. En código real se suele escribir `defer wg.Done()` para que se llame
 *       aunque la función termine con error.
 */

package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func doWork(id int, wg *sync.WaitGroup) {
    i := rand.Intn(5)
    time.Sleep(time.Duration(i) * time.Second)
    fmt.Println(id, "Done working after", i, "seconds")
    wg.Done()
}

func main() {
    wg := sync.WaitGroup{}
    wg.Add(4)
    for i := 1; i <= 4; i++ {
        go doWork(i, &wg)
    }
    wg.Wait()
    fmt.Println("All complete")
}
