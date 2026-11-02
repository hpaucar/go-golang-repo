/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar cómo importar varios paquetes de la biblioteca estándar
 *              en un bloque `import`. Usa `os.Getenv` para leer la variable de
 *              entorno `HOME` y `os.Create` para crear un archivo de texto
 *              vacío llamado `ejemplo.txt`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     La variable de entorno `HOME`.
 * Salida:      La ruta del directorio personal, o un mensaje si `HOME` no
 *              existe. Además crea (o vacía, si ya existe) el archivo
 *              `ejemplo.txt` en el directorio actual; si no puede crearlo,
 *              muestra un mensaje de error.
 *
 * Notas:
 *    1. Un bloque `import ( ... )` agrupa varios paquetes; cada paquete se usa
 *       con su nombre como prefijo (`os.Getenv`, `fmt.Println`).
 *    2. `os.Create` devuelve el archivo y un error; en Go los errores se
 *       revisan de forma explícita con `if err != nil`.
 *    3. `defer file.Close()` cierra el archivo al salir de `main`, aunque la
 *       función termine antes con un `return`.
 */

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
