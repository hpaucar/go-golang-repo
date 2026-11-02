/*
 * Archivo:     wordfrequencymutex.go
 *
 * Propósito:   Solución del ejercicio 4.4 del libro *Learn Concurrent
 *              Programming with Go*, que pide eliminar la condición de carrera
 *              del contador de palabras del ejercicio 3.1 usando un
 *              `sync.Mutex`. Cada goroutine descarga un RFC (1000 a 1020) y
 *              actualiza el mapa `frequency` dentro de una sección crítica.
 *
 * Compilación: go build wordfrequencymutex.go
 * Ejecución:   ./wordfrequencymutex
 *              (Alternativamente: go run wordfrequencymutex.go)
 *
 * Entrada:     Ninguna por teclado; descarga los documentos desde
 *              https://rfc-editor.org (requiere conexión a Internet).
 * Salida:      "Completed: <url>" por cada RFC y luego cada palabra con su
 *              frecuencia, ahora sin carreras de datos.
 *
 * Notas:
 *    1. La descarga y el análisis con la expresión regular ocurren fuera del
 *       mutex; solo la actualización del mapa está en la sección crítica.
 *    2. `main` sigue esperando con `time.Sleep(10 * time.Second)`; si la red es
 *       lenta, algunos RFC no se cuentan.
 *    3. No se verifica el error de `http.Get`: si falla la red, `resp` es `nil`
 *       y el programa entra en pánico.
 */

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

/*
Note: In this program we have a timer at the end which you might need to adjust
depending on how fast your internet connection is.
In later chapters we cover how to wait for threads to complete their work
*/
func countLetters(url string, frequency map[string]int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server's error: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
	mutex.Lock()
	for _, word := range wordRegex.FindAllString(string(body), -1) {
		wordLower := strings.ToLower(word)
		frequency[wordLower] += 1
	}
	mutex.Unlock()
	fmt.Println("Completed:", url)
}

func main() {
	mutex := sync.Mutex{}
	var frequency = make(map[string]int)
	for i := 1000; i <= 1020; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency, &mutex)
	}
	time.Sleep(10 * time.Second)
	mutex.Lock()
	for k, v := range frequency {
		fmt.Println(k, "->", v)
	}
	mutex.Unlock()
}
