/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar la diferencia entre bytes y runas (runes) en las cadenas
 *              UTF-8. Compara la longitud en bytes (`len`) con el número de
 *              caracteres (`utf8.RuneCountInString`) y recorre una cadena en
 *              chino por bytes y por runas con `range`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      La longitud 5 de "Hello", 6 bytes y 2 runas para "你好", seis
 *              caracteres sin sentido al recorrer "你好" byte por byte, cada runa
 *              con su código Unicode (U+4F60 y U+597D) y los códigos 20320 y
 *              241 de '你' y 'ñ'.
 *
 * Notas:
 *    1. Una cadena en Go es una secuencia de bytes; `s[i]` devuelve un byte, no
 *       un carácter. Por eso el recorrido por índice muestra caracteres
 *       extraños.
 *    2. `for idx, valor := range s` decodifica UTF-8: `valor` es una runa e
 *       `idx` es la posición en bytes donde empieza (0 y 3).
 *    3. `rune` es un alias de `int32` y guarda un punto de código Unicode. Las
 *       etiquetas impresas dicen "saludoEspanol", pero la constante se llama
 *       `saludoIngles`.
 */

package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	const saludoIngles = "Hello"
	fmt.Println("\nsaludoEspanol:", saludoIngles)
	fmt.Println("Longitud de saludoEspanol:", len(saludoIngles))

	const saludoChino = "你好"
	fmt.Println("saludoChino:", saludoChino)
	fmt.Println("Longitud de saludoChino:", len(saludoChino))

	// Iterar sobre los caracteres de saludoChino
	fmt.Println("\nIterando sobre los caracteres de saludoChino:")
	for i := range len(saludoChino) {
		fmt.Printf("saludoChino[%d]: %c\n", i, saludoChino[i])
	}

	// Esto cuenta los caracteres correctamente, incluyendo los multibyte
	fmt.Println("\nLongitud de saludoChino:", utf8.RuneCountInString(saludoChino))

	for idx, valor := range saludoChino {
		fmt.Printf("%#U en saludoChino[%d]: %c\n", valor, idx, valor)
	}

	var palabraChino rune = '你'
	fmt.Println("\nCódigo de palabraChino:", palabraChino)

	var r rune = 'ñ'
	fmt.Println("\nCódigo de r:", r)
}
