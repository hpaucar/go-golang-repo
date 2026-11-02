/*
 * Archivo:     dirhashconcurrent.go
 *
 * Propósito:   Solución del ejercicio 10.1 del libro *Learn Concurrent
 *              Programming with Go*, que pide calcular el hash de un directorio
 *              en paralelo pero combinando los resultados en orden: una
 *              goroutine por archivo y una cadena de `sync.WaitGroup` que fija
 *              el orden de escritura en el MD5 final.
 *
 * Compilación: go build dirhashconcurrent.go
 * Ejecución:   ./dirhashconcurrent ../../commonfiles
 *              (Alternativamente: go run dirhashconcurrent.go ../../commonfiles)
 *
 * Entrada:     Argumento: la ruta de un directorio.
 * Salida:      "Processing <ruta>" por cada archivo y al final "<dir> - <hash
 *              MD5 en hexadecimal>".
 *
 * Notas:
 *    1. Actualmente no compila: importa `.../chapter10/listing10.1`, pero la
 *       carpeta se llama `listing10.01`.
 *    2. Cada goroutine espera a la anterior (`prev.Wait()`) antes de escribir y
 *       avisa a la siguiente (`next.Done()`): el cálculo es paralelo, la
 *       combinación es secuencial.
 *    3. Si el directorio no contiene archivos regulares, `next` queda en `nil`
 *       y `next.Wait()` provoca un pánico.
 */

package main

import (
    "crypto/md5"
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.1"
    "os"
    "path/filepath"
    "sync"
)

func main() {
    dir := os.Args[1]
    files, _ := os.ReadDir(dir)
    hMd5 := md5.New()
    var prev, next *sync.WaitGroup
    for _, file := range files {
        if !file.IsDir() {
            next = &sync.WaitGroup{}
            next.Add(1)
            go func(filename string, prev, next *sync.WaitGroup) {
                fpath := filepath.Join(dir, filename)
                fmt.Println("Processing", fpath)
                hashOnFile := listing10_1.FHash(fpath)
                // If not the first iteration
                if prev != nil {
                    prev.Wait()
                }
                hMd5.Write(hashOnFile)
                next.Done()
            }(file.Name(), prev, next)
            prev = next
        }
    }
    next.Wait()
    fmt.Printf("%s - %x\n", dir, hMd5.Sum(nil))
}
