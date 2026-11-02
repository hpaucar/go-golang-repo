/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar cómo aplicar un tiempo límite (timeout) al esperar una
 *              respuesta, con `select` y `time.After`. Hay dos ejercicios: una
 *              "API" simulada que tarda 3 s (supera el límite de 2 s) y otra
 *              que tarda 1 s (responde a tiempo).
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      En el ejercicio 1, "⏰ Timeout: La API no respondió" después de 2
 *              s; en el ejercicio 2, "✅ Datos recibidos: Datos de la API"
 *              después de 1 s más.
 *
 * Notas:
 *    1. `select` toma el primer caso que esté listo: la respuesta o el canal de
 *       `time.After`.
 *    2. En el ejercicio 1, la goroutine que responde tarde queda bloqueada para
 *       siempre al enviar en `respuesta1`, porque nadie recibe (goroutine
 *       leak). Un canal con búfer de tamaño 1 evitaría ese bloqueo.
 *    3. En programas reales se usa `context.WithTimeout`, que además permite
 *       cancelar la operación en curso.
 */

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
