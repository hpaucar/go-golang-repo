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
