/*
 * Archivo:     dirfilehash.go
 *
 * Propósito:   Ilustrar el patrón de bucle paralelo (fork-join): por cada
 *              archivo regular de un directorio se lanza una goroutine que
 *              calcula su hash SHA-256 con `FHash`, y un `sync.WaitGroup`
 *              espera a que todas terminen.
 *
 * Compilación: go build dirfilehash.go
 * Ejecución:   ./dirfilehash ./ruta/directorio
 *              (Alternativamente: go run dirfilehash.go ./ruta/directorio)
 *
 * Entrada:     Un argumento de línea de comandos: la ruta de un directorio.
 *              Se procesan solo sus archivos (no los subdirectorios).
 * Salida:      Una línea por archivo con el formato `nombre - hash`. El
 *              orden de las líneas cambia entre ejecuciones, porque depende
 *              de qué goroutine termina primero.
 *
 * Notas:
 *    1. El nombre del archivo se pasa como parámetro a la goroutine, para
 *       que cada una trabaje con su propia copia del valor.
 *    2. `wg.Add(1)` se llama antes de lanzar la goroutine y `wg.Done()` al
 *       final de ella; `wg.Wait()` bloquea `main` hasta que todas terminan.
 *    3. La ruta de importación `.../chapter10/listing10.1` no coincide con
 *       la carpeta real `listing10.01`; hay que corregir la importación (o
 *       renombrar la carpeta) para que el programa compile.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.1"
    "os"
    "path/filepath"
    "sync"
)

func main() {
    dir := os.Args[1]
    files, _ := os.ReadDir(dir)
    wg := sync.WaitGroup{}
    for _, file := range files {
        if !file.IsDir() {
            wg.Add(1)
            go func(filename string) {
                fPath := filepath.Join(dir, filename)
                hash := listing10_1.FHash(fPath)
                fmt.Printf("%s - %x\n", filename, hash)
                wg.Done()
            }(file.Name())
        }
    }
    wg.Wait()
}
