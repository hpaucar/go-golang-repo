package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	cadenaBase := "Hola, mundo!"

	hash := sha256.New()           // Crear un nuevo hash SHA-256
	hash.Write([]byte(cadenaBase)) // Escribir la cadena en el hash

	// Obtener el resultado del hash
	resultado := hash.Sum(nil)

	fmt.Printf("Hash SHA-256 de la cadena '%s': %x\n", cadenaBase, resultado)
}
