/*
 * Archivo:     dirhashconcurrent.go
 *
 * Propósito:   Paralelizar el cálculo del hash de un directorio manteniendo
 *              el orden. Cada goroutine calcula en paralelo el hash de su
 *              archivo, pero antes de escribirlo en el hash común espera una
 *              señal de la goroutine anterior (canal `prev`) y, al terminar,
 *              avisa a la siguiente (canal `next`).
 *
 * Compilación: go build dirhashconcurrent.go
 * Ejecución:   ./dirhashconcurrent ./ruta/directorio
 *              (Alternativamente: go run dirhashconcurrent.go ./ruta/directorio)
 *
 * Entrada:     Un argumento de línea de comandos: la ruta de un directorio.
 *              Se procesan solo sus archivos (no los subdirectorios).
 * Salida:      El hash combinado del directorio en hexadecimal. Debe ser el
 *              mismo que imprime dirhashsequential.go (10.03).
 *
 * Notas:
 *    1. La parte costosa (`FHash`) se ejecuta en paralelo; solo la escritura
 *       en `sha` se hace en orden, como una cadena de relevos entre
 *       goroutines conectadas por canales sin búfer.
 *    2. `main` espera en `<-next` la señal de la última goroutine, que es la
 *       que garantiza que todas las escrituras ya terminaron.
 *    3. Si el directorio no tiene archivos, `next` queda en nil y `<-next`
 *       se bloquea para siempre: el runtime termina con "all goroutines are
 *       asleep - deadlock!".
 *    4. La ruta de importación `.../chapter10/listing10.1` no coincide con
 *       la carpeta real `listing10.01`; hay que corregirla para compilar.
 */

package main

import (
    "crypto/sha256"
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.1"
    "os"
    "path/filepath"
)

func main() {
    dir := os.Args[1]
    files, _ := os.ReadDir(dir)
    sha := sha256.New()
    var prev, next chan int
    for _, file := range files {
        if !file.IsDir() {
            next = make(chan int)
            go func(filename string, prev, next chan int) {
                fpath := filepath.Join(dir, filename)
                hashOnFile := listing10_1.FHash(fpath)
                if prev != nil {
                    <-prev
                }
                sha.Write(hashOnFile)
                next <- 0
            }(file.Name(), prev, next)
            prev = next
        }
    }
    <-next
    fmt.Printf("%x\n", sha.Sum(nil))
}
