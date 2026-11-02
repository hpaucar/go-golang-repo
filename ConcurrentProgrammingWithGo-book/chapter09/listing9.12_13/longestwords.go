/*
 * Archivo:     longestwords.go
 *
 * Propósito:   Ilustrar una etapa final que necesita todos los datos antes de
 *              responder: tras el pipeline fan-out/fan-in de listing9.09_11,
 *              `longestWords` guarda las palabras únicas, las ordena por
 *              longitud y envía las 10 más largas.
 *
 * Compilación: go build longestwords.go
 * Ejecución:   ./longestwords
 *              (Alternativamente: go run longestwords.go)
 *
 * Entrada:     Red: descarga los RFC 100 a 130 desde
 *              https://rfc-editor.org/rfc/rfc<N>.txt.
 * Salida:      Una línea "Longest Words: ..." con las 10 palabras únicas más
 *              largas, separadas por comas.
 *
 * Notas:
 *    1. `longestWords` solo envía su resultado cuando su canal de entrada se
 *       cierra, es decir, cuando todo el pipeline terminó.
 *    2. `sort.Slice` no es estable: entre palabras de la misma longitud el
 *       orden puede cambiar en cada ejecución.
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
    "sort"
    "strings"
)

func longestWords(quit <-chan int, words <-chan string) <-chan string {
    longWords := make(chan string)
    go func() {
        defer close(longWords)
        uniqueWordsMap := make(map[string]bool)
        uniqueWords := make([]string, 0)
        moreData, word := true, ""
        for moreData {
            select {
            case word, moreData = <-words:
                if moreData && !uniqueWordsMap[word] {
                    uniqueWordsMap[word] = true
                    uniqueWords = append(uniqueWords, word)
                }
            case <-quit:
                return
            }
        }
        sort.Slice(uniqueWords, func(a, b int) bool {
            return len(uniqueWords[a]) > len(uniqueWords[b])
        })
        longWords <- strings.Join(uniqueWords[:10], ", ")
    }()
    return longWords
}

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
    results := longestWords(quit,
        extractWords(quit, listing9_10.FanIn(quit, pages...)))
    fmt.Println("Longest Words:", <-results)
}
