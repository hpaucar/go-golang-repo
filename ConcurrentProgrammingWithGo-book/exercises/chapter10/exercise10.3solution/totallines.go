/*
 * Archivo:     totallines.go
 *
 * Propósito:   Solución del ejercicio 10.3 del libro *Learn Concurrent
 *              Programming with Go*, que pide contar las líneas de 30 RFC con
 *              el patrón fan-out/fan-in: 30 goroutines descargan en paralelo y
 *              envían sus conteos al canal `linesOnPage`; otra goroutine los
 *              suma y entrega el total por `finalResult`.
 *
 * Compilación: go build totallines.go
 * Ejecución:   ./totallines
 *              (Alternativamente: go run totallines.go)
 *
 * Entrada:     Ninguna por teclado; descarga los documentos desde
 *              https://rfc-editor.org (requiere conexión a Internet).
 * Salida:      "Downloading <url>" por cada RFC (en orden variable) y al final
 *              "Total lines: <n>", el mismo total que la versión secuencial.
 *
 * Notas:
 *    1. Fan-out: muchas goroutines producen; fan-in: una sola goroutine consume
 *       y combina.
 *    2. La goroutine sumadora sabe cuántos valores leer (`pagesToDownload`),
 *       así que no necesita que se cierre el canal.
 *    3. No se verifica el error de `http.Get`: si falla la red, el programa
 *       entra en pánico.
 */

package main

import (
    "fmt"
    "io"
    "net/http"
    "strings"
)

func main() {
    const pagesToDownload = 30
    linesOnPage := make(chan int)
    finalResult := make(chan int)
    for i := 1000; i < 1000 + pagesToDownload; i++ {
        go func(id int) {
            url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", id)
            fmt.Println("Downloading", url)
            resp, _ := http.Get(url)
            if resp.StatusCode != 200 {
                panic("Server's error: " + resp.Status)
            }
            bodyBytes, _ := io.ReadAll(resp.Body)
            linesOnPage <- strings.Count(string(bodyBytes), "\n")
            resp.Body.Close()
        }(i)
    }
    go func() {
        totalLines := 0
        for i := 0; i < pagesToDownload; i++ {
            totalLines += <-linesOnPage
        }
        finalResult <- totalLines
    }()
    fmt.Println("Total lines:", <-finalResult)
}

