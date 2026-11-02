/*
 * Archivo:     matrixmultiplywg.go
 *
 * Propósito:   Solución del ejercicio 6.4 del libro *Learn Concurrent
 *              Programming with Go*, que pide reemplazar la barrera de la
 *              multiplicación de matrices por un `sync.WaitGroup`: en cada
 *              ronda se lanza una goroutine por fila y `main` espera con
 *              `wg.Wait()` antes de imprimir.
 *
 * Compilación: go build matrixmultiplywg.go
 * Ejecución:   ./matrixmultiplywg
 *              (Alternativamente: go run matrixmultiplywg.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Cuatro bloques (uno por ronda) de 3 líneas; cada línea muestra
 *              una fila de A, la fila correspondiente de B y la fila del
 *              resultado.
 *
 * Notas:
 *    1. A diferencia de la barrera, aquí se crean goroutines nuevas en cada
 *       ronda; el WaitGroup se crea de nuevo en cada iteración.
 *    2. Se usa `matrixSize = 3` para que la salida sea legible.
 */

package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const matrixSize = 3

func generateRandMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			matrix[row][col] = rand.Intn(10) - 5
		}
	}
}

func rowMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int, row int, wg *sync.WaitGroup) {
	for col := 0; col < matrixSize; col++ {
		sum := 0
		for i := 0; i < matrixSize; i++ {
			sum += matrixA[row][i] * matrixB[i][col]
		}
		result[row][col] = sum
	}
	wg.Done()
}

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int
	for i := 0; i < 4; i++ {
		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)
		wg := sync.WaitGroup{}
		wg.Add(matrixSize)
		for row := 0; row < matrixSize; row++ {
			go rowMultiply(&matrixA, &matrixB, &result, row, &wg)
		}
		wg.Wait()
		for i := 0; i < matrixSize; i++ {
			fmt.Println(matrixA[i], matrixB[i], result[i])
		}
		fmt.Println()
	}
}
