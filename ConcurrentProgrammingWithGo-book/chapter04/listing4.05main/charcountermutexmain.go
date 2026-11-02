/*
 * Archivo:     charcountermutexmain.go
 *
 * Propósito:   Contar las letras de los RFC 1000 a 1030 con una goroutine por
 *              documento, usando `listing4_5.CountLetters`, que bloquea el
 *              mutex solo durante el conteo. Las descargas se solapan y el
 *              programa es mucho más rápido que la versión con el candado mal
 *              ubicado.
 *
 * Compilación: go build charcountermutexmain.go
 * Ejecución:   ./charcountermutexmain
 *              (Alternativamente: go run charcountermutexmain.go)
 *
 * Entrada:     Ninguna por teclado. Requiere conexión a Internet (31 peticiones
 *              HTTP).
 * Salida:      Líneas "Completed: <url> HH:MM:SS" en orden no determinista y,
 *              tras 10 segundos, las frecuencias "a-N ... z-N".
 *
 * Notas:
 *    1. El mismo mutex protege la escritura en las goroutines y la lectura
 *       final en `main`.
 *    2. La espera fija de 10 segundos no garantiza que todas las descargas
 *       hayan terminado; en el capítulo 6 se usa `sync.WaitGroup`.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter4/listing4.5"
    "sync"
    "time"
)

func main() {
    mutex := sync.Mutex{}
    var frequency = make([]int, 26)
    for i := 1000; i <= 1030; i++ {
        url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
        go listing4_5.CountLetters(url, frequency, &mutex)
    }
    time.Sleep(10 * time.Second)
    mutex.Lock()
    for i, c := range listing4_5.AllLetters {
        fmt.Printf("%c-%d ", c, frequency[i])
    }
    mutex.Unlock()
}
