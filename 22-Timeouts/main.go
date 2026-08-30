package main

import (
	"fmt"
	"time"
)

// !timeout
// Este ejercicio muestra cómo manejar un timeout al esperar una respuesta de una API.
// Utiliza un canal para recibir la respuesta y un select para manejar el timeout.

func main() {
	// ✅ Ejercicio 1
	fmt.Println("Ejercicio 1 - Iniciando consulta a la API...")
	// Creamos un canal para recibir la respuesta1 de la API
	respuesta1 := make(chan string)

	go func() {
		// Simula una llamada HTTP que puede tardar
		time.Sleep(3 * time.Second) // Simula un retraso de 3 segundos
		// Aquí se simula la respuesta de la API
		respuesta1 <- "Datos de la API"
	}()

	// Usamos select para esperar la respuesta o un timeout
	select {
	// Si la respuesta llega antes del timeout, se procesa
	case datos := <-respuesta1:
		fmt.Println("✅ Datos recibidos:", datos)

	// Si no llega respuesta en 2 segundos, se maneja el timeout
	case <-time.After(2 * time.Second):
		fmt.Println("⏰ Timeout: La API no respondió")
	}

	// ✅ Ejercicio 2
	fmt.Println("\nEjercicio 2 - Iniciando segunda consulta a la API...")
	respuesta2 := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		respuesta2 <- "Datos de la API"
	}()
	select {
	case datos := <-respuesta2:
		fmt.Println("✅ Datos recibidos:", datos)
	case <-time.After(2 * time.Second):
		fmt.Println("⏰ Timeout: La API no respondió")
	}
}
