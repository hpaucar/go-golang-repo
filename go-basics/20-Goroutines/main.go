/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar las goroutines. Primero llama a `saludar` de forma
 *              normal (secuencial) y luego lanza dos goroutines con `go
 *              saludar(...)`, que se ejecutan de manera concurrente mientras
 *              `main` espera un segundo con `time.Sleep`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      50 líneas "👋 Hola a todos" en orden; luego 50 líneas "🌞 Good
 *              morning" y 50 líneas "🌍 Hello world", intercaladas de forma
 *              distinta en cada ejecución; al final, "✅ Programa finalizado".
 *
 * Notas:
 *    1. La palabra `go` inicia una función en una nueva goroutine y no espera a
 *       que termine.
 *    2. Cuando `main` termina, el programa termina y se detienen todas las
 *       goroutines. `time.Sleep` solo da tiempo; no garantiza que hayan
 *       terminado. La forma correcta de esperarlas es `sync.WaitGroup` o
 *       canales.
 *    3. Las goroutines son concurrentes; se ejecutan en paralelo solo si hay
 *       varios núcleos disponibles.
 */

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
	for range 50 {
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
