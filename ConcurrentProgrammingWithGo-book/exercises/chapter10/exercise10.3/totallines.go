/*
 * Archivo:     totallines.go
 *
 * Propósito:   Código de partida del ejercicio 10.3 del libro *Learn Concurrent
 *              Programming with Go*, que pide volver concurrente este programa
 *              con el patrón fan-out/fan-in. Esta versión secuencial descarga
 *              los RFC 1000 a 1029 uno por uno y suma sus líneas.
 *
 * Compilación: go build totallines.go
 * Ejecución:   ./totallines
 *              (Alternativamente: go run totallines.go)
 *
 * Entrada:     Ninguna por teclado; descarga los documentos desde
 *              https://rfc-editor.org (requiere conexión a Internet).
 * Salida:      "Downloading <url>" por cada RFC y al final "Total lines: <n>".
 *
 * Notas:
 *    1. Como las descargas son secuenciales, el tiempo total es la suma de las
 *       30 descargas.
 *    2. La solución concurrente está en
 *       `../exercise10.3solution/totallines.go`.
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
    totalLines := 0
    for i := 1000; i < 1000 + pagesToDownload; i++ {
        url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
        fmt.Println("Downloading", url)
        resp, _ := http.Get(url)
        if resp.StatusCode != 200 {
            panic("Server's error: " + resp.Status)
        }
        bodyBytes, _ := io.ReadAll(resp.Body)
        totalLines += strings.Count(string(bodyBytes), "\n")
        resp.Body.Close()
    }
    fmt.Println("Total lines:", totalLines)
}

