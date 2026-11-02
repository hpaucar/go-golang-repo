/*
 * Archivo:     main.go
 *
 * Propósito:   Simular la preparación de un pedido con tres goroutines
 *              (hamburguesa, papas y soda), cada una con un tiempo distinto (3,
 *              2 y 1 s). `main` usa `select` dentro de un bucle de tres
 *              iteraciones para recibir cada ítem en el orden en que queda
 *              listo.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "🥤 Pedido listo: Soda fría", "🍟 Pedido listo: Papas fritas" y "🍔
 *              Pedido listo: Hamburguesa con queso", en ese orden, y luego "✅
 *              ¡El cliente ya tiene todo su pedido!". Tarda unos 3 segundos.
 *
 * Notas:
 *    1. Las tres preparaciones ocurren al mismo tiempo, así que el tiempo total
 *       es el de la más lenta (3 s) y no la suma (6 s).
 *    2. El bucle se repite tantas veces como mensajes se esperan, y cada
 *       `select` recibe exactamente un mensaje.
 *    3. Contar mensajes es una forma simple de sincronización; cuando el número
 *       de tareas varía, se usa `sync.WaitGroup` y se cierra un canal de
 *       resultados.
 */

package main

import (
	"fmt"
	"time"
)

// Este ejercicio implementa:
// !channels:
// Para manejar la comunicación entre goroutines
// !goroutines
// Para ejecutar tareas simultáneamente
// !select:
// Para esperar múltiples operaciones

func main() {
	// Creamos un canal por cada ítem del pedido
	hamburguesa := make(chan string)
	papas := make(chan string)
	soda := make(chan string)

	// Simulamos la preparación de cada cosa en una goroutine distinta
	go prepararHamburguesa(hamburguesa)
	go prepararPapas(papas)
	go prepararSoda(soda)

	// Esperamos a que cada parte del pedido esté lista
	const cantidadGoroutines = 3

	// Usando el for range para recibir mensajes de los canales
	// hasta que hayamos recibido todos los pedidos
	for range cantidadGoroutines {
		select {
		case item := <-hamburguesa:
			fmt.Println("🍔 Pedido listo:", item)
		case item := <-papas:
			fmt.Println("🍟 Pedido listo:", item)
		case item := <-soda:
			fmt.Println("🥤 Pedido listo:", item)
		}
	}

	fmt.Println("✅ ¡El cliente ya tiene todo su pedido!")
}

// Simulan funciones de preparación con distintos tiempos

func prepararHamburguesa(ch chan string) {
	time.Sleep(3 * time.Second)
	ch <- "Hamburguesa con queso"
}

func prepararPapas(ch chan string) {
	time.Sleep(2 * time.Second)
	ch <- "Papas fritas"
}

func prepararSoda(ch chan string) {
	time.Sleep(1 * time.Second)
	ch <- "Soda fría"
}
