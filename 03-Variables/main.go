package main

import "fmt"

func main() {
	var cadena string = "Inicial"
	fmt.Println("Cadena:", cadena)

	var entero int = 10
	fmt.Println("Entero:", entero)

	// Declarar múltiples variables
	var entero1, entero2 int = 20, 30
	fmt.Println("Entero1:", entero1)
	fmt.Println("Entero2:", entero2)

	var estaActivo bool = true
	fmt.Println("Esta activo:", estaActivo)

	// Declara una variable sin inicializar
	var enteroSimple int
	fmt.Println("Entero simple:", enteroSimple) // Imprime: 0 (Valor por defecto de int)
	enteroSimple = 100
	fmt.Println("Entero simple:", enteroSimple) // Imprime: 100

	// Declarar una variable con inferencia de tipo
	fruta := "Manzana"
	fmt.Println("Fruta:", fruta)
}
