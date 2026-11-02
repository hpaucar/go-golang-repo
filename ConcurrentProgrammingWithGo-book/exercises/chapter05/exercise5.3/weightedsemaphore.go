/*
 * Archivo:     weightedsemaphore.go
 *
 * Propósito:   Solución del ejercicio 5.3 del libro *Learn Concurrent
 *              Programming with Go*, que pide un semáforo ponderado (weighted
 *              semaphore) donde `Acquire` y `Release` toman o devuelven varios
 *              permisos a la vez. Se implementa con un contador y una
 *              `sync.Cond`.
 *
 * Compilación: go build weightedsemaphore.go
 * Ejecución:   ./weightedsemaphore
 *              (Alternativamente: go run weightedsemaphore.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El hilo padre A toma los 2 permisos; los hijos B y C quedan
 *              bloqueados 3 segundos hasta que A los libera, y luego cada uno
 *              toma, usa y libera 1 permiso.
 *
 * Notas:
 *    1. `Acquire(n)` espera mientras `permits - n < 0`, es decir, hasta que
 *       haya suficientes permisos.
 *    2. `Release` usa `Broadcast` porque varias goroutines pueden esperar
 *       cantidades distintas de permisos.
 *    3. Para uso real existe `golang.org/x/sync/semaphore.Weighted`.
 */

package main

import (
	"fmt"
	"sync"
	"time"
)

type WeightedSemaphore struct {
	permits int
	cond    *sync.Cond
}

func NewSemaphore(n int) *WeightedSemaphore {
	return &WeightedSemaphore{
		permits: n,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (rw *WeightedSemaphore) Acquire(permits int) {
	rw.cond.L.Lock()
	for rw.permits-permits < 0 {
		rw.cond.Wait()
	}
	rw.permits -= permits
	rw.cond.L.Unlock()
}

func (rw *WeightedSemaphore) Release(permits int) {
	rw.cond.L.Lock()
	rw.permits += permits
	rw.cond.Broadcast()
	rw.cond.L.Unlock()
}

func main() {
	sema := NewSemaphore(2)
	sema.Acquire(2)
	fmt.Println("Parent thread A acquired semaphore")
	go func() {
		sema.Acquire(1)
		fmt.Println("Child B thread acquired semaphore")
		time.Sleep(1 * time.Second)
		sema.Release(1)
		fmt.Println("Child B thread released semaphore")
	}()
	go func() {
		sema.Acquire(1)
		fmt.Println("Child C thread acquired semaphore")
		time.Sleep(1 * time.Second)
		sema.Release(1)
		fmt.Println("Child C thread released semaphore")
	}()
	time.Sleep(3 * time.Second)
	fmt.Println("Parent thread A releasing semaphore")
	sema.Release(2)
	time.Sleep(3 * time.Second)
}
