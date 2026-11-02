/*
 * Archivo:     extractwords.go
 *
 * Propósito:   Ilustrar un pipeline de tres etapas: generar URL de los RFC 100
 *              a 130, descargar las páginas y extraer sus palabras
 *              (`extractWords`, con la expresión regular `[a-zA-Z]+`).
 *
 * Compilación: go build extractwords.go
 * Ejecución:   ./extractwords
 *              (Alternativamente: go run extractwords.go)
 *
 * Entrada:     Red: descarga los RFC 100 a 130 desde
 *              https://rfc-editor.org/rfc/rfc<N>.txt.
 * Salida:      Todas las palabras de los documentos en minúsculas, una por
 *              línea (muchas líneas).
 *
 * Notas:
 *    1. Cada etapa añadida solo necesita un canal de entrada, uno de salida
 *       y el canal `quit`; las etapas se componen como funciones.
 *    2. Las etapas trabajan en paralelo: mientras se descarga una página, se
 *       extraen las palabras de la anterior.
 *    3. Igual que en listing9.05_06, sin red o con un código distinto de 200
 *       el programa termina con un pánico.
 */

package main

import (
    "fmt"
    "io"
    "net/http"
    "regexp"
    "strings"
)

func extractWords(quit <-chan int, pages <-chan string) <-chan string {
    words := make(chan string)
    go func() {
        defer close(words)
        wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
        moreData, pg := true, ""
        for moreData {
            select {
            case pg, moreData = <-pages:
                if moreData {
                    for _, word := range wordRegex.FindAllString(pg, -1) {
                        words <- strings.ToLower(word)
                    }
                }
            case <-quit:
                return
            }
        }
    }()
    return words
}

func downloadPages(quit <-chan int, urls <-chan string) <-chan string {
    pages := make(chan string)
    go func() {
        defer close(pages)
        moreData, url := true, ""
        for moreData {
            select {
            case url, moreData = <-urls:
                if moreData {
                    resp, _ := http.Get(url)
                    if resp.StatusCode != 200 {
                        panic("Server's error: " + resp.Status)
                    }
                    body, _ := io.ReadAll(resp.Body)
                    pages <- string(body)
                    resp.Body.Close()
                }
            case <-quit:
                return
            }
        }
    }()
    return pages
}

func generateUrls(quit <-chan int) <-chan string {
    urls := make(chan string)
    go func() {
        defer close(urls)
        for i := 100; i <= 130; i++ {
            url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
            select {
            case urls <- url:
            case <-quit:
                return
            }
        }
    }()
    return urls
}

func main() {
    quit := make(chan int)
    defer close(quit)
    results := extractWords(quit, downloadPages(quit, generateUrls(quit)))
    for result := range results {
        fmt.Println(result)
    }
}
