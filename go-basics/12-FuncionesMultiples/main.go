package main

import "fmt"

// función que devuelve múltiples valores: suma, resta, multiplicación y división
func calcular(a, b float64) (float64, float64, float64, float64) {
	suma := a + b
	resta := a - b
	multiplicacion := a * b
	var division float64
	if b != 0 {
		division = a / b
	}

	return suma, resta, multiplicacion, division
}

// función que usa un return nombrado para calcular la exponenciación
func exponenciar(a, b float64) (resultado float64) {
	resultado = 1
	for range int(b) {
		resultado *= a
	}
	return
}

func main() {
	// Valores de ejemplo
	x, y := 10.0, 5.0

	// Llamando a la función con múltiples valores de retorno
	suma, resta, mult, div := calcular(x, y)

	fmt.Printf("Números: %.1f y %.1f\n", x, y)
	fmt.Printf("Suma: %.1f\n", suma)   // Corregí este error que tenía en el original
	fmt.Printf("Resta: %.1f\n", resta) // Corregí el nombre de la variable
	fmt.Printf("Multiplicación: %.1f\n", mult)
	fmt.Printf("División: %.1f\n", div)

	// Usando identificador en blanco para ignorar algunos valores de retorno
	soloSuma, _, _, _ := calcular(x, y)
	fmt.Printf("Solo suma: %.1f\n", soloSuma)

	// Ignorando la suma y la división
	_, _, soloMultiplicacion, _ := calcular(x, y)
	fmt.Printf("Solo multiplicación: %.1f\n", soloMultiplicacion)

	// Ignorando la multiplicación y la suma
	_, soloResta, _, soloDivision := calcular(x, y)
	fmt.Printf("Solo resta: %.1f\n", soloResta)
	fmt.Printf("Solo división: %.1f\n", soloDivision)

	// Función con return nombrado
	resultado := exponenciar(2, 3)
	fmt.Printf("2 elevado a 3 es: %.1f\n", resultado)
}
