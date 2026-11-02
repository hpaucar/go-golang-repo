/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar cómo decodificar (deserializar) un texto JSON en una
 *              estructura de Go con `json.Unmarshal`, usando etiquetas de campo
 *              (struct tags) como `json:"name"` para relacionar las claves del
 *              JSON con los campos.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna (el JSON está escrito en el código).
 * Salida:      El JSON original, la estructura decodificada `{Alice 30
 *              alice@example.com}` y una línea con el nombre, la edad y el
 *              email.
 *
 * Notas:
 *    1. `json.Unmarshal` necesita un puntero (`&persona`) para poder escribir
 *       en la variable.
 *    2. Solo se decodifican los campos exportados (que empiezan con mayúscula);
 *       las claves del JSON que no coinciden se ignoran.
 *    3. Si el JSON está mal formado o un tipo no coincide, `Unmarshal` devuelve
 *       un error.
 */

package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func main() {
	// JSON de entrada
	jsonData := `{"name":"Alice","age":30,"email":"alice@example.com"}`

	// Variable para almacenar el resultado
	var persona Person

	// Decodificar el JSON
	fmt.Println("Estructura original en JSON:", jsonData)

	err := json.Unmarshal([]byte(jsonData), &persona)
	if err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return
	}

	// Mostrar el resultado
	fmt.Println("Estructura decodificada:", persona)
	fmt.Printf("Nombre: %s, Edad: %d, Email: %s\n", persona.Name, persona.Age, persona.Email)
}
