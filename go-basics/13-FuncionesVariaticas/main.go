/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar funciones variádicas: `imprimirValores` recibe cualquier
 *              cantidad de argumentos de cualquier tipo (`...any`) y los une en
 *              una sola cadena; `sumarVarios` recibe varios enteros y devuelve
 *              su suma.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Hola Pepito Perez de 30 años" y "El resultado de la suma de
 *              varios es: 131" (cada línea termina con un espacio extra).
 *
 * Notas:
 *    1. `any` es un alias de `interface{}`: acepta valores de cualquier tipo.
 *    2. Un slice existente se puede pasar a una función variádica con
 *       `f(slice...)`.
 *    3. Concatenar cadenas con `+=` dentro de un bucle crea una cadena nueva en
 *       cada paso; para muchos valores conviene `strings.Builder` o
 *       `fmt.Sprint`.
 */

package main

import "fmt"

// Usando parámetros variádicos
func imprimirValores(valores ...any) {
	var stringValores string
	for _, valor := range valores {
		stringValores += fmt.Sprintf("%v ", valor) // Convierte cada valor a string y lo concatena
	}
	fmt.Println(stringValores)
}

func sumarVarios(valores ...int) int {
	suma := 0
	for _, valor := range valores {
		suma += valor
	}
	return suma
}

func main() {
	nombre := "Pepito"
	apellido := "Perez"
	edad := 30
	imprimirValores("Hola", nombre, apellido, "de", edad, "años") // Función variática

	resultadoVarios := sumarVarios(5, 4, 22, 100)
	imprimirValores("El resultado de la suma de varios es:", resultadoVarios)
}
