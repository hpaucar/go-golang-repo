/*
 * Archivo:     waitgrptrywait.go
 *
 * Propósito:   Solución del ejercicio 6.2 del libro *Learn Concurrent
 *              Programming with Go*, que pide agregar `TryWait()` a la
 *              implementación propia de WaitGroup basada en `sync.Cond`.
 *              `TryWait` no bloquea: devuelve `true` si el grupo ya llegó a
 *              cero y `false` en otro caso.
 *
 * Compilación: go build waitgrptrywait.go
 * Ejecución:   ./waitgrptrywait
 *              (Alternativamente: go run waitgrptrywait.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Wait group is not done. Trying again later." una vez por
 *              segundo durante unos 5 segundos, luego "Marking Wait group as
 *              done." y "Wait group is done.".
 *
 * Notas:
 *    1. `TryWait` lee `groupSize` dentro del mutex para evitar carreras.
 *    2. `main` sondea con `TryWait` en un bucle con `time.Sleep`; esto muestra
 *       la diferencia con `Wait`, que bloquea hasta que el grupo termina.
 */

package main

import (
	"fmt"
	"sync"
	"time"
)

type WaitGrp struct {
	groupSize int
	cond      *sync.Cond
}

func NewWaitGrp() *WaitGrp {
	return &WaitGrp{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (wg *WaitGrp) Add(delta int) {
	wg.cond.L.Lock()
	wg.groupSize += delta
	wg.cond.L.Unlock()
}

func (wg *WaitGrp) TryWait() bool {
	wg.cond.L.Lock()
	result := wg.groupSize == 0
	wg.cond.L.Unlock()
	return result
}

func (wg *WaitGrp) Wait() {
	wg.cond.L.Lock()
	for wg.groupSize > 0 {
		wg.cond.Wait()
	}
	wg.cond.L.Unlock()
}

func (wg *WaitGrp) Done() {
	wg.cond.L.Lock()
	wg.groupSize--
	if wg.groupSize == 0 {
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}

func main() {
	wg := NewWaitGrp()
	wg.Add(1)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Marking Wait group as done.")
		wg.Done()
	}()
	for !wg.TryWait() {
		fmt.Println("Wait group is not done. Trying again later.")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Wait group is done.")
}
