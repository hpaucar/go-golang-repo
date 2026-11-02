/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar la sentencia condicional `if`/`else`/`else if`, incluida
 *              la forma con una declaración corta (`if numero := 10; numero <
 *              0`), cuyo alcance se limita al bloque `if`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Hola Pepito", "Eres mayor de edad", "El número 8 es par" y "El
 *              número es positivo".
 *
 * Notas:
 *    1. En Go la condición del `if` no lleva paréntesis, pero las llaves son
 *       obligatorias.
 *    2. La variable declarada en `if numero := 10; ...` solo existe dentro del
 *       `if` y de sus ramas `else`.
 *    3. El programa usa la función predefinida `println` y no `fmt.Println`:
 *       `println` escribe en la salida de error (stderr) y se usa sobre todo
 *       para depuración. En programas reales se prefiere `fmt.Println`.
 */

package main

func main() {
	nombre := "Pepito"
	edad := 32

	if nombre == "Pepito" {
		println("Hola Pepito")
	} else {
		println("Hola desconocido")
	}

	if edad >= 18 {
		println("Eres mayor de edad")
	} else {
		println("Eres menor de edad")
	}

	if 8%2 == 0 {
		println("El número 8 es par")
	} else {
		println("El número 8 es impar")
	}

	if numero := 10; numero < 0 {
		println("El número es negativo")
	} else if numero == 0 {
		println("El número es cero")
	} else {
		println("El número es positivo")
	}
}
