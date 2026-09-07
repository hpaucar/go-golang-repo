// Ejemplo de un algoritmo concurrente en Go.
//
// El programa crea dos goroutines que acceden y modifican una misma
// variable global. Su finalidad es demostrar una condición de carrera.
package main

import (
	"fmt"
	"time"
)

// n es una variable global compartida.
//
// La goroutine principal y las goroutines que ejecutan p() y q()
// pueden acceder a esta misma posición de memoria.
var n int

// p representa el primer proceso o tarea concurrente.
func p() {
	// k1 es una variable local.
	// Solo puede ser utilizada por la goroutine que ejecuta p().
	k1 := 1

	// Escritura sobre la variable global compartida.
	// Esta instrucción compite con la escritura realizada en q().
	n = k1
}

// q representa el segundo proceso o tarea concurrente.
func q() {
	// k2 es una variable local.
	// Solo puede ser utilizada por la goroutine que ejecuta q().
	k2 := 2

	// Escritura sobre la variable global compartida.
	// No existe sincronización con la escritura realizada en p().
	n = k2
}

func main() {
	// Inicialización de la variable compartida.
	n = 0

	// Se crea una goroutine que ejecutará concurrentemente la función p().
	//
	// La instrucción "go p()" inicia la tarea, pero main no espera
	// automáticamente a que termine.
	go p()

	// Se crea otra goroutine que ejecutará concurrentemente la función q().
	//
	// El planificador de Go decide cuándo y en qué orden se ejecutan
	// las goroutines p y q.
	go q()

	// Pausa la goroutine main durante 600 milisegundos.
	//
	// Se utiliza para dar tiempo a que p y q terminen, pero no constituye
	// un mecanismo formal de sincronización. No existe una garantía lógica
	// de que ambas goroutines hayan terminado cuando transcurra este tiempo.
	time.Sleep(600 * time.Millisecond)

	// Lectura de la variable compartida.
	//
	// El resultado dependerá de cuál goroutine escribió n al final:
	//   - Si p escribe al final, n será 1.
	//   - Si q escribe al final, n será 2.
	fmt.Printf("El valor final de n es %d\n", n)
}
