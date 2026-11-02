/*
 * Archivo:     grepdirrec.go
 *
 * Propósito:   Solución del ejercicio 2.4 del libro *Learn Concurrent
 *              Programming with Go*, que extiende el grep concurrente para
 *              recorrer el directorio de forma recursiva: cada subdirectorio
 *              encontrado genera nuevas goroutines (una por entrada) que buscan
 *              la cadena en sus archivos.
 *
 * Compilación: go build grepdirrec.go
 * Ejecución:   ./grepdirrec Mozilla ../../commonfiles
 *              (Alternativamente: go run grepdirrec.go Mozilla ../../commonfiles)
 *
 * Entrada:     Argumentos: la cadena a buscar y la ruta del directorio raíz.
 * Salida:      Una línea por cada archivo del árbol de directorios, indicando
 *              si contiene o no la cadena. El orden varía entre ejecuciones.
 *
 * Notas:
 *    1. El número de goroutines crece con el tamaño del árbol: cada
 *       subdirectorio lanza una goroutine por cada entrada que contiene.
 *    2. `main` espera solo 2 segundos con `time.Sleep`; en árboles grandes
 *       algunos archivos pueden no llegar a reportarse.
 *    3. Las goroutines se crean de forma recursiva, un patrón que en el
 *       capítulo 6 se coordina con `sync.WaitGroup`.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
  You can run this by executing:
  go run grepdirrec.go Mozilla ../../commonfiles
*/

func grepPath(path string, dirEntry os.DirEntry, searchStr string) {
	fullPath := filepath.Join(path, dirEntry.Name())
	if dirEntry.IsDir() {
		files, err := os.ReadDir(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			go grepPath(fullPath, file, searchStr)
		}
	} else {
		content, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(string(content), searchStr) {
			fmt.Println(fullPath, "contains a match with ", searchStr)
		} else {
			fmt.Println(fullPath, "does NOT contain a match with", searchStr)
		}
	}
}

func main() {
	searchStr := os.Args[1]
	path := os.Args[2]
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, dirEntry := range files {
		go grepPath(path, dirEntry, searchStr)
	}
	time.Sleep(2 * time.Second)
}
