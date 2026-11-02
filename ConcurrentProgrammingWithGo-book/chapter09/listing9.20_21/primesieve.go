/*
 * Archivo:     primesieve.go
 *
 * Propósito:   Ilustrar una cadena dinámica de goroutines (criba de primos
 *              concurrente). Cada `primeMultipleFilter` imprime el primer
 *              número que recibe (un primo), descarta sus múltiplos y crea la
 *              siguiente goroutine de la cadena cuando la necesita.
 *
 * Compilación: go build primesieve.go
 * Ejecución:   ./primesieve
 *              (Alternativamente: go run primesieve.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Todos los números primos menores que 100 000, uno por línea, en
 *              orden creciente.
 *
 * Notas:
 *    1. Se crea una goroutine por cada primo encontrado (9 592 en total);
 *       las goroutines de Go son lo bastante ligeras para esto.
 *    2. Al cerrar `numbers`, cada filtro cierra el canal de su derecha; el
 *       último filtro cierra `quit` y así `main` sabe que todo terminó.
 *    3. Es concurrente pero no eficiente: cada número pasa por muchos
 *       canales. Sirve para aprender, no como algoritmo rápido.
 */

package main

import "fmt"

func primeMultipleFilter(numbers <-chan int, quit chan<- int) {
    var right chan int
    p := <-numbers
    fmt.Println(p)
    for n := range numbers {
        if n%p != 0 {
            if right == nil {
                right = make(chan int)
                go primeMultipleFilter(right, quit)
            }
            right <- n
        }
    }
    if right == nil {
        close(quit)
    } else {
        close(right)
    }
}

func main() {
    numbers := make(chan int)
    quit := make(chan int)
    go primeMultipleFilter(numbers, quit)
    for i := 2; i < 100000; i++ {
        numbers <- i
    }
    close(numbers)
    <-quit
}
