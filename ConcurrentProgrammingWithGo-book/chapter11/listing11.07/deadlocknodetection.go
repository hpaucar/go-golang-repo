/*
 * Archivo:     deadlocknodetection.go
 *
 * Propósito:   Mostrar que el runtime de Go no detecta un deadlock parcial.
 *              Es el mismo bloqueo de 11.06, pero la espera del `WaitGroup`
 *              se hace en otra goroutine y `main` duerme 30 segundos, así
 *              que no todas las goroutines están bloqueadas.
 *
 * Compilación: go build deadlocknodetection.go
 * Ejecución:   ./deadlocknodetection
 *              (Alternativamente: go run deadlocknodetection.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Si ocurre el deadlock (lo normal), no hay ningún error: a los
 *              30 segundos solo se imprime "Done", sin el mensaje
 *              "Done waiting on waitgroup".
 *
 * Notas:
 *    1. Las goroutines bloqueadas quedan colgadas en silencio; en un
 *       servidor real esto se ve como peticiones que nunca terminan.
 *    2. No se debe confiar en la detección del runtime: hay que evitar los
 *       deadlocks en el diseño (por ejemplo, con un orden fijo de bloqueo,
 *       como en 11.13 y 11.14).
 */

package main

import (
    "fmt"
    "sync"
    "time"
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
    go func() {
        wg.Wait()
        fmt.Println("Done waiting on waitgroup")
    }()
    time.Sleep(30 * time.Second)
    fmt.Println("Done")
}
