package main

import "fmt"

// modificarArreglo recibe un puntero a un slice de enteros y modifica su primer elemento
func modificarArreglo(arg *[5]int) {
	(*arg)[0] = 55  // Cambia el primer elemento a 55
	(*arg)[4] = 333 // Cambia el segundo elemento a 333
}

func main() {
	// 1. Ejercicio básico
	a := 10
	p := &a // p es un puntero a la variable a

	fmt.Println("Valor de a:", a)      // Imprime el valor de a
	fmt.Println("Dirección de a:", &a) // Imprime la dirección de memoria de a

	*p = 20 // Cambia el valor de a a través del puntero p

	fmt.Println("Valor de a:", a) // Imprime el valor de a

	fmt.Println("Valor de p:", *p) // Imprime el valor de p (valor apuntado por p)
	fmt.Println("Remoto de p:", p) // Imprime la dirección de memoria de a

	// 2. Ejercicio con punteros y arreglos
	numeros := [5]int{1, 2, 3, 4, 5}
	fmt.Println("\nNúmeros originales:", numeros)
	modificarArreglo(&numeros)                   // Pasa el puntero al arreglo
	fmt.Println("Números modificados:", numeros) // Imprime el arreglo modificado
}
