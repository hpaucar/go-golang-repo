/*
 * Archivo:     countdown.go
 *
 * Propósito:   Código de partida del ejercicio 4.1 del libro *Learn Concurrent
 *              Programming with Go*, que pide encontrar y corregir con un mutex
 *              la condición de carrera del programa. Una goroutine decrementa
 *              `count` cada segundo mientras `main` lo lee e imprime cada medio
 *              segundo, sin ninguna sincronización.
 *
 * Compilación: go build countdown.go
 * Ejecución:   ./countdown
 *              (Alternativamente: go run countdown.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Una cuenta regresiva de 5 a 1, con cada valor impreso unas dos
 *              veces. Aparentemente funciona, pero tiene una carrera de datos
 *              sobre `count`.
 *
 * Notas:
 *    1. `go run -race countdown.go` reporta la carrera (DATA RACE) entre la
 *       escritura en `countdown` y la lectura en `main`.
 *    2. La solución está en `../exercise4.1solution/countdown.go`.
 */

package main

import (
    "fmt"
    "time"
)

func countdown(seconds *int) {
    for *seconds > 0 {
        time.Sleep(1 * time.Second)
        *seconds -= 1
    }
}

func main() {
    count := 5
    go countdown(&count)
    for count > 0 {
        time.Sleep(500 * time.Millisecond)
        fmt.Println(count)
    }
}

