/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar los arreglos (arrays) de tamaño fijo: declaración con
 *              valores cero, asignación por índice, `len`, inicialización con
 *              valores y tamaño inferido con `[...]`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El arreglo `[0 0 0 0 0]`, luego `[0 0 0 0 100]`, su longitud 5,
 *              una lista de cinco números y otra de nueve números con su
 *              longitud.
 *
 * Notas:
 *    1. El tamaño forma parte del tipo: `[5]int` y `[9]int` son tipos
 *       distintos, y el tamaño no puede cambiar.
 *    2. Un índice constante fuera de rango, como `arreglo[8]`, es un error de
 *       compilación; con un índice variable, el error ocurre en ejecución
 *       (panic).
 *    3. `[...]` deja que el compilador cuente los elementos.
 */

package main

import "fmt"

func main() {
	// Declaración de un arreglo de enteros
	var arreglo [5]int
	fmt.Println("Arreglo inicial:", arreglo)

	arreglo[4] = 100 // Asignación de un valor al índice 4 del arreglo
	// arreglo[8] = 200 // ❌ Esto causaría un error de índice fuera de rango
	fmt.Println("Arreglo después de asignar valor:", arreglo)
	fmt.Println("Elemento en la posición 4:", arreglo[4])

	fmt.Println("Longitud del arreglo:", len(arreglo))

	fmt.Println("----------")

	// Declaración de un arreglo con valores iniciales
	listaNumeros := [5]int{1, 2, 3, 4, 5}
	fmt.Println("Lista de números:", listaNumeros)

	fmt.Println("----------")

	// [...] => Inferencia de tamaño del arreglo
	listaSinLimites := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println("Lista sin límites:", listaSinLimites)
	fmt.Println("Longitud de la lista sin límites:", len(listaSinLimites))
}
