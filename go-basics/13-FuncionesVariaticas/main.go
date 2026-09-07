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
