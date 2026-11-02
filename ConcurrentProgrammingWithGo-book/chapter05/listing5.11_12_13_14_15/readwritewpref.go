/*
 * Archivo:     readwritewpref.go
 *
 * Propósito:   Implementar un mutex de lectura-escritura con preferencia de
 *              escritura usando `sync.Cond`: cuando un escritor espera, los
 *              nuevos lectores se bloquean. Se prueba con el mismo escenario
 *              del listado 5.10 (dos lectores en bucle infinito y un escritor).
 *
 * Compilación: go build readwritewpref.go
 * Ejecución:   ./readwritewpref
 *              (Alternativamente: go run readwritewpref.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Algunas líneas "Read done" y luego "Write finished"; el programa
 *              termina cuando el escritor obtiene el candado.
 *
 * Notas:
 *    1. `writersWaiting` y `writerActive` forman el estado compartido; todas
 *       las comprobaciones se hacen en bucles `for` alrededor de `Wait()`.
 *    2. `ReadLock()` espera mientras haya escritores esperando o activos, por
 *       eso los lectores ya no pueden matar de hambre al escritor.
 *    3. Se usa `Broadcast()` porque pueden esperar lectores y escritores a la
 *       vez y todos deben revisar su condición.
 *    4. Con esta política pueden sufrir inanición los lectores si llegan
 *       escritores sin parar.
 */

package main

import (
    "fmt"
    "sync"
    "time"
)

type ReadWriteMutex struct {
    readersCounter int
    writersWaiting int
    writerActive   bool
    cond           *sync.Cond
}

func NewReadWriteMutex() *ReadWriteMutex {
    return &ReadWriteMutex{cond: sync.NewCond(&sync.Mutex{})}
}

func (rw *ReadWriteMutex) ReadLock() {
    rw.cond.L.Lock()
    for rw.writersWaiting > 0 || rw.writerActive {
        rw.cond.Wait()
    }
    rw.readersCounter++
    rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
    rw.cond.L.Lock()
    rw.writersWaiting++
    for rw.readersCounter > 0 || rw.writerActive {
        rw.cond.Wait()
    }
    rw.writersWaiting--
    rw.writerActive = true
    rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
    rw.cond.L.Lock()
    rw.readersCounter--
    if rw.readersCounter == 0 {
        rw.cond.Broadcast()
    }
    rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
    rw.cond.L.Lock()
    rw.writerActive = false
    rw.cond.Broadcast()
    rw.cond.L.Unlock()
}

func main() {
    rwMutex := NewReadWriteMutex()
    for i := 0; i < 2; i++ {
        go func() {
            for {
                rwMutex.ReadLock()
                time.Sleep(1 * time.Second)
                fmt.Println("Read done")
                rwMutex.ReadUnlock()
            }
        }()
    }
    time.Sleep(1 * time.Second)
    rwMutex.WriteLock()
    fmt.Println("Write finished")
}
