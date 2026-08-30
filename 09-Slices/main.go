package main

import (
	"fmt"
	"slices"
)

// Diferencia entre arreglos y slices:
//
// - Un arreglo tiene un tamaño fijo y se define con un número específico de elementos.
//
// - Un slice es una vista dinámica de un arreglo, que puede crecer o reducir su tamaño según sea necesario.
// - Los slices son más flexibles y se utilizan comúnmente en Go para manejar colecciones de datos.
// - Los slices se crean a partir de arreglos y pueden ser modificados sin necesidad de definir un tamaño fijo.
func main() {
	// 1. Uso de slices con make()

	// El "make" se usa para crear un slice con un tamaño y capacidad inicial.
	var arregloCadenas = make([]string, 3) // Inicializa el arreglo con 3 elementos
	arregloCadenas[0] = "A"
	arregloCadenas[1] = "B"
	arregloCadenas[2] = "C"
	fmt.Println("✅ arregloCadenas:", arregloCadenas)
	fmt.Println("dato 1:", arregloCadenas[1])
	fmt.Println("Longitud:", len(arregloCadenas))
	fmt.Println("Capacidad:", cap(arregloCadenas))

	arregloCadenas = append(arregloCadenas, "D", "E", "F") // Agrega más elementos al slice
	fmt.Println("\narregloCadenas después de append:", arregloCadenas)
	fmt.Println("Longitud:", len(arregloCadenas))
	fmt.Println("Capacidad:", cap(arregloCadenas))

	// 2. Uso de slices con make() y asignación de valores

	// Aquí se crea un slice con una longitud de 3 y una capacidad de 5.
	var animales = make([]string, 2, 4)                   // Crea un slice con longitud 2 y capacidad 4
	animales = []string{"Perro", "Gato", "Vaca", "Cerdo"} // Forma 1: Asigna valores al slice
	fmt.Println("\n✅ animales:", animales)
	fmt.Println("Longitud:", len(animales))
	fmt.Println("Capacidad:", cap(animales))

	animales = append(animales, "Pájaro") // Se pasa la capacidad del slice
	fmt.Println("\nanimales después de agregar Pájaro:", animales)
	fmt.Println("Longitud:", len(animales))
	fmt.Println("Capacidad:", cap(animales)) // Cambió a 8 debido al append
	// NOTA: Si se agrega un elemento más, la capacidad se duplicará nuevamente.

	// 3. Uso de slices con valores iniciales

	// Aquí se crea un slice directamente con valores iniciales.
	flores := []string{"Rosa", "Tulipán", "Lirio"}
	fmt.Println("\n✅ plantas:", flores)
	fmt.Println("Longitud:", len(flores))
	fmt.Println("Capacidad:", cap(flores))

	// 4. Comparar slices
	fmt.Println("\n✅ Comparando slices:")
	frutas := []string{"Manzana", "Banana", "Naranja"}
	if slices.Equal(flores, frutas) {
		fmt.Println("Las flores y las frutas son iguales.")
	} else {
		fmt.Println("Las flores y las frutas son diferentes.")
	}

	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	if slices.Equal(slice1, slice2) {
		fmt.Println("slice1 y slice2 son iguales.")
	} else {
		fmt.Println("slice1 y slice2 son diferentes.")
	}
}
