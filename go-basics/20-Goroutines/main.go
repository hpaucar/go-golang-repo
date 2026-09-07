package main

import (
	"fmt"
	"time"
)

//
// NOTE: Gorutines
//
// Las goroutines permiten ejecutar código concurrentemente,
// facilitando la creación de programas eficientes y rápidos en Go.
//
// Similar al concepto de hilos (threads) en otros lenguajes, o al
// async/await en JavaScript, las goroutines son ligeras y se ejecutan
// en paralelo, permitiendo que múltiples tareas se realicen al mismo tiempo.
//

func saludar(saludo string) {
	for range 5 {
		fmt.Println(saludo)
	}
}

func main() {
	saludar("👋 Hola a todos")

	go saludar("🌞 Good morning")
	go saludar("🌍 Hello world")

	// Espera para que la goroutine termine antes de que finalice el programa
	time.Sleep(1 * time.Second)

	fmt.Println("✅ Programa finalizado")
}
