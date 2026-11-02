/*
 * Archivo:     secuencial.go
 *
 * Propósito:   Mostrar un algoritmo secuencial como punto de comparación con
 *              los ejemplos concurrentes. Las instrucciones se ejecutan una
 *              después de otra, en el orden en que aparecen: `n` recibe primero
 *              el valor de `k1` y luego el de `k2`.
 *
 * Compilación: go build secuencial.go
 * Ejecución:   ./secuencial
 *              (Alternativamente: go run secuencial.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "El resultado de n es 2", siempre.
 *
 * Notas:
 *    1. En un programa secuencial el resultado es determinista: con los mismos
 *       datos, siempre se obtiene el mismo valor.
 *    2. Compárelo con `concurrente1.go`, que hace las mismas asignaciones en
 *       dos goroutines; allí el resultado depende del orden de ejecución.
 */

package main

import "fmt"

func main() {
	// Se declaran e inicializan tres variables locales.
	n := 0
	k1 := 1
	k2 := 2

	// Primera asignación:
	// el valor de k1 se copia en n.
	// Después de esta instrucción, n vale 1.
	n = k1

	// Segunda asignación:
	// el valor anterior de n se reemplaza por el valor de k2.
	// Después de esta instrucción, n vale 2.
	n = k2

	// Como la ejecución es secuencial, el resultado siempre será 2.
	fmt.Printf("El resultado de n es %d\n", n)
}