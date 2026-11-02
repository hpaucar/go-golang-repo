/*
 * Archivo:     concurrente1.go
 *
 * Propósito:   Mostrar una condición de carrera (race condition): dos
 *              goroutines, `p` y `q`, escriben en la misma variable global `n`
 *              (los valores 1 y 2) sin sincronización, y `main` imprime el
 *              valor final después de esperar 600 ms con `time.Sleep`.
 *
 * Compilación: go build concurrente1.go
 * Ejecución:   ./concurrente1
 *              (Alternativamente: go run concurrente1.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Ya terminé desde q()" y "El valor final de n es 2". En la
 *              práctica el resultado siempre es 2: `p` espera 3600 ms antes de
 *              escribir, pero `main` termina a los 600 ms, así que `p` nunca
 *              llega a escribir ni a imprimir su mensaje.
 *
 * Notas:
 *    1. La instrucción `go` inicia una goroutine y `main` no espera a que
 *       termine; cuando `main` termina, el programa termina y las demás
 *       goroutines se detienen.
 *    2. `time.Sleep` no es un mecanismo de sincronización. Si se quita la pausa
 *       de `p` (o se alarga la de `main`), el resultado puede ser 1 o 2 según
 *       qué goroutine escriba al final.
 *    3. `go run -race concurrente1.go` detecta el acceso sin sincronización a
 *       `n`. Para coordinar goroutines se usan `sync.WaitGroup`, `sync.Mutex` o
 *       canales.
 */

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
    time.Sleep(3600 * time.Millisecond)
	// Escritura sobre la variable global compartida.
	// Esta instrucción compite con la escritura realizada en q().
	n = k1
	fmt.Println("Ya terminé desde p()")
}

// q representa el segundo proceso o tarea concurrente.
func q() {
	// k2 es una variable local.
	// Solo puede ser utilizada por la goroutine que ejecuta q().
	k2 := 2

	// Escritura sobre la variable global compartida.
	// No existe sincronización con la escritura realizada en p().
	n = k2
	fmt.Println("Ya terminé desde q()")
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
