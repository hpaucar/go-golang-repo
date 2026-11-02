/*
 * Archivo:     goscheduler.go
 *
 * Propósito:   Mostrar el efecto de `runtime.Gosched()`: `main` lanza la
 *              goroutine `sayHello` y cede su turno al planificador para que
 *              esta tenga la oportunidad de ejecutarse antes de imprimir
 *              "Finished". Es un ejemplo de lo que NO se debe hacer para
 *              sincronizar.
 *
 * Compilación: go build goscheduler.go
 * Ejecución:   ./goscheduler
 *              (Alternativamente: go run goscheduler.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Normalmente "Hello" y luego "Finished", pero el orden no está
 *              garantizado: "Hello" puede salir después o no salir, porque el
 *              programa termina al acabar `main`.
 *
 * Notas:
 *    1. `runtime.Gosched()` solo cede el procesador; no espera a que otra
 *       goroutine termine.
 *    2. Con varios núcleos, `sayHello` puede ejecutarse en paralelo, antes o
 *       después de `Gosched()`; el resultado depende del planificador.
 *    3. Para esperar a una goroutine se deben usar mecanismos de sincronización
 *       (`sync.WaitGroup`, canales), no el planificador.
 */

package main

import (
    "fmt"
    "runtime"
)
/*
Note: this program is an example of what not to do; using go scheduler
to synchronize executions
*/
func sayHello() {
    fmt.Println("Hello")
}

func main() {
    go sayHello()
    
    // Cede voluntariamente la ejecución para permitir que
    // el planificador ejecute otras goroutines.
    // main permanece lista para reanudar su ejecución.
    // Esto NO equivale a esperar la finalización de sayHello.
    runtime.Gosched()

    // No existe una sincronización que garantice que "Hello"
    // se imprima antes de "Finished".
    fmt.Println("Finished")

    // Al terminar main, finaliza el programa.
    // Si sayHello todavía no imprimió su mensaje, este podría no aparecer.
}
