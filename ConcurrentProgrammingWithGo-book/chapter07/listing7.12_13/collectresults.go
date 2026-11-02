/*
 * Archivo:     collectresults.go
 *
 * Propósito:   Ilustrar cómo obtener el resultado de una goroutine mediante un
 *              canal: una goroutine calcula los divisores de 3419110721
 *              mientras `main` calcula los de 4033836233 al mismo tiempo.
 *
 * Compilación: go build collectresults.go
 * Ejecución:   ./collectresults
 *              (Alternativamente: go run collectresults.go)
 *
 * Entrada:     Ninguna (los números están fijos en el código).
 * Salida:      Dos líneas con los divisores de 4033836233 y de 3419110721
 *              (en ese orden). Puede tardar varios segundos.
 *
 * Notas:
 *    1. `findFactors` busca divisores por fuerza bruta, de 1 a n: es una
 *       tarea intensiva en CPU (CPU-bound).
 *    2. El canal sin búfer `resultCh` lleva el resultado de la goroutine a
 *       `main`; `<-resultCh` espera hasta que el cálculo termine.
 *    3. Los dos cálculos corren en paralelo si hay más de un núcleo, así que
 *       el tiempo total es cercano al del cálculo más largo, no a la suma.
 */

package main

import (
    "fmt"
)

func findFactors(number int) []int {
    result := make([]int, 0)
    for i := 1; i <= number; i++ {
        if number%i == 0 {
            result = append(result, i)
        }
    }
    return result
}

func main() {
    resultCh := make(chan []int)
    
    go func() {
        resultCh <- findFactors(3419110721)
    }()
    
    fmt.Println(findFactors(4033836233))
    fmt.Println(<-resultCh)
}