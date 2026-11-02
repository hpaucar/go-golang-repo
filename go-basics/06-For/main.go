/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar las formas del bucle `for`, el único bucle de Go: el
 *              `for` clásico con inicialización, condición e incremento; `for i
 *              := range 5` sobre un entero; el bucle infinito `for {}` que
 *              termina con `break`; y el uso de `continue`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los números del 0 al 4 dos veces, el texto "Bucle infinito" una
 *              sola vez y los valores impares 1, 3 y 5.
 *
 * Notas:
 *    1. Go no tiene `while`; las formas `for condicion {}` y `for {}` cumplen
 *       ese papel.
 *    2. `range` sobre un entero (`range 5`) recorre 0, 1, 2, 3 y 4, y requiere
 *       Go 1.22 o posterior.
 *    3. Los títulos impresos pueden confundir: el bloque llamado "Bucle for
 *       tradicional" en realidad usa `range`, y el clásico es el primero.
 */

package main

import "fmt"

func main() {
	// Usando un bucle for con una condición
	fmt.Println("Bucle for con condición:")
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// Usando un rango para iterar
	fmt.Println("Bucle for tradicional:")
	for i := range 5 {
		fmt.Println(i)
	}

	// Usando un bucle for con una condición de parada
	fmt.Println("Bucle for con condición de parada:")
	for {
		fmt.Println("Bucle infinito")
		break // Rompe el bucle infinito
	}

	// Usando un bucle for con una condición de parada y continue
	fmt.Println("Bucle for con continue:")
	for valor := range 6 {
		if valor%2 == 0 {
			continue
		}
		fmt.Println("Valor impar:", valor)
	}
}
