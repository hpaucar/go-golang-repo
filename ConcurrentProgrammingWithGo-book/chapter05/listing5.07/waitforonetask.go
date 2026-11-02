/*
 * Archivo:     waitforonetask.go
 *
 * Propósito:   Corregir la señal perdida del listado anterior: `doWork`
 *              adquiere el candado antes de llamar a `cond.Signal()`. Como
 *              `main` tiene el candado hasta que entra en `cond.Wait()`, la
 *              señal nunca llega antes de que `main` esté esperando.
 *
 * Compilación: go build waitforonetask.go
 * Ejecución:   ./waitforonetask
 *              (Alternativamente: go run waitforonetask.go)
 *
 * Entrada:     Ninguna.
 * Salida:      La secuencia "Waiting for child goroutine", "Work started",
 *              "Work finished", "Child goroutine finished" repetida 50 000
 *              veces, sin interbloqueo.
 *
 * Notas:
 *    1. `main` bloquea `cond.L` antes de lanzar la goroutine; `doWork` se queda
 *       esperando el candado.
 *    2. `cond.Wait()` libera el candado de forma atómica al suspender a `main`;
 *       solo entonces `doWork` puede adquirirlo y emitir la señal.
 *    3. Aun así, la práctica recomendada es acompañar `sync.Cond` con una
 *       variable de estado (por ejemplo, `done bool`) y comprobarla en un bucle
 *       `for` antes de `Wait()`.
 */

package main

import (
    "fmt"
    "sync"
)

func doWork(cond *sync.Cond) {
    fmt.Println("Work started")
    fmt.Println("Work finished")
    cond.L.Lock()
    cond.Signal()
    cond.L.Unlock()
}

func main() {
    cond := sync.NewCond(&sync.Mutex{})
    cond.L.Lock()
    for i := 0; i < 50000; i++ {
        go doWork(cond)
        fmt.Println("Waiting for child goroutine")
        cond.Wait()
        fmt.Println("Child goroutine finished")
    }
    cond.L.Unlock()
}
