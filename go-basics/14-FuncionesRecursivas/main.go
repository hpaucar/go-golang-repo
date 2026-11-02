/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar la recursión con dos ejemplos: la función `factorial` y
 *              una función anónima recursiva que calcula el n-ésimo número de
 *              Fibonacci, asignada a una variable declarada antes (`var
 *              fibonacci func(numeroBase int) int`).
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El factorial de 5 (120) y de 7 (5040), y los números de
 *              Fibonacci del 1 al 7 (1, 1, 2, 3, 5, 8 y 13).
 *
 * Notas:
 *    1. Toda recursión necesita un caso base (`n <= 1`) que detenga las
 *       llamadas.
 *    2. Una función anónima solo puede llamarse a sí misma si la variable se
 *       declara antes de asignarla; con `fibonacci := func...` no compila.
 *    3. Este Fibonacci recursivo tiene un costo exponencial porque repite
 *       cálculos; para valores grandes se usa memoización o una versión
 *       iterativa.
 */

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
