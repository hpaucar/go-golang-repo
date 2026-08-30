package main

import (
	"fmt"
	"os"
)

func main() {
	envVar := os.Getenv("HOME")
	// La variable de entorno HOME contiene la ruta del directorio personal del usuario
	if envVar == "" {
		fmt.Println("No se encontró la variable de entorno HOME")
	} else {
		fmt.Printf("La variable de entorno HOME es: %s\n", envVar)
	}

	// Crear un archivo de texto
	file, err := os.Create("ejemplo.txt")
	// Si el archivo no se puede crear, se imprime un mensaje de error
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return
	}
	defer file.Close()
}
