/*
 * Archivo:     allfilesinfo.go
 *
 * Propósito:   Mostrar un deadlock con canales. Dos goroutines se comunican
 *              en ciclo: `handleDirectories` lee un directorio y envía sus
 *              entradas a `files`, y `handleFiles` imprime los datos de cada
 *              archivo y devuelve los subdirectorios a `dirs`. Ambos canales
 *              son sin búfer.
 *
 * Compilación: go build allfilesinfo.go
 * Ejecución:   ./allfilesinfo ./ruta/directorio
 *              (Alternativamente: go run allfilesinfo.go ./ruta/directorio)
 *
 * Entrada:     Un argumento de línea de comandos: el directorio raíz.
 * Salida:      Nombre, tamaño en KB y fecha de modificación de los archivos.
 *              Cuando aparece un subdirectorio, la salida se detiene: cada
 *              goroutine espera a que la otra reciba. El programa termina a
 *              los 60 segundos.
 *
 * Notas:
 *    1. `handleFiles` se bloquea al enviar a `dirs` mientras
 *       `handleDirectories` está bloqueada enviando a `files`: es una
 *       espera circular entre canales, igual que con mutex.
 *    2. El runtime no lo detecta porque `main` sigue en `time.Sleep`.
 *    3. Los archivos abiertos con `os.Open` nunca se cierran; en árboles
 *       grandes esto agota los descriptores de archivo.
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
            files <- filepath.Join(fullpath, file.Name())
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