/*
 * Archivo:     charcountersequential.go
 *
 * Propósito:   Descargar secuencialmente los documentos RFC 1000 a RFC 1030
 *              desde rfc-editor.org y contar las letras de la "a" a la "z" sin
 *              distinguir mayúsculas. Acumula las frecuencias de todos los
 *              documentos y muestra el total de apariciones de cada letra. Es
 *              la versión base, sin concurrencia.
 *
 * Compilación: go build charcountersequential.go
 * Ejecución:   ./charcountersequential
 *              (Alternativamente: go run charcountersequential.go)
 *
 * Entrada:     Ninguna por teclado. Requiere conexión a Internet (31 peticiones
 *              HTTP).
 * Salida:      Una línea "Completed: <url>" por documento, en orden, y al final
 *              las frecuencias en el formato "a-N b-N ... z-N".
 *
 * Notas:
 *    1. Las descargas se hacen una tras otra, así que el tiempo total es la
 *       suma de los tiempos de red.
 *    2. Si el servidor devuelve un código distinto de 200, el programa termina
 *       con `panic`.
 *    3. Los errores de `http.Get` se ignoran (`resp, _ :=`): si falla la red,
 *       `resp` es nil y el programa falla al usarlo.
 */

package main

import (
    "fmt"
    "io"
    "net/http"
    "strings"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency []int) {
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        panic("Server returning error status code: " + resp.Status)
    }
    body, _ := io.ReadAll(resp.Body)
    for _, b := range body {
        c := strings.ToLower(string(b))
        cIndex := strings.Index(allLetters, c)
        if cIndex >= 0 {
            frequency[cIndex] += 1
        }
    }
    fmt.Println("Completed:", url)
}

func main() {
    // Crea 26 contadores, inicializados automáticamente en cero.
    var frequency = make([]int, 26)
    for i := 1000; i <= 1030; i++ {
        url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
        countLetters(url, frequency)
    }
    for i, c := range allLetters {
        fmt.Printf("%c-%d ", c, frequency[i])
    }
}
