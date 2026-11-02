/*
 * Archivo:     charcountermessagepass.go
 *
 * Propósito:   Ilustrar el paso de mensajes en lugar de memoria compartida:
 *              31 goroutines descargan los RFC 1000 a 1030 y cada una
 *              devuelve por su propio canal un arreglo con la frecuencia de
 *              cada letra; `main` suma los resultados.
 *
 * Compilación: go build charcountermessagepass.go
 * Ejecución:   ./charcountermessagepass
 *              (Alternativamente: go run charcountermessagepass.go)
 *
 * Entrada:     Red: descarga los RFC 1000 a 1030 desde
 *              https://rfc-editor.org/rfc/rfc<N>.txt.
 * Salida:      Una línea "Completed: <url>" por documento, en orden variable,
 *              y al final la frecuencia total de cada letra, "a-N b-N ..."
 *              (sin salto de línea final).
 *
 * Notas:
 *    1. Cada goroutine trabaja con su propio arreglo `frequency`; no hay datos
 *       compartidos, así que no se necesitan mutex.
 *    2. `main` lee los canales en orden, pero las descargas corren en
 *       paralelo; el orden de los "Completed:" es no determinista.
 *    3. El error de `http.Get` se ignora: sin red, `resp` es `nil` y el
 *       programa termina con un pánico (nil pointer dereference). Si el
 *       servidor no devuelve 200, también hay pánico.
 */

package main

import (
    "fmt"
    "io"
    "net/http"
    "strings"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string) <-chan []int {
    result := make(chan []int)
    go func() {
        defer close(result)
        frequency := make([]int, 26)
        resp, _ := http.Get(url)
        defer resp.Body.Close()
        if resp.StatusCode != 200 {
            panic("Server returning error code: " + resp.Status)
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
        result <- frequency
    }()
    return result
}

func main() {
    results := make([]<-chan []int, 0)
    totalFrequencies := make([]int, 26)
    for i := 1000; i <= 1030; i++ {
        url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
        results = append(results, countLetters(url))
    }
    for _, c := range results {
        frequencyResult := <-c
        for i := 0; i < 26; i++ {
            totalFrequencies[i] += frequencyResult[i]
        }
    }
    for i, c := range allLetters {
        fmt.Printf("%c-%d ", c, totalFrequencies[i])
    }
}
