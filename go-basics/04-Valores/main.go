package main

import "fmt"

func main() {
	var numero1 int = 10
	numero2 := 20 // Inferencia de tipo
	fmt.Println(numero1, numero2)

	var numeroEntero = 10
	numeroDouble := 20.5 // Inferencia de tipo
	fmt.Println(numeroDouble)

	// Sumar un entero y un double con conversión de tipos
	resultado := float64(numeroEntero) + numeroDouble
	fmt.Println(resultado)

	var nombre string = "Pepito"
	apellido := "Pérez"
	edad := 32
	nombreCompleto := nombre + " " + apellido
	fmt.Println(nombreCompleto)
	fmt.Println(nombreCompleto, "tiene", edad, "años")
}
