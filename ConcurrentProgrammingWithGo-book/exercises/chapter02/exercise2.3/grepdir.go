/*
 * Archivo:     grepdir.go
 *
 * Propósito:   Solución del ejercicio 2.3 del libro *Learn Concurrent
 *              Programming with Go*, que modifica el grep concurrente del
 *              ejercicio 2.2 para recibir un directorio en lugar de una lista
 *              de archivos: se lista el directorio con `ioutil.ReadDir` y se
 *              lanza una goroutine por cada entrada.
 *
 * Compilación: go build grepdir.go
 * Ejecución:   ./grepdir Mozilla ../../commonfiles
 *              (Alternativamente: go run grepdir.go Mozilla ../../commonfiles)
 *
 * Entrada:     Argumentos: la cadena a buscar y la ruta de un directorio.
 * Salida:      Una línea por cada archivo (no subdirectorio) del directorio,
 *              indicando si contiene o no la cadena. Los subdirectorios se
 *              ignoran. El orden varía entre ejecuciones.
 *
 * Notas:
 *    1. La búsqueda no es recursiva: `grepPath` omite las entradas que son
 *       directorios.
 *    2. `ioutil.ReadDir` está obsoleto desde Go 1.16; el ejercicio 2.4 usa
 *       `os.ReadDir`, que es la alternativa recomendada.
 *    3. Igual que en los ejercicios anteriores, la espera se hace con
 *       `time.Sleep`, no con una sincronización real.
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
  You can run this by executing:
  go run grepdir.go Mozilla ../../commonfiles
*/

func grepPath(path string, fileInfo os.FileInfo, searchStr string) {
	fullPath := filepath.Join(path, fileInfo.Name())
	if !fileInfo.IsDir() {
		content, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(string(content), searchStr) {
			fmt.Println(fullPath, "contains a match with", searchStr)
		} else {
			fmt.Println(fullPath, "does NOT contain a match with", searchStr)
		}
	}
}

func main() {
	searchStr := os.Args[1]
	path := os.Args[2]
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range files {
		go grepPath(path, fileInfo, searchStr)
	}
	time.Sleep(2 * time.Second)
}
