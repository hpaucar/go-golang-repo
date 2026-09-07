package main

import "fmt"

// Función recursiva que calcula el factorial de un número entero n.
// Ejemplo: factorial(5) = 5 * 4 * 3 * 2 * 1 = 120
// La función se llama a sí misma con n-1 hasta que n es 1 o menor.
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1) // Llamada recursiva
}

func main() {
	fmt.Println("Factorial de 5:", factorial(5))
	fmt.Println("Factorial de 7:", factorial(7))

	var fibonacci func(numeroBase int) int // Declaración de la función recursiva anónima para Fibonacci.

	// Función recursiva anónima para calcular el n-ésimo número de Fibonacci.
	fibonacci = func(numeroBase int) int {
		if numeroBase <= 1 {
			return numeroBase
		}
		return fibonacci(numeroBase-1) + fibonacci(numeroBase-2)
	}

	fmt.Println("Fibonacci de 1:", fibonacci(1))
	fmt.Println("Fibonacci de 2:", fibonacci(2))
	fmt.Println("Fibonacci de 3:", fibonacci(3))
	fmt.Println("Fibonacci de 4:", fibonacci(4))
	fmt.Println("Fibonacci de 5:", fibonacci(5))
	fmt.Println("Fibonacci de 6:", fibonacci(6))
	fmt.Println("Fibonacci de 7:", fibonacci(7))
}
