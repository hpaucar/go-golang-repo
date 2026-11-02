/*
 * Archivo:     atomiccharcounter.go
 *
 * Propósito:   Contar la frecuencia de cada letra en 31 documentos RFC
 *              (rfc1000 a rfc1030) descargados en paralelo. Cada goroutine
 *              descarga un documento y suma sus letras en un slice
 *              compartido de 26 contadores con `atomic.AddInt32`.
 *
 * Compilación: go build atomiccharcounter.go
 * Ejecución:   ./atomiccharcounter
 *              (Alternativamente: go run atomiccharcounter.go)
 *
 * Entrada:     Ninguna por teclado; necesita conexión a Internet para
 *              descargar los RFC desde https://rfc-editor.org.
 * Salida:      Una línea "Completed: <url>" por documento (en cualquier
 *              orden) y, al final, la frecuencia de cada letra en el formato
 *              `a-123 b-45 ...`.
 *
 * Notas:
 *    1. Cada contador se actualiza de forma atómica, así que no hace falta
 *       un mutex aunque 31 goroutines escriban en el mismo slice.
 *    2. Si el servidor responde con un código distinto de 200, el programa
 *       termina con `panic`.
 *    3. El error de `http.Get` se ignora: sin conexión, `resp` es nil y el
 *       programa termina con un pánico de puntero nulo.
 */

package main

import (
    "fmt"
    "io"
    "net/http"
    "strings"
    "sync"
    "sync/atomic"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency []int32) {
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
            atomic.AddInt32(&frequency[cIndex], 1)
        }
    }
    fmt.Println("Completed:", url)
}


func main() {
    wg := sync.WaitGroup{}
    wg.Add(31)
    var frequency = make([]int32, 26)
    for i := 1000; i <= 1030; i++ {
        url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
        go func() {
            countLetters(url, frequency)
            wg.Done()
        }()
    }
    wg.Wait()
    for i, c := range allLetters {
        fmt.Printf("%c-%d ", c, atomic.LoadInt32(&frequency[i]))
    }
}

