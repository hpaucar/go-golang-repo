/*
 * Archivo:     filesearch.go
 *
 * Propósito:   Buscar archivos de forma recursiva y concurrente: `fileSearch`
 *              recorre un directorio, imprime las rutas cuyo nombre contiene el
 *              texto buscado y lanza una nueva goroutine por cada
 *              subdirectorio. Un `sync.WaitGroup` compartido crece con
 *              `wg.Add(1)` a medida que aparecen subdirectorios.
 *
 * Compilación: go build filesearch.go
 * Ejecución:   ./filesearch /ruta/directorio texto
 *              (Alternativamente: go run filesearch.go /ruta/directorio texto)
 *
 * Entrada:     Dos argumentos: el directorio inicial y el texto que debe
 *              contener el nombre del archivo, por ejemplo: ./filesearch
 *              /home/usuario README
 * Salida:      Las rutas completas de los archivos y directorios cuyo nombre
 *              contiene el texto, en orden no determinista.
 *
 * Notas:
 *    1. `wg.Add(1)` se llama antes de lanzar cada goroutine, así el contador
 *       nunca llega a 0 mientras quede trabajo.
 *    2. Si faltan argumentos, el programa termina con `panic` (índice fuera de
 *       rango en `os.Args`).
 *    3. Los errores de `os.ReadDir` se ignoran (por ejemplo, sin permisos); en
 *       árboles muy grandes se crean muchas goroutines.
 */

package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

func fileSearch(dir string, filename string, wg *sync.WaitGroup) {
    files, _ := os.ReadDir(dir)
    for _, file := range files {
        fpath := filepath.Join(dir, file.Name())
        if strings.Contains(file.Name(), filename) {
            fmt.Println(fpath)
        }
        if file.IsDir() {
            wg.Add(1)
            go fileSearch(fpath, filename, wg)
        }
    }
    wg.Done()
}

func main() {
    wg := sync.WaitGroup{}
    wg.Add(1)
    go fileSearch(os.Args[1], os.Args[2], &wg)
    wg.Wait()
}
