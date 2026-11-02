/*
 * Archivo:     tryreadwritemutex.go
 *
 * Propósito:   Solución del ejercicio 4.3 del libro *Learn Concurrent
 *              Programming with Go*, que pide agregar `TryReadLock()` al mutex
 *              de lectores-escritor. El método intenta tomar `readersLock` y,
 *              si es el primer lector, también `globalLock`, sin bloquearse.
 *              `main` prueba ambos métodos.
 *
 * Compilación: go build tryreadwritemutex.go
 * Ejecución:   ./tryreadwritemutex
 *              (Alternativamente: go run tryreadwritemutex.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes que muestran el resultado de cada intento: con dos
 *              lectores activos, `TryReadLock` da `true` y `TryWriteLock` da
 *              `false`; tras liberar los tres bloqueos de lectura,
 *              `TryWriteLock` da `true` y `TryReadLock` da `false`.
 *
 * Notas:
 *    1. `main` llama tres veces a `ReadUnlock` porque el `TryReadLock` exitoso
 *       también cuenta como lector.
 *    2. Si `readersLock` está ocupado en ese momento, `TryReadLock` devuelve
 *       `false` aunque no haya un escritor; es un fallo "espurio" aceptable en
 *       una operación Try.
 */

package main

import (
    "fmt"
    "sync"
)

type ReadWriteMutex struct {
    readersCounter int
    readersLock    sync.Mutex
    globalLock     sync.Mutex
}

func (rw *ReadWriteMutex) ReadLock() {
    rw.readersLock.Lock()
    rw.readersCounter++
    if rw.readersCounter == 1 {
        rw.globalLock.Lock()
    }
    rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) TryReadLock() bool {
    if rw.readersLock.TryLock() {
        globalSuccess := true
        if rw.readersCounter == 0 {
            globalSuccess = rw.globalLock.TryLock()
        }
        if globalSuccess {
            rw.readersCounter++
        }
        rw.readersLock.Unlock()
        return globalSuccess
    } else {
        return false
    }
}

func (rw *ReadWriteMutex) WriteLock() {
    rw.globalLock.Lock()
}

func (rw *ReadWriteMutex) TryWriteLock() bool {
    return rw.globalLock.TryLock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
    rw.readersLock.Lock()
    rw.readersCounter--
    if rw.readersCounter == 0 {
        rw.globalLock.Unlock()
    }
    rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
    rw.globalLock.Unlock()
}

func main() {
    rwMutex := ReadWriteMutex{}
    fmt.Println("Acquiring Readlock")
    rwMutex.ReadLock()
    fmt.Println("Acquiring Readlock again")
    rwMutex.ReadLock()
    fmt.Println("Trying Readlock", rwMutex.TryReadLock())
    fmt.Println("Trying Writelock", rwMutex.TryWriteLock())
    rwMutex.ReadUnlock()
    rwMutex.ReadUnlock()
    rwMutex.ReadUnlock()
    fmt.Println("Trying Writelock", rwMutex.TryWriteLock())
    fmt.Println("Trying Readlock", rwMutex.TryReadLock())
}
