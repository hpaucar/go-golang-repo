/*
 * Archivo:     wordfrequency.go
 *
 * Propósito:   Solución del ejercicio 3.1 del libro *Learn Concurrent
 *              Programming with Go*, que pide adaptar el programa de frecuencia
 *              de letras para contar palabras: descarga los RFC 1000 a 1020 en
 *              21 goroutines y cada una actualiza un mapa compartido
 *              `frequency` con las palabras encontradas.
 *
 * Compilación: go build wordfrequency.go
 * Ejecución:   ./wordfrequency
 *              (Alternativamente: go run wordfrequency.go)
 *
 * Entrada:     Ninguna por teclado; descarga los documentos desde
 *              https://rfc-editor.org (requiere conexión a Internet).
 * Salida:      "Completed: <url>" por cada RFC descargado y luego cada palabra
 *              con su frecuencia. Los conteos pueden ser incorrectos y el
 *              programa puede abortar con "fatal error: concurrent map writes".
 *
 * Notas:
 *    1. El programa tiene una condición de carrera intencional: varias
 *       goroutines escriben en el mismo mapa sin protección. Ejecútelo con `go
 *       run -race wordfrequency.go` para verla.
 *    2. El temporizador de 10 segundos al final puede necesitar ajustes según
 *       la velocidad de la conexión.
 *    3. No se verifica el error de `http.Get`: si falla la red, `resp` es `nil`
 *       y el programa entra en pánico. La versión corregida con mutex es el
 *       ejercicio 4.4.
 */

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

/*
Note: this program has a race condition for demonstration purposes
Additionally we have a timer at the end which you might need to adjust
depending on how fast your internet connection is.
In later chapters we cover how to wait for threads to complete their work
*/
func countLetters(url string, frequency map[string]int) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server's error: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
	for _, word := range wordRegex.FindAllString(string(body), -1) {
		wordLower := strings.ToLower(word)
		frequency[wordLower] += 1
	}
	fmt.Println("Completed:", url)
}

func main() {
	var frequency = make(map[string]int)
	for i := 1000; i <= 1020; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency)
	}
	time.Sleep(10 * time.Second)
	for k, v := range frequency {
		fmt.Println(k, "->", v)
	}
}
