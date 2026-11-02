/*
 * Archivo:     nonblockingmutex.go
 *
 * Propósito:   Ilustrar el bloqueo no bloqueante con `mutex.TryLock()`: se
 *              lanzan 201 goroutines que cuentan letras de los RFC 2000 a 2200
 *              y, mientras tanto, `main` intenta cada 100 ms leer los
 *              contadores. Si el candado está libre, imprime el progreso; si
 *              no, sigue sin esperar.
 *
 * Compilación: go build nonblockingmutex.go
 * Ejecución:   ./nonblockingmutex
 *              (Alternativamente: go run nonblockingmutex.go)
 *
 * Entrada:     Ninguna por teclado. Requiere conexión a Internet (201
 *              peticiones HTTP).
 * Salida:      Durante unos 10 segundos (100 intentos), líneas de frecuencias
 *              parciales "a-N ... z-N" o el mensaje "Mutex already being used",
 *              mezcladas con líneas "Completed: <url>".
 *
 * Notas:
 *    1. `TryLock()` (Go 1.18 o superior) devuelve `true` si obtuvo el candado y
 *       `false` si está ocupado, sin bloquear a la goroutine.
 *    2. Es útil para tareas opcionales, como mostrar progreso, que no deben
 *       esperar. La documentación de Go advierte que su uso correcto es raro.
 *    3. El programa termina tras los 100 intentos aunque queden descargas
 *       pendientes.
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
    for i := 2000; i <= 2200; i++ {
        url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
        go listing4_5.CountLetters(url, frequency, &mutex)
    }
    for i := 0; i < 100; i++ {
        time.Sleep(100 * time.Millisecond)
        if mutex.TryLock() {
            for i, c := range listing4_5.AllLetters {
                fmt.Printf("%c-%d ", c, frequency[i])
            }
            fmt.Println()
            mutex.Unlock()
        } else {
            fmt.Println("Mutex already being used")
        }
    }
}
