/*
 * Archivo:     charcounterconcurrent.go
 *
 * Propósito:   Descargar concurrentemente los RFC 1000 a 1030, una goroutine
 *              por documento, y acumular las frecuencias de letras en un slice
 *              compartido. Ilustra una carrera de datos al actualizar los
 *              contadores sin protección y una espera fija de 10 segundos que
 *              no garantiza que todas las tareas terminen.
 *
 * Compilación: go build charcounterconcurrent.go
 * Ejecución:   ./charcounterconcurrent
 *              (Alternativamente: go run charcounterconcurrent.go)
 *
 * Entrada:     Ninguna por teclado. Requiere conexión a Internet (31 peticiones
 *              HTTP).
 * Salida:      Líneas "Completed: <url>" en orden no determinista y, tras 10
 *              segundos, las frecuencias "a-N ... z-N". Los totales pueden ser
 *              menores que los de la versión secuencial por las actualizaciones
 *              perdidas.
 *
 * Notas:
 *    1. `frequency[cIndex] += 1` no es atómico: varias goroutines leen y
 *       escriben el mismo elemento a la vez. Use `go run -race
 *       charcounterconcurrent.go` para detectarlo.
 *    2. Si la conexión es lenta, algunas descargas no terminan antes de los 10
 *       segundos y sus letras no se cuentan; ajuste el tiempo según su red.
 *    3. El capítulo 4 corrige la carrera con un mutex y el capítulo 6 reemplaza
 *       la espera fija por `sync.WaitGroup`.
 */

package main

import (
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

/*
Note: this program has a race condition for demonstration purposes
Additionally we have a timer at the end which you might need to adjust
depending on how fast your internet connection is.
In later chapters we cover how to wait for threads to complete their work
*/
func countLetters(url string, frequency []int) {
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        panic("Server returning error status code: " + resp.Status)
    }
    body, _ := io.ReadAll(resp.Body)
    for _, b := range body {
        c := strings.ToLower(string(b))
        cIndex := strings.Index(allLetters, c)
        if cIndex >= 0 {
            frequency[cIndex] += 1
        }
    }
    fmt.Println("Completed:", url)
}

func main() {
    var frequency = make([]int, 26)
    for i := 1000; i <= 1030; i++ {
        url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
        go countLetters(url, frequency)
    }
    time.Sleep(10 * time.Second)
    for i, c := range allLetters {
        fmt.Printf("%c-%d ", c, frequency[i])
    }
}
