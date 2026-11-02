/*
 * Archivo:     selectsender.go
 *
 * Propósito:   Ilustrar `select` con un caso de envío y uno de recepción a la
 *              vez. `main` envía números aleatorios a un filtro de primos
 *              (`primesOnly`) y, en el mismo `select`, recibe los primos
 *              encontrados, hasta juntar 100.
 *
 * Compilación: go build selectsender.go
 * Ejecución:   ./selectsender
 *              (Alternativamente: go run selectsender.go)
 *
 * Entrada:     Ninguna (números aleatorios entre 1 y 1 000 000 000).
 * Salida:      100 líneas "Found prime: p" con números primos aleatorios.
 *
 * Notas:
 *    1. Un `select` puede mezclar envíos (`numbersChannel <- x`) y
 *       recepciones (`<-primes`); se ejecuta el que esté listo primero.
 *    2. Si `main` solo enviara, se bloquearía cuando el filtro espera a que
 *       alguien lea su resultado (interbloqueo); el `select` evita ese
 *       problema.
 *    3. La expresión `rand.Intn(...)` se evalúa en cada vuelta aunque gane el
 *       caso de recepción; ese número simplemente se descarta.
 */

package main

import (
    "fmt"
    "math"
    "math/rand"
)

func primesOnly(inputs <-chan int) <-chan int {
    results := make(chan int)
    go func() {
        for c := range inputs {
            isPrime := c != 1
            for i := 2; i <= int(math.Sqrt(float64(c))); i++ {
                if c%i == 0 {
                    isPrime = false
                    break
                }
            }
            if isPrime {
                results <- c
            }
        }
    }()
    return results
}

func main() {
    numbersChannel := make(chan int)
    primes := primesOnly(numbersChannel)
    for i := 0; i < 100; {
        select {
        case numbersChannel <- rand.Intn(1000000000) + 1:
        case p := <-primes:
            fmt.Println("Found prime:", p)
            i++
        }
    }
}
