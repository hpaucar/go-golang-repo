/*
 * Archivo:     downloadpages.go
 *
 * Propósito:   Ilustrar un pipeline de dos etapas: `generateUrls` produce las
 *              URL de los RFC 100 a 130 y `downloadPages` descarga cada una y
 *              envía su texto por otro canal. Ambas etapas respetan `quit`.
 *
 * Compilación: go build downloadpages.go
 * Ejecución:   ./downloadpages
 *              (Alternativamente: go run downloadpages.go)
 *
 * Entrada:     Red: descarga los RFC 100 a 130 desde
 *              https://rfc-editor.org/rfc/rfc<N>.txt.
 * Salida:      El texto completo de los 31 documentos, en orden.
 *
 * Notas:
 *    1. Cada etapa es una goroutine que lee un canal, procesa y escribe en
 *       otro; cuando su entrada se cierra, cierra su salida (`defer close`).
 *    2. `url, moreData = <-urls` dentro del `select` detecta el cierre del
 *       canal de entrada.
 *    3. El error de `http.Get` se ignora: sin red, `resp` es `nil` y el
 *       programa termina con un pánico. Si algún documento no devuelve 200,
 *       también hay pánico ("Server's error: ...").
 */

package main

import (
    "fmt"
    "io"
    "net/http"
)

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
    results := downloadPages(quit, generateUrls(quit))
    for result := range results {
        fmt.Println(result)
    }
}
