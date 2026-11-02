/*
 * Archivo:     concurrente3.go
 *
 * Propósito:   Mostrar el problema de la actualización perdida (lost update).
 *              Dos goroutines ejecutan `proceso`, que lee la variable global
 *              `n`, espera 10 ms y escribe el valor leído más uno. `main`
 *              espera 500 ms e imprime `n`.
 *
 * Compilación: go build concurrente3.go
 * Ejecución:   ./concurrente3
 *              (Alternativamente: go run concurrente3.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Casi siempre "El valor final de n es 1" en lugar del 2 esperado:
 *              las dos goroutines leen `n == 0` antes de que alguna escriba.
 *
 * Notas:
 *    1. Un incremento no es atómico: se divide en leer `n`, sumar uno al valor
 *       local y escribir el resultado. La pausa de 10 ms entre la lectura y la
 *       escritura hace visible el problema.
 *    2. La lectura y la escritura de `n` forman una sección crítica (critical
 *       section); se protegen con `sync.Mutex` o se reemplazan por una
 *       operación atómica de `sync/atomic` (por ejemplo, `atomic.Int64`).
 *    3. `go run -race concurrente3.go` reporta la condición de carrera, y
 *       `time.Sleep` no es un mecanismo formal de sincronización.
 */

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