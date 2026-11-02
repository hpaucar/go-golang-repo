/*
 * Archivo:     parallel.go
 *
 * Propósito:   Ilustrar la ejecución concurrente con goroutines: `main` lanza
 *              cinco llamadas a `doWork` con la palabra clave `go`, de modo que
 *              las tareas de 1 segundo se solapan en el tiempo. Luego `main`
 *              duerme 2 segundos para dar tiempo a que terminen.
 *
 * Compilación: go build parallel.go
 * Ejecución:   ./parallel
 *              (Alternativamente: go run parallel.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Las cinco líneas "A: Work N started" aparecen casi a la misma
 *              hora y en orden no determinista; un segundo después aparecen las
 *              líneas "B: Work N finished". El programa tarda unos 2 segundos.
 *
 * Notas:
 *    1. `go doWork(i)` crea una goroutine; `main` continúa sin esperar. El
 *       argumento `i` se evalúa en el momento de lanzar la goroutine.
 *    2. El orden de ejecución de las goroutines lo decide el planificador
 *       (scheduler) de Go, por eso puede cambiar en cada ejecución.
 *    3. Esperar con `time.Sleep` no es una sincronización real: si una tarea
 *       tardara más de 2 segundos, el programa terminaría antes. En el capítulo
 *       6 se usa `sync.WaitGroup` para esperar correctamente.
 */

package main

import (
    "fmt"
    "time"
)

func doWork(id int) {
    fmt.Printf("A: Work %d started at %s\n", id, time.Now().Format("15:04:05"))
    time.Sleep(1 * time.Second)
    fmt.Printf("B: Work %d finished at %s\n", id, time.Now().Format("15:04:05"))
}

func main() {
    for i := 0; i < 5; i++ {
        // "go" ejecuta doWork en una nueva goroutine.
        // El argumento i se evalúa al lanzar la llamada.
        // main continúa sin esperar a que termine la tarea.
        go doWork(i)
    }
    time.Sleep(2 * time.Second)
}
