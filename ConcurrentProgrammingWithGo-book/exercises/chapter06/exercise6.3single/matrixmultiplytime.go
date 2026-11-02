/*
 * Archivo:     matrixmultiplytime.go
 *
 * Propósito:   Versión secuencial del ejercicio 6.3 del libro *Learn Concurrent
 *              Programming with Go*, que pide medir el tiempo de multiplicar
 *              matrices de 1200x1200 sin concurrencia, para compararlo con la
 *              versión con barreras de `../exercise6.3multi`.
 *
 * Compilación: go build matrixmultiplytime.go
 * Ejecución:   ./matrixmultiplytime
 *              (Alternativamente: go run matrixmultiplytime.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Complete in <duración>" con el tiempo total de generar y
 *              multiplicar 4 pares de matrices aleatorias.
 *
 * Notas:
 *    1. Las tres matrices son arreglos fijos de 1200x1200 enteros (unos 11 MB
 *       cada una) declarados como variables locales.
 *    2. En una máquina con varios núcleos, la versión concurrente debería ser
 *       varias veces más rápida.
 */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const matrixSize = 1200

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
	start := time.Now()
	for i := 0; i < 4; i++ {
		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)
		matrixMultiply(&matrixA, &matrixB, &result)
	}
	fmt.Println("Complete in", time.Since(start))
}
