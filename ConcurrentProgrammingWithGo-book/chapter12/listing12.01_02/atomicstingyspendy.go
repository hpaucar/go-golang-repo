/*
 * Archivo:     atomicstingyspendy.go
 *
 * Propósito:   Resolver la condición de carrera de `stingy` y `spendy` con
 *              variables atómicas en lugar de un mutex. Una goroutine suma
 *              10 y la otra resta 10, un millón de veces cada una, sobre el
 *              mismo `int32` usando `atomic.AddInt32`.
 *
 * Compilación: go build atomicstingyspendy.go
 * Ejecución:   ./atomicstingyspendy
 *              (Alternativamente: go run atomicstingyspendy.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Stingy Done" y "Spendy Done" (en cualquier orden) y luego
 *              "Money in account:  100", siempre el mismo valor.
 *
 * Notas:
 *    1. `atomic.AddInt32` hace la lectura, la suma y la escritura como una
 *       sola operación indivisible, así que no se pierden actualizaciones.
 *    2. La lectura final también se hace con `atomic.LoadInt32`: todos los
 *       accesos a una variable compartida deben ser atómicos.
 *    3. Sirve para contadores simples; para actualizar varias variables a
 *       la vez sigue haciendo falta un mutex.
 */

package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func stingy(money *int32) {
    for i := 0; i < 1000000; i++ {
        atomic.AddInt32(money, 10)
    }
    fmt.Println("Stingy Done")
}

func spendy(money *int32) {
    for i := 0; i < 1000000; i++ {
        atomic.AddInt32(money, -10)
    }
    fmt.Println("Spendy Done")
}

func main() {
    money := int32(100)
    wg := sync.WaitGroup{}
    wg.Add(2)
    go func() {
        stingy(&money)
        wg.Done()
    }()
    go func() {
        spendy(&money)
        wg.Done()
    }()
    wg.Wait()
    fmt.Println("Money in account: ", atomic.LoadInt32(&money))
}
