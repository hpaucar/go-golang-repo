/*
 * Archivo:     cond_fix_demonstration.go
 *
 * Propósito:   Demostrar cómo solucionar (técnicamente) el problema de la 
 *              "señal perdida" del ejemplo anterior mediante el uso correcto 
 *              del mutex asociado a la variable de condición. Este programa 
 *              fuerza un orden de ejecución estricto entre el hilo principal 
 *              y la goroutine hija.
 *
 * Compilación: go build cond_fix_demonstration.go
 * Ejecución:   ./cond_fix_demonstration
 *              (Alternativamente: go run cond_fix_demonstration.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El programa imprimirá la secuencia de inicio y finalización 
 *              de 50,000 iteraciones de forma exitosa sin caer en interbloqueo.
 *
 * Notas:       
 *    1. El hilo principal adquiere el candado (`cond.L.Lock()`) ANTES de 
 *       lanzar la goroutine hija. Por lo tanto, cuando `doWork` intenta 
 *       adquirir el mismo candado, se bloquea (se queda esperando).
 *    2. `doWork` solo puede avanzar, adquirir el candado y ejecutar `Signal()` 
 *       una vez que el hilo principal llama a `cond.Wait()`, ya que esta 
 *       función libera internamente el candado antes de suspender al hilo.
 *    3. Esto garantiza que `cond.Signal()` jamás se ejecutará antes de que 
 *       `main` esté escuchando en `cond.Wait()`, evitando que la señal se pierda.
 *    4. ADVERTENCIA: Aunque esto evita el deadlock en este diseño específico, 
 *       sigue siendo una práctica "frágil". El estándar de la industria y la 
 *       forma más robusta de usar `sync.Cond` es acompañarlo SIEMPRE de una 
 *       variable de estado explícita (ej. `workDone boolean`) y comprobar esa 
 *       variable dentro de un bucle `for` antes del `Wait()`. Esto protege 
 *       contra los "despertares espurios" (spurious wakeups).
 */
 
package main

import (
    "fmt"
    "sync"
)

func doWork(cond *sync.Cond) {
    fmt.Println("Work started")
    fmt.Println("Work finished")
    cond.L.Lock()
    cond.Signal()
    cond.L.Unlock()
}

func main() {
    cond := sync.NewCond(&sync.Mutex{})
    cond.L.Lock()
    for i := 0; i < 50000; i++ {
        go doWork(cond)
        fmt.Println("Waiting for child goroutine")
        cond.Wait()
        fmt.Println("Child goroutine finished")
    }
    cond.L.Unlock()
}
