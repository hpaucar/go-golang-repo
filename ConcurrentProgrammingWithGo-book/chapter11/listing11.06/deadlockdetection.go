/*
 * Archivo:     deadlockdetection.go
 *
 * Propósito:   Mostrar la detección automática de deadlocks del runtime de
 *              Go. Dos goroutines bloquean los mismos dos mutex en orden
 *              inverso 10000 veces, y `main` espera con `wg.Wait()`.
 *
 * Compilación: go build deadlockdetection.go
 * Ejecución:   ./deadlockdetection
 *              (Alternativamente: go run deadlockdetection.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Casi siempre el programa termina con "fatal error: all
 *              goroutines are asleep - deadlock!" y la traza de las
 *              goroutines. Solo si no hay bloqueo por casualidad, imprime
 *              "Done".
 *
 * Notas:
 *    1. Aquí el runtime sí detecta el deadlock porque todas las goroutines,
 *       incluida `main`, quedan bloqueadas al mismo tiempo.
 *    2. Esta detección es limitada: basta con que una goroutine siga viva
 *       (ver 11.07) para que el bloqueo pase inadvertido.
 */

package main

import (
    "fmt"
    "sync"
)

func lockBoth(lock1, lock2 *sync.Mutex, wg *sync.WaitGroup) {
    for i := 0; i < 10000; i++ {
        lock1.Lock(); lock2.Lock()
        lock1.Unlock(); lock2.Unlock()
    }
    wg.Done()
}

func main() {
    lockA, lockB := sync.Mutex{}, sync.Mutex{}
    wg := sync.WaitGroup{}
    wg.Add(2)
    go lockBoth(&lockA, &lockB, &wg)
    go lockBoth(&lockB, &lockA, &wg)
    wg.Wait()
    fmt.Println("Done")
}
