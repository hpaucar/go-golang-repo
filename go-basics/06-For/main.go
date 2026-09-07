package main

import "fmt"

func main() {
	// Usando un bucle for con una condición
	fmt.Println("Bucle for con condición:")
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// Usando un rango para iterar
	fmt.Println("Bucle for tradicional:")
	for i := range 5 {
		fmt.Println(i)
	}

	// Usando un bucle for con una condición de parada
	fmt.Println("Bucle for con condición de parada:")
	for {
		fmt.Println("Bucle infinito")
		break // Rompe el bucle infinito
	}

	// Usando un bucle for con una condición de parada y continue
	fmt.Println("Bucle for con continue:")
	for valor := range 6 {
		if valor%2 == 0 {
			continue
		}
		fmt.Println("Valor impar:", valor)
	}
}
