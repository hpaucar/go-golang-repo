/*
 * Archivo:     matrixmultiplysimple.go
 *
 * Propósito:   Multiplicar matrices de forma secuencial: en cuatro rondas se
 *              generan dos matrices aleatorias de 3x3 (valores de -5 a 4) y se
 *              calcula su producto con tres bucles anidados. Es la versión base
 *              para comparar con la versión concurrente con barrera.
 *
 * Compilación: go build matrixmultiplysimple.go
 * Ejecución:   ./matrixmultiplysimple
 *              (Alternativamente: go run matrixmultiplysimple.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Cuatro bloques de tres líneas; cada línea muestra una fila de A,
 *              una fila de B y la fila correspondiente del resultado.
 *
 * Notas:
 *    1. La constante `matrixSize` define el tamaño de las matrices.
 *    2. No se fija la semilla de `math/rand`; según la versión de Go, las
 *       matrices pueden repetirse o cambiar entre ejecuciones.
 */

package main

import (
    "fmt"
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

func matrixMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
    for row := 0; row < matrixSize; row++ {
        for col := 0; col < matrixSize; col++ {
            sum := 0
            for i := 0; i < matrixSize; i++ {
                sum += matrixA[row][i] * matrixB[i][col]
            }
            result[row][col] = sum
        }
    }
}

func main() {
    var matrixA, matrixB, result [matrixSize][matrixSize]int
    for i := 0; i < 4; i++ {
        generateRandMatrix(&matrixA)
        generateRandMatrix(&matrixB)
        matrixMultiply(&matrixA, &matrixB, &result)
        for i := 0; i < matrixSize; i++ {
            fmt.Println(matrixA[i], matrixB[i], result[i])
        }
        fmt.Println()
    }
}
