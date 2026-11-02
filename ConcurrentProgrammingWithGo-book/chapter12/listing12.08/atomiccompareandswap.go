/*
 * Archivo:     atomiccompareandswap.go
 *
 * Propósito:   Mostrar cómo funciona `atomic.CompareAndSwapInt32` (CAS):
 *              cambia el valor solo si el valor actual es igual al esperado
 *              y devuelve si el cambio se hizo. Se prueba una vez con el
 *              valor esperado (17) y otra con un valor distinto (23).
 *
 * Compilación: go build atomiccompareandswap.go
 * Ejecución:   ./atomiccompareandswap
 *              (Alternativamente: go run atomiccompareandswap.go)
 *
 * Entrada:     Ninguna.
 * Salida:      17 <- swap(17,19): result true, value: 19
 *              23 <- swap(17,19): result false, value: 23
 *
 * Notas:
 *    1. La comparación y el cambio ocurren como una sola operación
 *       indivisible, sin que otra goroutine pueda intervenir en medio.
 *    2. CAS es la base de muchas estructuras sin bloqueo (lock-free) y del
 *       spin lock de 12.09.
 */

package main

import (
    "fmt"
    "sync/atomic"
)

func main() {
    number := int32(17)
    result := atomic.CompareAndSwapInt32(&number, 17, 19)
    fmt.Printf("17 <- swap(17,19): result %t, value: %d\n", result, number)
    number = int32(23)
    result = atomic.CompareAndSwapInt32(&number, 17, 19)
    fmt.Printf("23 <- swap(17,19): result %t, value: %d\n", result, number)
}
