/*
 * Archivo:     charcounterwaitgroup.go
 *
 * Propósito:   Contar las letras de los RFC 1000 a 1030 con una goroutine por
 *              documento (`listing4_5.CountLetters`) y esperar a que todas
 *              terminen con un `sync.WaitGroup`, en lugar de una espera fija.
 *
 * Compilación: go build charcounterwaitgroup.go
 * Ejecución:   ./charcounterwaitgroup
 *              (Alternativamente: go run charcounterwaitgroup.go)
 *
 * Entrada:     Ninguna por teclado. Requiere conexión a Internet (31 peticiones
 *              HTTP).
 * Salida:      Líneas "Completed: <url> HH:MM:SS" en orden no determinista y,
 *              cuando terminan las 31 descargas, las frecuencias "a-N ... z-N".
 *
 * Notas:
 *    1. `wg.Add(31)` coincide con el número de documentos (1000 a 1030).
 *    2. Cada goroutine anónima llama a `CountLetters` y luego a `wg.Done()`;
 *       `url` se declara dentro del bucle, así que cada goroutine recibe su
 *       propia copia.
 *    3. El programa termina en cuanto acaba la última descarga, sin esperar de
 *       más ni de menos.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter4/listing4.5"
    "sync"
)

func main() {
    wg := sync.WaitGroup{}
    wg.Add(31)
    mutex := sync.Mutex{}
    var frequency = make([]int, 26)
    for i := 1000; i <= 1030; i++ {
        url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
        go func() {
            listing4_5.CountLetters(url, frequency, &mutex)
            wg.Done()
        }()
    }
    wg.Wait()
    mutex.Lock()
    for i, c := range listing4_5.AllLetters {
        fmt.Printf("%c-%d ", c, frequency[i])
    }
    mutex.Unlock()
}
