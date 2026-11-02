/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar las formas de declarar variables en Go: `var` con tipo
 *              explícito, declaración múltiple en una línea, variables sin
 *              inicializar y declaración corta con inferencia de tipo (`:=`).
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los valores de variables de tipo `string`, `int` y `bool`,
 *              incluido el valor 0 de una variable `int` sin inicializar y su
 *              valor 100 después de asignarlo.
 *
 * Notas:
 *    1. Toda variable declarada sin valor recibe su valor cero (zero value): 0
 *       para números, "" para cadenas y false para booleanos.
 *    2. El operador `:=` declara e inicializa a la vez, y el compilador infiere
 *       el tipo; solo puede usarse dentro de funciones.
 *    3. Go no compila si una variable local se declara y no se usa.
 */

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
