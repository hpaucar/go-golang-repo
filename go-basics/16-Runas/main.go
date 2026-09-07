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
