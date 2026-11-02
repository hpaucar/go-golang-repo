/*
 * Archivo:     sequential.go
 *
 * Propósito:   Ilustrar la ejecución secuencial: `main` llama cinco veces a
 *              `doWork`, que simula una tarea de 1 segundo con `time.Sleep`.
 *              Cada tarea empieza solo cuando la anterior ha terminado. Sirve
 *              como punto de comparación con la versión concurrente
 *              (parallel.go).
 *
 * Compilación: go build sequential.go
 * Ejecución:   ./sequential
 *              (Alternativamente: go run sequential.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Para cada tarea, una línea "A: Work N started" y otra "B: Work N
 *              finished" con la hora. Las tareas aparecen en orden 0..4 y el
 *              programa tarda unos 5 segundos.
 *
 * Notas:
 *    1. El tiempo total es la suma de los tiempos de cada tarea, porque no hay
 *       solapamiento entre ellas.
 *    2. Compare la hora de inicio y fin de cada tarea con la salida de
 *       parallel.go para ver la diferencia entre ejecución secuencial y
 *       concurrente.
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
        doWork(i)
    }
}
