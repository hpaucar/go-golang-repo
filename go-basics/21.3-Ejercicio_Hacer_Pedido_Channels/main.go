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
