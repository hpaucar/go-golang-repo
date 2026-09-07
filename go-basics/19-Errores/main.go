package main

import (
	"errors"
	"fmt"
)

// Crear los errores personalizados
var ErrDivisionByZero = errors.New("❌ Division by zero")
var ErrGeneric = fmt.Errorf("❌ Error genérico")

// dividir realiza la división de dos números y devuelve un error si el divisor es cero.
func dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero // Retorna un error si el divisor es cero
	}
	return a / b, nil // Retorna el resultado de la división y nil como error
}

func generarError() error {
	return ErrGeneric
}

func main() {
	// Ejemplo de uso de la función dividir
	resultado, err := dividir(10, 0) // Intentamos dividir por cero, lo que generará un error
	// Manejo del error
	// Si err no es nil, significa que ocurrió un error
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Resultado:", resultado)
	}

	// Otro ejemplo de uso de la función dividir
	resultado, err = dividir(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Resultado:", resultado)
	}

	// Generar un error genérico
	err = generarError()
	fmt.Println("Error generado:", err)
}
