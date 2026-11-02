/*
 * Archivo:     spinsemaphoremain.go
 *
 * Propósito:   Programa de prueba del ejercicio 12.3 del libro *Learn
 *              Concurrent Programming with Go*: diez goroutines compiten por un
 *              `SpinSemaphore` de 2 permisos; cada una lo toma, espera 2
 *              segundos y lo libera. `main` espera a todas con un
 *              `sync.WaitGroup`.
 *
 * Compilación: go build spinsemaphoremain.go
 * Ejecución:   ./spinsemaphoremain
 *              (Alternativamente: go run spinsemaphoremain.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes "N has acquired the semaphore" y "N releasing the
 *              semaphore", nunca con más de 2 goroutines dentro a la vez.
 *              Termina en unos 10 segundos.
 *
 * Notas:
 *    1. Mientras esperan, las goroutines giran en el bucle de `Acquire` y
 *       consumen CPU; se nota con `top`.
 *    2. El orden en que las goroutines obtienen el semáforo no está
 *       garantizado.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/exercises/chapter12/exercise12.3"
    "sync"
    "time"
)

func acquireAndWait(id int, sema *exercise12_3.SpinSemaphore) {
    sema.Acquire()
    fmt.Println(id, "has acquired the semaphore")
    time.Sleep(2 * time.Second)
    fmt.Println(id, "releasing the semaphore")
    sema.Release()
}

func main() {
    sema := exercise12_3.NewSpinSemaphore(2)
    wg := sync.WaitGroup{}
    wg.Add(10)
    for i := 0; i < 10; i++ {
        go func(id int) {
            acquireAndWait(id, sema)
            wg.Done()
        }(i)
    }
    wg.Wait()
}
