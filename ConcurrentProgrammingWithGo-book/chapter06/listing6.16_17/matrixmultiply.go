/*
 * Archivo:     matrixmultiply.go
 *
 * Propósito:   Multiplicar matrices de forma concurrente con una barrera: una
 *              goroutine calcula cada fila del resultado. En cada ronda, `main`
 *              genera A y B y pasa dos veces por la barrera de tamaño
 *              `matrixSize + 1`: la primera inicia el cálculo y la segunda
 *              espera a que todas las filas estén listas.
 *
 * Compilación: go build matrixmultiply.go
 * Ejecución:   ./matrixmultiply
 *              (Alternativamente: go run matrixmultiply.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Cuatro bloques de tres líneas con una fila de A, una fila de B y
 *              la fila del resultado, igual que en la versión secuencial.
 *
 * Notas:
 *    1. La primera barrera evita que las goroutines lean las matrices mientras
 *       `main` las está generando; la segunda evita que `main` imprima un
 *       resultado incompleto.
 *    2. La barrera tiene tamaño `matrixSize + 1` porque participan las
 *       goroutines de las filas y `main`.
 *    3. Al terminar `main`, las goroutines quedan bloqueadas en la barrera y el
 *       programa finaliza con ellas.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter6/listing6.10"
    "math/rand"
)

const matrixSize = 3

func generateRandMatrix(matrix *[matrixSize][matrixSize]int) {
    for row := 0; row < matrixSize; row++ {
        for col := 0; col < matrixSize; col++ {
            matrix[row][col] = rand.Intn(10) - 5
        }
    }
}

func rowMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int,
    row int, barrier *listing6_10.Barrier) {
    for {
        barrier.Wait()
        for col := 0; col < matrixSize; col++ {
            sum := 0
            for i := 0; i < matrixSize; i++ {
                sum += matrixA[row][i] * matrixB[i][col]
            }
            result[row][col] = sum
        }
        barrier.Wait()
    }
}

func main() {
    var matrixA, matrixB, result [matrixSize][matrixSize]int
    barrier := listing6_10.NewBarrier(matrixSize + 1)
    for row := 0; row < matrixSize; row++ {
        go rowMultiply(&matrixA, &matrixB, &result, row, barrier)
    }

    for i := 0; i < 4; i++ {
        generateRandMatrix(&matrixA)
        generateRandMatrix(&matrixB)
        barrier.Wait()
        barrier.Wait()
        for i := 0; i < matrixSize; i++ {
            fmt.Println(matrixA[i], matrixB[i], result[i])
        }
        fmt.Println()
    }
}
