/*
 * Archivo:     signalbeforewait.go
 *
 * Propósito:   Demostrar el problema de la señal perdida (lost signal) con
 *              variables de condición. `main` lanza una goroutine `doWork` y
 *              espera con `cond.Wait()`, pero `doWork` llama a `cond.Signal()`
 *              sin adquirir el candado, así que puede avisar antes de que
 *              `main` esté esperando.
 *
 * Compilación: go build signalbeforewait.go
 * Ejecución:   ./signalbeforewait
 *              (Alternativamente: go run signalbeforewait.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Muchas líneas "Work started", "Work finished", "Waiting for
 *              child goroutine" y "Child goroutine finished"; tarde o temprano
 *              una señal se pierde y el programa termina con "fatal error: all
 *              goroutines are asleep - deadlock!".
 *
 * Notas:
 *    1. `sync.Cond` no tiene memoria: si `Signal()` ocurre cuando nadie espera
 *       en `Wait()`, la señal se pierde para siempre.
 *    2. No hay una variable de estado que `main` pueda comprobar en un bucle
 *       `for` (por ejemplo, `for !done { cond.Wait() }`).
 *    3. El código es intencionalmente incorrecto; el listado 5.7 lo corrige.
 *       Para esperar goroutines, use `sync.WaitGroup` o canales (channels).
 */

package main

import (
    "fmt"
    "sync"
)

/*
Note: this program has a bug for demonstration purposes
We demonstrate how to fix this problem in the next listing
*/
func doWork(cond *sync.Cond) {
    fmt.Println("Work started")
    fmt.Println("Work finished")
    cond.Signal()
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
