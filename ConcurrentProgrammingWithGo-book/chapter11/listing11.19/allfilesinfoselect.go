/*
 * Archivo:     allfilesinfoselect.go
 *
 * Propósito:   Corregir el deadlock de 11.15_16_17 con `select`. Ahora
 *              `handleDirectories` guarda en una cola (`toPush`) las rutas
 *              pendientes y, en cada iteración, hace lo que esté listo:
 *              recibir un nuevo directorio o enviar la siguiente ruta.
 *
 * Compilación: go build allfilesinfoselect.go
 * Ejecución:   ./allfilesinfoselect ./ruta/directorio
 *              (Alternativamente: go run allfilesinfoselect.go ./ruta/directorio)
 *
 * Entrada:     Un argumento de línea de comandos: el directorio raíz.
 * Salida:      Nombre, tamaño en KB y fecha de modificación de todos los
 *              archivos del árbol. El programa termina a los 60 segundos.
 *
 * Notas:
 *    1. Con `select` sobre una recepción y un envío, la goroutine nunca
 *       queda bloqueada en una sola operación, así que no hay espera
 *       circular.
 *    2. Si la cola está vacía, la goroutine solo espera un nuevo directorio
 *       (`<-dirs`); así nunca intenta enviar `toPush[0]` de una cola vacía.
 *    3. Los archivos abiertos con `os.Open` nunca se cierran.
 */

package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

func handleDirectories(dirs <-chan string, files chan<- string) {
    toPush := make([]string, 0)
    appendAllFiles := func(path string) {
        fmt.Println("Reading all files from", path)
        filesInDir, _ := os.ReadDir(path)
        fmt.Printf("Pushing %d files from %s\n", len(filesInDir), path)
        for _, f := range filesInDir {
            toPush = append(toPush, filepath.Join(path, f.Name()))
        }
    }
    for {
        if len(toPush) == 0 {
            appendAllFiles(<-dirs)
        } else {
            select {
            case fullpath := <-dirs:
                appendAllFiles(fullpath)
            case files <- toPush[0]:
                toPush = toPush[1:]
            }
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
