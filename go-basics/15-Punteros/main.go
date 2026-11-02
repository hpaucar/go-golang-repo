/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar los punteros: obtener la dirección de una variable con
 *              `&`, leer y modificar el valor apuntado con `*`, y pasar un
 *              puntero a un arreglo a una función (`modificarArreglo`) para
 *              cambiar el arreglo original.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El valor de `a` (10 y luego 20), su dirección de memoria (cambia
 *              en cada ejecución), el valor y la dirección guardados en `p`, y
 *              el arreglo antes (`[1 2 3 4 5]`) y después (`[55 2 3 4 333]`) de
 *              modificarlo.
 *
 * Notas:
 *    1. En Go los argumentos se pasan por valor: un arreglo se copia completo,
 *       así que para modificar el original se pasa un puntero.
 *    2. `(*arg)[0]` y `arg[0]` son equivalentes: Go desreferencia
 *       automáticamente un puntero a arreglo al indexar.
 *    3. Go no tiene aritmética de punteros. El comentario de `modificarArreglo`
 *       habla de un "slice", pero el parámetro es un puntero a un arreglo
 *       `[5]int`.
 */

package main

import "fmt"

// modificarArreglo recibe un puntero a un slice de enteros y modifica su primer elemento
func modificarArreglo(arg *[5]int) {
	(*arg)[0] = 55  // Cambia el primer elemento a 55
	(*arg)[4] = 333 // Cambia el ultimo elemento a 333
}

func main() {
	// 1. Ejercicio básico
	a := 10
	p := &a // p es un puntero a la variable a

	fmt.Println("Valor de a:", a)      // Imprime el valor de a
	fmt.Println("Dirección de a:", &a) // Imprime la dirección de memoria de a

	*p = 20 // Cambia el valor de a a través del puntero p

	fmt.Println("Valor de a:", a) // Imprime el valor de a

	fmt.Println("Valor de p:", *p) // Imprime el valor de p (valor apuntado por p)
	fmt.Println("Remoto de p:", p) // Imprime la dirección de memoria de a

	// 2. Ejercicio con punteros y arreglos
	numeros := [5]int{1, 2, 3, 4, 5}
	fmt.Println("\nNúmeros originales:", numeros)
	modificarArreglo(&numeros)                   // Pasa el puntero al arreglo
	fmt.Println("Números modificados:", numeros) // Imprime el arreglo modificado
}
