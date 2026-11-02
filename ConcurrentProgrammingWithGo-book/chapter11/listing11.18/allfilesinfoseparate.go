/*
 * Archivo:     allfilesinfoseparate.go
 *
 * Propósito:   Corregir el deadlock de 11.15_16_17 lanzando una goroutine
 *              por cada envío a `files`. Así `handleDirectories` nunca se
 *              bloquea al enviar y siempre vuelve a leer de `dirs`, lo que
 *              rompe la espera circular.
 *
 * Compilación: go build allfilesinfoseparate.go
 * Ejecución:   ./allfilesinfoseparate ./ruta/directorio
 *              (Alternativamente: go run allfilesinfoseparate.go
 *                                 ./ruta/directorio)
 *
 * Entrada:     Un argumento de línea de comandos: el directorio raíz.
 * Salida:      Nombre, tamaño en KB y fecha de modificación de todos los
 *              archivos del árbol, en un orden que cambia entre ejecuciones.
 *              El programa termina a los 60 segundos.
 *
 * Notas:
 *    1. La solución funciona, pero puede crear miles de goroutines en
 *       directorios grandes y no conserva el orden de los archivos.
 *    2. Los archivos abiertos con `os.Open` nunca se cierran.
 */

package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

func handleDirectories(dirs <-chan string, files chan<- string) {
    for fullpath := range dirs {
        fmt.Println("Reading all files from", fullpath)
        filesInDir, _ := os.ReadDir(fullpath)
        fmt.Printf("Pushing %d files from %s\n", len(filesInDir), fullpath)
        for _, file := range filesInDir {
            go func(fp string) {
                files <- fp
            }(filepath.Join(fullpath, file.Name()))
        }
    }
}

func handleFiles(files <-chan string, dirs chan<- string) {
    for path := range files {
        file, _ := os.Open(path)
        fileInfo, _ := file.Stat()
        if fileInfo.IsDir() {
            fmt.Printf("Pushing %s directory\n", fileInfo.Name())
            dirs <- path
        } else {
            fmt.Printf("File %s, size: %.2fKB, last modified: %s\n",
                fileInfo.Name(), float64(fileInfo.Size()) / 1024.0,
                fileInfo.ModTime().Format(time.ANSIC))
        }
    }
}

func main() {
    filesChannel := make(chan string)
    dirsChannel := make(chan string)
    go handleFiles(filesChannel, dirsChannel)
    go handleDirectories(dirsChannel, filesChannel)
    dirsChannel <- os.Args[1]
    time.Sleep(60 * time.Second)
}