// Ejemplo de un algoritmo concurrente en Go.
//
// Tres goroutines modifican una misma variable global.
// El programa permite observar una condición de carrera.
package main

import (
	"fmt"
	"time"
)

// n es una variable global compartida por la goroutine principal
// y por las goroutines que ejecutan las funciones p(), q() y r().
var n int

// p representa la primera tarea concurrente.
func p() {
	// k1 es una variable local de la goroutine que ejecuta p().
	// No es compartida con las otras goroutines.
	k1 := 1

	// Escritura en la variable global compartida.
	n = k1
}

// q representa la segunda tarea concurrente.
func q() {
	// k2 es una variable local de la goroutine que ejecuta q().
	k2 := 2

	// Escritura en la misma variable global compartida.
	n = k2
}

// r representa la tercera tarea concurrente.
func r() {
	// k3 es una variable local de la goroutine que ejecuta r().
	k3 := 3

	// Escritura en la misma variable global compartida.
	n = k3
}

func main() {
	// Inicialización de la variable compartida.
	n = 0

	// Se crean tres goroutines.
	//
	// La palabra reservada "go" inicia la ejecución concurrente
	// de cada función. El orden de estas instrucciones no garantiza
	// el orden en el que las funciones terminarán.
	go p()
	go q()
	go r()

	// Pausa únicamente la goroutine principal durante 600 milisegundos.
	//
	// Esta pausa busca dar tiempo para que p(), q() y r() terminen.
	// Sin embargo, Sleep no es un mecanismo formal de sincronización.
	time.Sleep(600 * time.Millisecond)

	// Se muestra el último valor escrito en la variable compartida n.
	//
	// El resultado podría ser:
	//   1: si p() realiza la última escritura.
	//   2: si q() realiza la última escritura.
	//   3: si r() realiza la última escritura.
	fmt.Printf("El valor final de n es %d\n", n)
}