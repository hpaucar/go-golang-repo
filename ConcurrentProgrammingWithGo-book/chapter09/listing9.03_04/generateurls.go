/*
 * Archivo:     generateurls.go
 *
 * Propósito:   Ilustrar la primera etapa de un pipeline: `generateUrls` es un
 *              generador que envía las URL de los RFC 100 a 130 por un canal,
 *              puede cancelarse con `quit` y cierra su canal al terminar.
 *
 * Compilación: go build generateurls.go
 * Ejecución:   ./generateurls
 *              (Alternativamente: go run generateurls.go)
 *
 * Entrada:     Ninguna (no descarga nada; solo genera las URL).
 * Salida:      31 líneas, de https://rfc-editor.org/rfc/rfc100.txt a
 *              https://rfc-editor.org/rfc/rfc130.txt.
 *
 * Notas:
 *    1. `defer close(urls)` cierra el canal al salir, lo que termina el
 *       `for range` de `main`.
 *    2. `select` con `case <-quit` permite detener al generador antes de
 *       tiempo; `defer close(quit)` en `main` lo avisa al terminar.
 *    3. Es el patrón que se reutiliza en los siguientes listados del
 *       capítulo.
 */

package main

import "fmt"

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
    results := generateUrls(quit)
    for result := range results {
        fmt.Println(result)
    }
}
