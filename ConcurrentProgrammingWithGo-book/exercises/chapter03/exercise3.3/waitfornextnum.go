/*
 * Archivo:     waitfornextnum.go
 *
 * Propósito:   Código del ejercicio 3.3 del libro *Learn Concurrent Programming
 *              with Go*, que pide encontrar la condición de carrera del
 *              programa. Cien goroutines comparten el arreglo `nextNum`: cada
 *              una escribe en la primera posición libre el valor anterior más
 *              uno, y `main` sondea hasta llenar la última.
 *
 * Compilación: go build waitfornextnum.go
 * Ejecución:   ./waitfornextnum
 *              (Alternativamente: go run waitfornextnum.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Varias líneas "Waiting for goroutines to complete" y luego el
 *              arreglo con los números del 1 al 101. Si ocurre la carrera, hay
 *              valores repetidos y la última posición puede no llenarse nunca,
 *              por lo que el programa no termina.
 *
 * Notas:
 *    1. La carrera está en leer y escribir `nextNum` sin sincronización: dos
 *       goroutines pueden ver la misma posición libre y escribir ambas en ella.
 *    2. Puede ser necesario ejecutarlo varias veces para observar el error; `go
 *       run -race waitfornextnum.go` lo detecta de inmediato.
 *    3. El sondeo con `time.Sleep` en `main` es una espera activa poco
 *       eficiente; los capítulos siguientes muestran alternativas.
 */

package main

import (
	"fmt"
	"time"
)

/*
  Note: this program has a race condition for demonstration purposes
  You might need to run this multiple times to trigger the race condition
*/

func addNextNumber(nextNum *[101]int) {
	i := 0
	for nextNum[i] != 0 {
		i++
	}
	nextNum[i] = nextNum[i-1] + 1
}

func main() {
	nextNum := [101]int{1}
	for i := 0; i < 100; i++ {
		go addNextNumber(&nextNum)
	}
	for nextNum[100] == 0 {
		fmt.Println("Waiting for goroutines to complete")
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println(nextNum)
}
