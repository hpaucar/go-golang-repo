/*
 * Archivo:     orderedfilesearch.go
 *
 * Propósito:   Solución del ejercicio 6.1 del libro *Learn Concurrent
 *              Programming with Go*, que pide modificar la búsqueda recursiva
 *              de archivos para que, en lugar de imprimir cada coincidencia al
 *              encontrarla, las acumule en una lista protegida por un mutex y
 *              las imprima ordenadas al final, tras `wg.Wait()`.
 *
 * Compilación: go build orderedfilesearch.go
 * Ejecución:   ./orderedfilesearch ../../commonfiles txtfile
 *              (Alternativamente: go run orderedfilesearch.go ../../commonfiles txtfile)
 *
 * Entrada:     Argumentos: el directorio raíz y la cadena que debe aparecer en
 *              el nombre del archivo o directorio.
 * Salida:      Las rutas completas que coinciden, una por línea y en orden
 *              alfabético.
 *
 * Notas:
 *    1. `sync.WaitGroup` permite esperar a un número de goroutines que no se
 *       conoce de antemano: se llama a `wg.Add(1)` antes de lanzar cada
 *       subbúsqueda.
 *    2. El slice `matches` se comparte entre goroutines, por eso cada `append`
 *       va dentro de `mutex.Lock()`/`Unlock()`.
 *    3. Los errores de `os.ReadDir` se ignoran: un directorio sin permisos
 *       simplemente no aporta resultados.
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func fileSearch(dir string, filename string,
	wg *sync.WaitGroup, mutex *sync.Mutex, matches *[]string) {
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		fpath := filepath.Join(dir, file.Name())
		if strings.Contains(file.Name(), filename) {
			mutex.Lock()
			*matches = append(*matches, fpath)
			mutex.Unlock()
		}
		if file.IsDir() {
			wg.Add(1)
			go fileSearch(fpath, filename, wg, mutex, matches)
		}
	}
	wg.Done()
}

func main() {
	results := make([]string, 0)
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go fileSearch(os.Args[1], os.Args[2], &wg, &mutex, &results)
	wg.Wait()
	mutex.Lock()
	sort.Strings(results)
	fmt.Println(strings.Join(results, "\n"))
	mutex.Unlock()
}
