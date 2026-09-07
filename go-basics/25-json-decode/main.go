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
