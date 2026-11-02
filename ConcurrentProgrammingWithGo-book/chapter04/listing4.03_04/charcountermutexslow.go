/*
 * Archivo:     charcountermutexslow.go
 *
 * Propósito:   Contar las letras de los RFC 1000 a 1030 con una goroutine por
 *              documento y un `sync.Mutex` que protege el slice de frecuencias.
 *              Muestra una MALA ubicación del candado: se bloquea toda la
 *              función, incluida la descarga HTTP, y el programa pierde la
 *              concurrencia.
 *
 * Compilación: go build charcountermutexslow.go
 * Ejecución:   ./charcountermutexslow
 *              (Alternativamente: go run charcountermutexslow.go)
 *
 * Entrada:     Ninguna por teclado. Requiere conexión a Internet (31 peticiones
 *              HTTP).
 * Salida:      Líneas "Completed: <url> HH:MM:SS" que aparecen de una en una y,
 *              tras 60 segundos, las frecuencias "a-N ... z-N".
 *
 * Notas:
 *    1. Como el candado cubre `http.Get`, las descargas se ejecutan una tras
 *       otra: el resultado es correcto, pero tan lento como la versión
 *       secuencial.
 *    2. Regla práctica: mantenga la sección crítica lo más pequeña posible y no
 *       incluya E/S lenta dentro de ella. El listado 4.5 lo corrige.
 *    3. La espera fija de 60 segundos puede no bastar con una red lenta.
 */

package main

import (
    "fmt"
    "io"
    "net/http"
    "strings"
    "sync"
    "time"
)

const AllLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
    mutex := sync.Mutex{}
    var frequency = make([]int, 26)
    for i := 1000; i <= 1030; i++ {
        url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
        go CountLetters(url, frequency, &mutex)
    }
    time.Sleep(60 * time.Second)
    mutex.Lock()
    for i, c := range AllLetters {
        fmt.Printf("%c-%d ", c, frequency[i])
    }
    mutex.Unlock()
}

// CountLetters
// Note: this program us locking the entire goroutine with mutex on purpose to demonstrate
// bad placement of the lock and unlock. We fix this in the next listing
func CountLetters(url string, frequency []int, mutex *sync.Mutex) {
    mutex.Lock()
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        panic("Server returning error status code: " + resp.Status)
    }
    body, _ := io.ReadAll(resp.Body)
    for _, b := range body {
        c := strings.ToLower(string(b))
        cIndex := strings.Index(AllLetters, c)
        if cIndex >= 0 {
            frequency[cIndex] += 1
        }
    }
    fmt.Println("Completed:", url, time.Now().Format("15:04:05"))
    mutex.Unlock()
}
