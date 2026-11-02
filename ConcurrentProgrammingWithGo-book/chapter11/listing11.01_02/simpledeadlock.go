/*
 * Archivo:     simpledeadlock.go
 *
 * Propósito:   Provocar un interbloqueo (deadlock) clásico con dos mutex.
 *              La goroutine `red` adquiere `lock1` y luego `lock2`, mientras
 *              `blue` los adquiere en orden inverso (`lock2` y luego
 *              `lock1`). Tarde o temprano cada una retiene un mutex y espera
 *              para siempre el que tiene la otra.
 *
 * Compilación: go build simpledeadlock.go
 * Ejecución:   ./simpledeadlock
 *              (Alternativamente: go run simpledeadlock.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes de `Red` y `Blue` que adquieren y liberan los locks.
 *              En algún momento los mensajes se detienen (normalmente con
 *              "Red: Acquiring lock2" y "Blue: Acquiring lock1" como últimos
 *              mensajes) y, a los 20 segundos, `main` imprime "Done".
 *
 * Notas:
 *    1. Se cumplen las cuatro condiciones de Coffman: exclusión mutua,
 *       retención y espera, no expropiación y espera circular.
 *    2. El runtime de Go no detecta este deadlock porque `main` sigue viva
 *       en `time.Sleep`; solo informa un deadlock cuando todas las
 *       goroutines están bloqueadas.
 *    3. El momento en que aparece el bloqueo cambia en cada ejecución; a
 *       veces ocurre casi de inmediato y otras tras muchas iteraciones.
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
        fmt.Println("Blue: Acquiring lock2")
        lock2.Lock()
        fmt.Println("Blue: Acquiring lock1")
        lock1.Lock()
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