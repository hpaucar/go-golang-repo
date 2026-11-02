/*
 * Archivo:     countdown.go
 *
 * Propósito:   Solución del ejercicio 4.1 del libro *Learn Concurrent
 *              Programming with Go*: corrige la condición de carrera del
 *              contador regresivo protegiendo cada lectura y escritura de
 *              `count` con un `sync.Mutex` compartido entre `main` y la
 *              goroutine `countdown`.
 *
 * Compilación: go build countdown.go
 * Ejecución:   ./countdown
 *              (Alternativamente: go run countdown.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Una cuenta regresiva de 5 a 0, con cada valor impreso unas dos
 *              veces. Con `go run -race` ya no se reporta ninguna carrera.
 *
 * Notas:
 *    1. El valor se copia a la variable local `remaining` dentro de la sección
 *       crítica, de modo que la condición del bucle no lee `count` sin el
 *       mutex.
 *    2. La goroutine duerme fuera de la sección crítica, así el mutex se
 *       mantiene tomado el menor tiempo posible.
 *    3. Compare con `../exercise4.1/countdown.go`, el código original con la
 *       carrera.
 */

package main

import (
    "fmt"
    "sync"
    "time"
)

func countdown(seconds *int, mutex *sync.Mutex) {
    mutex.Lock()
    remaining := *seconds
    mutex.Unlock()
    for remaining > 0 {
        time.Sleep(1 * time.Second)
        mutex.Lock()
        *seconds -= 1
        remaining = *seconds
        mutex.Unlock()
    }
}

func main() {
    mutex := sync.Mutex{}
    count := 5
    go countdown(&count, &mutex)
    remaining := count
    for remaining > 0 {
        time.Sleep(500 * time.Millisecond)
        mutex.Lock()
        fmt.Println(count)
        remaining = count
        mutex.Unlock()
    }
}

