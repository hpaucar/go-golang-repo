/*
 * Archivo:     contador.go
 *
 * Propósito:   Observar la intercalación (interleaving) de dos goroutines. Cada
 *              una incrementa diez veces la variable global `n` en dos pasos
 *              (leer y luego escribir), con pausas aleatorias muy cortas entre
 *              ellos. Si los incrementos estuvieran sincronizados, el resultado
 *              sería 20.
 *
 * Compilación: go build contador.go
 * Ejecución:   ./contador
 *              (Alternativamente: go run contador.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "El valor de n es X", con X entre 10 y 20. El valor cambia en
 *              cada ejecución; un valor menor que 20 muestra actualizaciones
 *              perdidas (lost updates).
 *
 * Notas:
 *    1. Cuando una goroutine lee `n`, se detiene y la otra escribe, la primera
 *       sobrescribe después ese incremento con un valor desactualizado.
 *    2. Las pausas de 50 a 99 nanosegundos son muy cortas; a veces no provocan
 *       un cambio de goroutine y el resultado es 20. Con pausas más largas, la
 *       pérdida se ve con más frecuencia.
 *    3. `go run -race contador.go` detecta la condición de carrera. Para un
 *       resultado correcto se protege el incremento con `sync.Mutex` o se usa
 *       `sync/atomic`, y se espera a las goroutines con `sync.WaitGroup` en
 *       lugar de `time.Sleep`.
 */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// n es una variable global compartida por main y las dos goroutines.
var n int

// pausa introduce un pequeño retardo aleatorio para favorecer
// diferentes intercalaciones entre las goroutines.
func pausa() {
	// Genera un número aleatorio entre 50 y 99.
	t := rand.Intn(50) + 50

	// Suspende la goroutine durante t nanosegundos.
	//
	// Este intervalo es extremadamente pequeño y podría no provocar
	// un cambio de goroutine de manera observable.
	time.Sleep(time.Duration(t) * time.Nanosecond)
}

// proceso realiza diez incrementos sobre la variable compartida n.
func proceso() {
	// temp es una variable local.
	// Cada goroutine tendrá su propia instancia de temp.
	var temp int

	for i := 0; i < 10; i++ {
		// Primera parte del incremento:
		// se lee la variable compartida.
		temp = n

		// Otra goroutine podría ejecutarse durante esta pausa
		// y leer o modificar n.
		pausa()

		// Segunda parte del incremento:
		// se escribe en n el valor leído anteriormente más uno.
		//
		// Como el valor de temp puede estar desactualizado,
		// esta escritura puede sobrescribir el incremento
		// realizado por la otra goroutine.
		n = temp + 1

		// Se introduce otra pausa antes de la siguiente iteración.
		pausa()
	}
}

func main() {
	// Valor inicial de la variable compartida.
	n = 0

	// Se crean dos goroutines que ejecutan la misma función.
	// Cada una intenta incrementar n diez veces.
	go proceso() // Goroutine p
	go proceso() // Goroutine q

	// Se pausa la goroutine main para dar tiempo a que las otras terminen.
	//
	// Sleep no es un mecanismo formal de sincronización.
	time.Sleep(500 * time.Millisecond)

	// Sin condición de carrera, el resultado esperado sería:
	//
	//     10 incrementos + 10 incrementos = 20
	//
	// Pero puede obtenerse un valor menor debido a actualizaciones perdidas.
	fmt.Printf("El valor de n es %d\n", n)
}