/*
 * Archivo:     spinlockboolmain.go
 *
 * Propósito:   Programa de prueba del ejercicio 12.2 del libro *Learn
 *              Concurrent Programming with Go*: comprueba el `TryLock` del spin
 *              lock del paquete `exercise12_2`. Primero lo intenta con el
 *              bloqueo tomado y luego con el bloqueo libre.
 *
 * Compilación: go build spinlockboolmain.go
 * Ejecución:   ./spinlockboolmain
 *              (Alternativamente: go run spinlockboolmain.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "This should be false:  false" y "This should be true:  true".
 *
 * Notas:
 *    1. Debe ejecutarse dentro del módulo (hay un `go.mod` en la raíz del
 *       libro) para que se resuelva el import.
 *    2. El último `Unlock` libera el bloqueo obtenido por el `TryLock` exitoso.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/exercises/chapter12/exercise12.2"
)

func main() {
    spin := exercise12_2.SpinLock{}

    spin.Lock()

    fmt.Println("This should be false: ", spin.TryLock())

    spin.Unlock()

    fmt.Println("This should be true: ", spin.TryLock())

    spin.Unlock()
}
