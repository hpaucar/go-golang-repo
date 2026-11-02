/*
 * Archivo:     simplenodeadlockorder.go
 *
 * Propósito:   Corregir el deadlock de 11.01_02 con un orden fijo de
 *              bloqueo: ahora `red` y `blue` adquieren siempre `lock1`
 *              primero y `lock2` después, así no puede haber espera circular.
 *
 * Compilación: go build simplenodeadlockorder.go
 * Ejecución:   ./simplenodeadlockorder
 *              (Alternativamente: go run simplenodeadlockorder.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes de `Red` y `Blue` que adquieren y liberan los locks
 *              sin detenerse durante 20 segundos; después `main` imprime
 *              "Done" y el programa termina.
 *
 * Notas:
 *    1. Ordenar los recursos y bloquearlos siempre en ese orden elimina la
 *       condición de espera circular de Coffman.
 *    2. El programa genera muchísima salida; conviene redirigirla, por
 *       ejemplo con `| tail`, para ver el final.
 */

package main

import (
    "fmt"
    "sync"
    "time"
)

func red(lock1, lock2 *sync.Mutex) {
    for {
        fmt.Println("Red: Acquiring lock1")
        lock1.Lock()
        fmt.Println("Red: Acquiring lock2")
        lock2.Lock()
        fmt.Println("Red: Both locks Acquired")
        lock1.Unlock(); lock2.Unlock()
        fmt.Println("Red: Locks Released")
    }
}

func blue(lock1, lock2 *sync.Mutex) {
    for {
        fmt.Println("Blue: Acquiring lock1")
        lock1.Lock()
        fmt.Println("Blue: Acquiring lock2")
        lock2.Lock()
        fmt.Println("Blue: Both locks Acquired")
        lock1.Unlock(); lock2.Unlock()
        fmt.Println("Blue: Locks Released")
    }
}

func main() {
    lockA := sync.Mutex{}
    lockB := sync.Mutex{}
    go red(&lockA, &lockB)
    go blue(&lockA, &lockB)
    time.Sleep(20 * time.Second)
    fmt.Println("Done")
}