/*
 * Archivo:     grepfiles.go
 *
 * Propósito:   Solución del ejercicio 2.2 del libro *Learn Concurrent
 *              Programming with Go*, que pide una versión concurrente simple de
 *              `grep`: el primer argumento es la cadena a buscar y los
 *              siguientes son archivos; cada archivo se revisa en su propia
 *              goroutine con `strings.Contains`.
 *
 * Compilación: go build grepfiles.go
 * Ejecución:   ./grepfiles Mozilla ../../commonfiles/txtfile1
 *              (Alternativamente: go run grepfiles.go Mozilla ../../commonfiles/txtfile1)
 *
 * Entrada:     Argumentos: la cadena a buscar seguida de uno o más nombres de
 *              archivo.
 * Salida:      Una línea por archivo indicando si contiene ("contains a match")
 *              o no ("does NOT contain a match") la cadena. El orden de las
 *              líneas varía entre ejecuciones.
 *
 * Notas:
 *    1. Cada archivo se procesa en paralelo; la salida aparece en el orden en
 *       que terminan las goroutines, no en el orden de los argumentos.
 *    2. `main` espera con `time.Sleep(2 * time.Second)` en lugar de una
 *       sincronización real; un archivo lento podría quedar sin reportar.
 *    3. Sin argumentos, `os.Args[1]` provoca un pánico (index out of range).
 */

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

/*
  You can run this by executing:
  go run grepfiles.go Mozilla ../../commonfiles/txtfile1 ../../commonfiles/txtfile2 ../../commonfiles/txtfile3
*/

func grepFile(filename string, searchStr string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(string(content), searchStr) {
		fmt.Println(filename, "contains a match with", searchStr)
	} else {
		fmt.Println(filename, "does NOT contain a match with", searchStr)
	}
}

func main() {
	searchStr := os.Args[1]
	filenames := os.Args[2:]
	for _, filename := range filenames {
		go grepFile(filename, searchStr)
	}
	time.Sleep(2 * time.Second)
}
