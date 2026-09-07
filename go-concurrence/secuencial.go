// Ejemplo de un algoritmo secuencial.
//
// Todas las instrucciones se ejecutan una después de otra,
// siguiendo el orden en el que aparecen en el programa.
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