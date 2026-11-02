/*
 * Archivo:     catfiles.go
 *
 * Propósito:   Solución del ejercicio 2.1 del libro *Learn Concurrent
 *              Programming with Go*, que pide escribir una versión concurrente
 *              del comando `cat`: recibe nombres de archivo como argumentos y
 *              lanza una goroutine por archivo para leerlo con `os.ReadFile` e
 *              imprimir su contenido.
 *
 * Compilación: go build catfiles.go
 * Ejecución:   ./catfiles ../../commonfiles/txtfile1
 *              (Alternativamente: go run catfiles.go ../../commonfiles/txtfile1)
 *
 * Entrada:     Uno o más nombres de archivo como argumentos de línea de
 *              comandos (`os.Args[1:]`).
 * Salida:      El contenido de cada archivo. El orden entre archivos no está
 *              garantizado: depende de qué goroutine termine primero.
 *
 * Notas:
 *    1. `main` no espera a las goroutines de forma explícita: usa `time.Sleep(2
 *       * time.Second)` y confía en que terminen antes. Si un archivo tarda
 *       más, su contenido no se imprime.
 *    2. Si un archivo no existe, `log.Fatal` termina todo el programa,
 *       incluidas las demás goroutines.
 *    3. En capítulos posteriores se reemplaza el `Sleep` por mecanismos de
 *       sincronización como `sync.WaitGroup`.
 */

package main

import (
    "fmt"
    "log"
    "os"
    "time"
)

/*
  You can run this by executing:
  go run catfiles.go ../../commonfiles/txtfile1 ../../commonfiles/txtfile2 ../../commonfiles/txtfile3
*/

func printFile(filename string) {
    content, err := os.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(content))
}

func main() {
    filenames := os.Args[1:]
    for _, filename := range filenames {
        go printFile(filename)
    }
    time.Sleep(2 * time.Second)
}
