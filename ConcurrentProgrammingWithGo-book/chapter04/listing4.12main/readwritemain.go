/*
 * Archivo:     readwritemain.go
 *
 * Propósito:   Probar el `ReadWriteMutex` del listado 4.12: diez goroutines
 *              adquieren el candado de lectura a la vez y lo mantienen 5
 *              segundos. Un segundo después, `main` pide el candado de
 *              escritura y debe esperar a que terminen todos los lectores.
 *
 * Compilación: go build readwritemain.go
 * Ejecución:   ./readwritemain
 *              (Alternativamente: go run readwritemain.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Diez "Read started", luego "Write started"; unos 4 segundos
 *              después, diez "Read done" y finalmente "Write finished".
 *
 * Notas:
 *    1. Los lectores trabajan en paralelo, lo que demuestra que el candado de
 *       lectura es compartido.
 *    2. El escritor queda bloqueado hasta que el último lector llama a
 *       `ReadUnlock()`.
 *    3. `main` no llama a `WriteUnlock()` porque el programa termina justo
 *       después.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter4/listing4.12"
    "time"
)

func main() {
    rwMutex := listing4_12.ReadWriteMutex{}
    for i := 0; i < 10; i++ {
        go func() {
            rwMutex.ReadLock()
            fmt.Println("Read started")
            time.Sleep(5 * time.Second)
            fmt.Println("Read done")
            rwMutex.ReadUnlock()
        }()
    }
    time.Sleep(1 * time.Second)
    fmt.Println("Write started")
    rwMutex.WriteLock()
    fmt.Println("Write finished")
}
