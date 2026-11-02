/*
 * Archivo:     releasebeforeacquire.go
 *
 * Propósito:   Resolver el problema de la señal perdida con un semáforo: `main`
 *              lanza 50 000 goroutines, una tras otra, y espera a cada una con
 *              `Acquire()` sobre un semáforo iniciado en 0. Cada goroutine
 *              termina con `Release()`.
 *
 * Compilación: go build releasebeforeacquire.go
 * Ejecución:   ./releasebeforeacquire
 *              (Alternativamente: go run releasebeforeacquire.go)
 *
 * Entrada:     Ninguna.
 * Salida:      50 000 repeticiones de "Waiting for child goroutine", "Work
 *              started", "Work finished" y "Child goroutine finished", sin
 *              interbloqueo.
 *
 * Notas:
 *    1. Si `Release()` ocurre antes de `Acquire()`, el permiso queda guardado y
 *       `Acquire()` no se bloquea; no se pierde ninguna señal.
 *    2. Compare con el listado 5.6, donde `Signal()` antes de `Wait()` producía
 *       un deadlock.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter5/listing5.16"
)

func main() {
    semaphore := listing5_16.NewSemaphore(0)
    for i := 0; i < 50000; i++ {
        go doWork(semaphore)
        fmt.Println("Waiting for child goroutine")
        semaphore.Acquire()
        fmt.Println("Child goroutine finished")
    }
}

func doWork(semaphore *listing5_16.Semaphore) {
    fmt.Println("Work started")
    fmt.Println("Work finished")
    semaphore.Release()
}
