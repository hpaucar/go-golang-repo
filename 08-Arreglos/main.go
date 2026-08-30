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
