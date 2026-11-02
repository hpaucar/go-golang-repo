/*
 * Archivo:     matrixmultiplytime.go
 *
 * Propósito:   Versión concurrente del ejercicio 6.3 del libro *Learn
 *              Concurrent Programming with Go*, que pide medir el tiempo de
 *              multiplicar matrices de 1200x1200 en forma secuencial y
 *              concurrente. Usa una goroutine por fila coordinada con la
 *              barrera (`Barrier`) del listado 6.10.
 *
 * Compilación: go build matrixmultiplytime.go
 * Ejecución:   ./matrixmultiplytime
 *              (Alternativamente: go run matrixmultiplytime.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Complete in <duración>" con el tiempo total de 4
 *              multiplicaciones. Compare con `../exercise6.3single`.
 *
 * Notas:
 *    1. Actualmente no compila: importa `.../chapter6/listing6.10`, pero la
 *       carpeta del repositorio se llama `chapter06`.
 *    2. Cada ronda usa dos `barrier.Wait()`: el primero inicia el cálculo y el
 *       segundo espera a que todas las filas terminen.
 *    3. Las 1200 goroutines de `rowMultiply` nunca terminan; quedan bloqueadas
 *       en la barrera hasta que `main` sale.
 */

package main

import (
	"fmt"
	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter6/listing6.10"
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
	start := time.Now()
	for i := 0; i < 4; i++ {
		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)
		barrier.Wait()
		barrier.Wait()
	}
	fmt.Println("Complete in", time.Since(start))
}
