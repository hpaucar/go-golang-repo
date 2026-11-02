/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar el manejo de errores en Go: errores definidos como
 *              variables (`errors.New`, `fmt.Errorf`), funciones que devuelven
 *              `(resultado, error)` y la comprobación `if err != nil`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Error: ❌ Division by zero" para 10/0, "Resultado: 5" para 10/2
 *              y "Error generado: ❌ Error genérico".
 *
 * Notas:
 *    1. En Go los errores son valores que se devuelven, no excepciones; quien
 *       llama decide qué hacer con ellos.
 *    2. Los errores definidos como variables (errores centinela, sentinel
 *       errors) se pueden comparar con `errors.Is(err, ErrDivisionByZero)`.
 *    3. `fmt.Errorf` sin verbos de formato equivale a `errors.New`; su uso
 *       típico es agregar contexto y envolver otro error con `%w`.
 */

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
