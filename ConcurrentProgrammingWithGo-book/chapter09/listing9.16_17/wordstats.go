/*
 * Archivo:     wordstats.go
 *
 * Propósito:   Ilustrar el patrón broadcast en un pipeline: las palabras de
 *              los RFC 100 a 130 (fan-out/fan-in) se copian con
 *              `listing9_14.Broadcast` a dos etapas, `longestWords` y
 *              `frequentWords`, que trabajan en paralelo.
 *
 * Compilación: go build wordstats.go
 * Ejecución:   ./wordstats
 *              (Alternativamente: go run wordstats.go)
 *
 * Entrada:     Red: descarga los RFC 100 a 130 desde
 *              https://rfc-editor.org/rfc/rfc<N>.txt.
 * Salida:      "Longest Words: ..." con las 10 palabras más largas y "Most
 *              frequent Words: ..." con las 10 más frecuentes.
 *
 * Notas:
 *    1. Gracias a `Broadcast`, el mismo flujo de palabras alimenta dos
 *       cálculos distintos sin descargar las páginas dos veces.
 *    2. `frequentWords` cuenta las apariciones con un mapa y ordena la lista
 *       de palabras por frecuencia.
 *    3. La ruta de importación usa `chapter9/...`, pero la carpeta real es
 *       `chapter09/`; tal como está, `go build` falla con "no required module
 *       provides package". Hay que corregir la ruta para compilar.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter9/listing9.10"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter9/listing9.14"
    "io"
    "net/http"
    "regexp"
    "sort"
    "strings"
)

func frequentWords(quit <-chan int, words <-chan string) <-chan string {
    mostFrequentWords := make(chan string)
    go func() {
        defer close(mostFrequentWords)
        freqMap := make(map[string]int)
        freqList := make([]string, 0)
        moreData, word := true, ""
        for moreData {
            select {
            case word, moreData = <-words:
                if moreData {
                    if freqMap[word] == 0 {
                        freqList = append(freqList, word)
                    }
                    freqMap[word] += 1
                }
            case <-quit:
                return
            }
        }
        sort.Slice(freqList, func(a, b int) bool {
            return freqMap[freqList[a]] > freqMap[freqList[b]]
        })
        mostFrequentWords <- strings.Join(freqList[:10], ", ")
    }()
    return mostFrequentWords
}

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
    words := extractWords(quit, listing9_10.FanIn(quit, pages...))
    wordsMulti := listing9_14.Broadcast(quit, words, 2)
    longestResults := longestWords(quit, wordsMulti[0])
    frequentResults := frequentWords(quit, wordsMulti[1])
    fmt.Println("Longest Words:", <-longestResults)
    fmt.Println("Most frequent Words:", <-frequentResults)
}
