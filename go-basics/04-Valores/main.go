/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar valores de distintos tipos (enteros, reales y cadenas),
 *              la inferencia de tipo y la conversión explícita de tipos. Suma
 *              un `int` con un `float64` usando `float64(...)` y concatena
 *              cadenas con el operador `+`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los números 10 y 20, el real 20.5, la suma 30.5, el nombre
 *              completo "Pepito Pérez" y la frase "Pepito Pérez tiene 32 años".
 *
 * Notas:
 *    1. Go no hace conversiones implícitas entre tipos numéricos: `numeroEntero
 *       + numeroDouble` no compila; hay que convertir con
 *       `float64(numeroEntero)`.
 *    2. Un literal como `20.5` se infiere como `float64`. Go no tiene un tipo
 *       `double`; el nombre de la variable viene de otros lenguajes.
 *    3. `fmt.Println` separa sus argumentos con un espacio.
 */

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
