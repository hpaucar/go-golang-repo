/*
 * Archivo:     closingchannel.go
 *
 * Propósito:   Ilustrar el uso de un canal `quit` para que un consumidor avise
 *              al productor que debe detenerse. `main` genera números
 *              triangulares (1, 3, 6, 10, ...) y la goroutine imprime 10 y
 *              luego cierra `quit`.
 *
 * Compilación: go build closingchannel.go
 * Ejecución:   ./closingchannel
 *              (Alternativamente: go run closingchannel.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los 10 primeros números triangulares (1 a 55), uno por línea, y
 *              "Quitting number generation".
 *
 * Notas:
 *    1. Cerrar un canal es un aviso de difusión (broadcast): toda lectura de
 *       `quit` queda lista al instante.
 *    2. El `select` de `main` envía un número o detecta el cierre de `quit`,
 *       lo que ocurra primero.
 *    3. Así el consumidor controla cuándo termina el productor, sin
 *       variables compartidas.
 */

package main

import "fmt"

func printNumbers(numbers <-chan int, quit chan int) {
    go func() {
        for i := 0; i < 10; i++ {
            fmt.Println(<-numbers)
        }
        close(quit)
    }()
}

func main() {
    numbers := make(chan int)
    quit := make(chan int)
    printNumbers(numbers, quit)
    next := 0
    for i := 1; ; i++ {
        next += i
        select {
        case numbers <- next:
        case <-quit:
            fmt.Println("Quitting number generation")
            return
        }
    }
}
