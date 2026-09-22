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
