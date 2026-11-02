/*
 * Archivo:     blockingnils.go
 *
 * Propósito:   Ilustrar que enviar a un canal `nil` bloquea para siempre.
 *
 * Compilación: go build blockingnils.go
 * Ejecución:   ./blockingnils
 *              (Alternativamente: go run blockingnils.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Nada útil: el programa aborta con "fatal error: all goroutines
 *              are asleep - deadlock!" (goroutine 1 [chan send (nil chan)]).
 *              "This is never printed" nunca se imprime.
 *
 * Notas:
 *    1. Enviar o recibir en un canal `nil` bloquea siempre; no causa pánico.
 *    2. Aquí `main` es la única goroutine, por eso el runtime detecta el
 *       interbloqueo (deadlock).
 *    3. Esta propiedad se aprovecha en `select` para desactivar casos (ver
 *       listing8.11_12).
 */

package main

import "fmt"

func main() {
    var ch chan string = nil
    ch <- "message"
    fmt.Println("This is never printed")
}
