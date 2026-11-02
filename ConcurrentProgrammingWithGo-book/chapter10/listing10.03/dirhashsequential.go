/*
 * Archivo:     dirhashsequential.go
 *
 * Propósito:   Calcular un único hash SHA-256 para todo un directorio de
 *              forma secuencial: se obtiene el hash de cada archivo con
 *              `FHash` y se escribe, en orden, en un hash acumulado. Es la
 *              versión de referencia para la versión concurrente (10.04).
 *
 * Compilación: go build dirhashsequential.go
 * Ejecución:   ./dirhashsequential ./ruta/directorio
 *              (Alternativamente: go run dirhashsequential.go ./ruta/directorio)
 *
 * Entrada:     Un argumento de línea de comandos: la ruta de un directorio.
 *              Se procesan solo sus archivos (no los subdirectorios).
 * Salida:      Una sola línea `directorio - hash` con el hash combinado.
 *
 * Notas:
 *    1. El resultado depende del orden en que se escriben los hashes;
 *       `os.ReadDir` devuelve las entradas ordenadas por nombre, así que el
 *       hash es siempre el mismo para el mismo contenido.
 *    2. Este orden obligatorio es lo que impide paralelizar el bucle de forma
 *       ingenua: hay una dependencia entre iteraciones (loop-carried
 *       dependence) en la variable `sha`.
 *    3. La ruta de importación `.../chapter10/listing10.1` no coincide con
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
    for _, file := range files {
        if !file.IsDir() {
            fpath := filepath.Join(dir, file.Name())
            hashOnFile := listing10_1.FHash(fpath)
            sha.Write(hashOnFile)
        }
    }
    fmt.Printf("%s - %x\n", dir, sha.Sum(nil))
}
