/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar cómo declarar funciones con parámetros y valor de
 *              retorno (`sumar`) y con parámetros variádicos (`sumarVarios`).
 *              El programa también lee dos números enteros desde el teclado con
 *              `fmt.Scanln` y muestra su suma.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Dos números enteros escritos por el usuario en la entrada
 *              estándar, uno por línea.
 * Salida:      La suma 150 de los valores fijos, los mensajes que piden cada
 *              número y la suma de los dos números ingresados.
 *
 * Notas:
 *    1. El tipo se escribe después del nombre del parámetro, y el tipo de
 *       retorno después de la lista de parámetros.
 *    2. `valores ...int` recibe cualquier cantidad de enteros; dentro de la
 *       función se usa como un slice `[]int`.
 *    3. El programa ignora el error que devuelve `fmt.Scanln`: si se escribe un
 *       texto no numérico, la variable queda en 0 sin aviso.
 */

package main

import "fmt"

// Usando parámetros normales
func sumar(valor1 int, valor2 int) int {
	return valor1 + valor2
}

// Usando parámetros variádicos
func sumarVarios(valores ...int) int {
	suma := 0
	for _, valor := range valores {
		suma += valor // suma = suma + valor
	}
	return suma
}

func main() {
	resultadoVarios := sumarVarios(10, 20, 30, 40, 50)
	fmt.Println("El resultado de la suma de varios es:", resultadoVarios)

	var numero1, numero2 int
	fmt.Print("Ingrese el primer número: ")
	fmt.Scanln(&numero1) // Solicita al usuario el primer número
	fmt.Print("Ingrese el segundo número: ")
	fmt.Scanln(&numero2) // Solicita al usuario el segundo número

	resultado := sumar(numero1, numero2)
	fmt.Println("El resultado de la suma es:", resultado)
}
