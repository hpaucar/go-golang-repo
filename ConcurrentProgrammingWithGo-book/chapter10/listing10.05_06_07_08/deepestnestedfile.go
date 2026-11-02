/*
 * Archivo:     deepestnestedfile.go
 *
 * Propósito:   Ilustrar el patrón fork-join con un canal de resultados. Se
 *              recorre un árbol de directorios con `filepath.Walk`; por cada
 *              archivo `.go` se lanza (fork) una goroutine que calcula la
 *              profundidad máxima de llaves `{}`. La goroutine de
 *              `joinResults` (join) reúne los resultados y elige el mayor.
 *
 * Compilación: go build deepestnestedfile.go
 * Ejecución:   ./deepestnestedfile ./ruta/directorio
 *              (Alternativamente: go run deepestnestedfile.go ./ruta/directorio)
 *
 * Entrada:     Un argumento de línea de comandos: la ruta del directorio
 *              raíz. Se recorren también todos sus subdirectorios.
 * Salida:      Una línea con el archivo `.go` que tiene el bloque de código
 *              más anidado y su nivel de anidamiento.
 *
 * Notas:
 *    1. `main` espera con `wg.Wait()` a que terminen todas las goroutines
 *       de trabajo y solo entonces cierra `partialResults`; el cierre hace
 *       que el `for range` de `joinResults` termine y envíe el máximo.
 *    2. `joinResults` debe estar en marcha antes de `wg.Wait()`: el canal
 *       no tiene búfer y, sin un lector, las goroutines de trabajo quedarían
 *       bloqueadas en el envío y `wg.Wait()` nunca volvería.
 *    3. El conteo de llaves es aproximado: también cuenta las `{` y `}` que
 *       aparecen dentro de cadenas o comentarios.
 */

package main

import (
    "fmt"
    "math"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

type CodeDepth struct {
    file  string
    level int
}

func deepestNestedBlock(filename string) CodeDepth {
    code, _ := os.ReadFile(filename)
    max := 0
    level := 0
    for _, c := range code {
        if c == '{' {
            level += 1
            max = int(math.Max(float64(max), float64(level)))
        } else if c == '}' {
            level -= 1
        }
    }
    return CodeDepth{filename, max}
}

func forkIfNeeded(path string, info os.FileInfo,
    wg *sync.WaitGroup, results chan CodeDepth) {
    if !info.IsDir() && strings.HasSuffix(path, ".go") {
        wg.Add(1)
        go func() {
            results <- deepestNestedBlock(path)
            wg.Done()
        }()
    }
}

func joinResults(partialResults chan CodeDepth) chan CodeDepth {
    finalResult := make(chan CodeDepth)
    max := CodeDepth{"", 0}
    go func() {
        for pr := range partialResults {
            if pr.level > max.level {
                max = pr
            }
        }
        finalResult <- max
    }()
    return finalResult
}

func main() {
    dir := os.Args[1]
    partialResults := make(chan CodeDepth)
    wg := sync.WaitGroup{}
    filepath.Walk(dir,
        func(path string, info os.FileInfo, err error) error {
            forkIfNeeded(path, info, &wg, partialResults)
            return nil
        })
    finalResult := joinResults(partialResults)
    wg.Wait()
    close(partialResults)
    result := <-finalResult
    fmt.Printf("%s has the deepest nested code block of %d\n",
        result.file, result.level)
}
