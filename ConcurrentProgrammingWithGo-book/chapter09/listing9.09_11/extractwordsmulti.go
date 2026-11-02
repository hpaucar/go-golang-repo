/*
 * Archivo:     extractwordsmulti.go
 *
 * Propósito:   Ilustrar el patrón fan-out/fan-in: 20 goroutines
 *              `downloadPages` leen del mismo canal de URL (fan-out) y
 *              `listing9_10.FanIn` une sus salidas en un solo canal (fan-in)
 *              antes de extraer las palabras.
 *
 * Compilación: go build extractwordsmulti.go
 * Ejecución:   ./extractwordsmulti
 *              (Alternativamente: go run extractwordsmulti.go)
 *
 * Entrada:     Red: descarga los RFC 100 a 130 desde
 *              https://rfc-editor.org/rfc/rfc<N>.txt.
 * Salida:      Todas las palabras de los documentos en minúsculas, una por
 *              línea; el orden de los documentos varía en cada ejecución.
 *
 * Notas:
 *    1. Varias goroutines pueden leer del mismo canal; cada URL la recibe
 *       solo una de ellas, así las descargas se reparten.
 *    2. `FanIn` (listing9.10) une varios canales en uno y lo cierra cuando
 *       todos los de entrada se cierran.
 *    3. La ruta de importación usa `chapter9/...`, pero la carpeta real es
 *       `chapter09/`; tal como está, `go build` falla con "no required module
 *       provides package". Hay que corregir la ruta para compilar.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter9/listing9.10"
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
        moreData, page := true, ""
        for moreData {
            select {
            case page, moreData = <-pages:
                if moreData {
                    for _, word := range wordRegex.FindAllString(page, -1) {
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

const downloaders = 20

func main() {
    quit := make(chan int)
    defer close(quit)
    urls := generateUrls(quit)
    pages := make([]<-chan string, downloaders)
    for i := 0; i < downloaders; i++ {
        pages[i] = downloadPages(quit, urls)
    }
    results := extractWords(quit, listing9_10.FanIn(quit, pages...))
    for result := range results {
        fmt.Println(result)
    }
}
