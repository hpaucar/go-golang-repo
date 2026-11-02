/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar cómo codificar (serializar) una estructura de Go a JSON
 *              con `json.Marshal`, usando etiquetas de campo (struct tags) para
 *              definir los nombres de las claves.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      La estructura original `{Alice 30 alice@example.com}` y el JSON
 *              `{"name":"Alice","age":30,"email":"alice@example.com"}`.
 *
 * Notas:
 *    1. `json.Marshal` devuelve un `[]byte`; se convierte con `string(...)`
 *       para imprimirlo.
 *    2. Los campos no exportados (en minúscula) no aparecen en el JSON.
 *    3. Para un JSON legible con sangría se usa `json.MarshalIndent`; la opción
 *       `omitempty` en la etiqueta omite los campos vacíos.
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
	// Crear una instancia de Person
	persona := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	// Convertir la estructura a JSON
	fmt.Println("Estructura original:", persona)

	jsonData, err := json.Marshal(persona)
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	fmt.Println("JSON:", string(jsonData))
}
