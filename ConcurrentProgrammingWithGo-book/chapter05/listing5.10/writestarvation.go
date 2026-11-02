/*
 * Archivo:     writestarvation.go
 *
 * Propósito:   Mostrar la inanición de escritores (writer starvation) con el
 *              `ReadWriteMutex` del listado 4.12, que da preferencia a los
 *              lectores. Dos goroutines leen en bucle infinito con lecturas
 *              solapadas y `main` intenta obtener el candado de escritura.
 *
 * Compilación: go build writestarvation.go
 * Ejecución:   ./writestarvation
 *              (Alternativamente: go run writestarvation.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Líneas "Read done" sin fin. Como casi siempre hay un lector
 *              activo, el escritor no obtiene el candado y "Write finished"
 *              normalmente nunca aparece; detenga el programa con Ctrl+C.
 *
 * Notas:
 *    1. Un nuevo lector puede entrar mientras otro lector tiene el candado, así
 *       que el contador de lectores rara vez llega a 0.
 *    2. El escritor espera indefinidamente aunque el programa no esté en
 *       interbloqueo: es inanición, no deadlock.
 *    3. El listado 5.11 resuelve el problema con un mutex que da preferencia a
 *       los escritores.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter4/listing4.12"
    "time"
)

func main() {
    rwMutex := listing4_12.ReadWriteMutex{}
    for i := 0; i < 2; i++ {
        go func() {
            for {
                rwMutex.ReadLock()
                time.Sleep(1 * time.Second)
                fmt.Println("Read done")
                rwMutex.ReadUnlock()
            }
        }()
    }
    time.Sleep(1 * time.Second)
    rwMutex.WriteLock()
    fmt.Println("Write finished")
}
