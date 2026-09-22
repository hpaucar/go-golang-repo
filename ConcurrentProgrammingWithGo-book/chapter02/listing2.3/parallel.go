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
