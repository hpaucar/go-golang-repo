// Ejemplo de concurrencia con una variable global compartida.
//
// Dos goroutines ejecutan la misma función e intentan incrementar
// simultáneamente la variable global n.
package main

import (
	"fmt"
	"time"
)

// n es una variable global compartida por todas las goroutines.
var n int

// proceso lee, modifica y escribe la variable compartida n.
func proceso() {
	// temp es una variable local.
	// Cada goroutine posee su propia variable temp.
	var temp int

	// Se lee el valor de la variable compartida n.
	temp = n

	// Se introduce un retardo de 10 milisegundos.
	// Durante este tiempo, otra goroutine puede leer o modificar n.
	// El retardo aumenta la probabilidad de observar la condición de carrera.
	time.Sleep(10 * time.Millisecond)

	// Se incrementa el valor previamente leído y se escribe en n.
	//
	// El incremento no es una operación atómica, pues se divide en:
	//   1. Lectura de n.
	//   2. Incremento del valor local.
	//   3. Escritura del resultado en n.
	n = temp + 1
}

func main() {
	// Se inicializa la variable global compartida.
	n = 0

	// Se crea la primera goroutine, denominada conceptualmente p.
	go proceso() // p

	// Se crea la segunda goroutine, denominada conceptualmente q.
	go proceso() // q

	// Se pausa la goroutine principal para dar tiempo a que p y q terminen.
	//
	// Sleep no es un mecanismo formal de sincronización porque no comprueba
	// realmente que las dos goroutines hayan finalizado.
	time.Sleep(500 * time.Millisecond)

	// Se muestra el valor final de la variable compartida.
	fmt.Printf("El valor final de n es %d\n", n)
}