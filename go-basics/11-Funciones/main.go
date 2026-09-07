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
		suma += valor
	}
	return suma
}

func main() {
	resultadoVarios := sumarVarios(1, 2, 3, 4, 5)
	fmt.Println("El resultado de la suma de varios es:", resultadoVarios)

	var numero1, numero2 int
	fmt.Print("Ingrese el primer número: ")
	fmt.Scanln(&numero1) // Solicita al usuario el primer número
	fmt.Print("Ingrese el segundo número: ")
	fmt.Scanln(&numero2) // Solicita al usuario el segundo número

	resultado := sumar(numero1, numero2)
	fmt.Println("El resultado de la suma es:", resultado)
}
