/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar la estructura mínima de un programa ejecutable en Go: la
 *              declaración `package main`, la importación del paquete `fmt` y
 *              la función `main`, que imprime un saludo en la consola.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El texto "Hola Mundo" seguido de un salto de línea.
 *
 * Notas:
 *    1. El paquete `main` indica que el archivo pertenece a un programa
 *       ejecutable y no a una biblioteca.
 *    2. La función `main` es el punto de inicio del programa; no recibe
 *       argumentos ni devuelve valores.
 *    3. `fmt.Println` escribe en la salida estándar y agrega un salto de línea
 *       al final.
 */

package main

// Importa el paquete fmt, utilizado para mostrar información en la consola.
import "fmt"

// main es la función principal y el punto de inicio del programa.
func main() {
	// Imprime el mensaje "Hola Mundo" en la consola y añade un salto de línea.
	fmt.Println("Hola Mundo")
}