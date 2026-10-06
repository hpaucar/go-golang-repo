/*
 * Archivo:     cond_bug_demonstration.go
 *
 * Propósito:   Demostrar un error de concurrencia clásico al utilizar variables 
 *              de condición: el problema de la "señal perdida" (lost signal). 
 *              El programa pretende lanzar goroutines hijas y esperar a que 
 *              terminen usando sync.Cond, pero falla por un diseño incorrecto.
 *
 * Compilación: go build cond_bug_demonstration.go
 * Ejecución:   ./cond_bug_demonstration
 *              (Alternativamente: go run cond_bug_demonstration.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El programa imprimirá algunos mensajes de inicio y luego 
 *              entrará en un estado de interbloqueo (deadlock). El runtime 
 *              de Go probablemente detectará esto y arrojará el error:
 *              "fatal error: all goroutines are asleep - deadlock!".
 *
 * Notas:       
 *    1. El error principal radica en que no se utiliza una variable de estado 
 *       junto con la variable de condición. `sync.Cond` no tiene "memoria"; 
 *       si la goroutine `doWork` ejecuta `cond.Signal()` ANTES de que el hilo 
 *       principal alcance `cond.Wait()`, la señal se pierde para siempre.
 *    2. La regla de oro de las variables de condición es que `Wait()` siempre 
 *       debe estar envuelto en un bucle `for` que verifique una condición 
 *       lógica (ej. `for !workDone { cond.Wait() }`). Aquí, esa comprobación 
 *       brilla por su ausencia.
 *    3. Este código es solo para fines demostrativos de lo que *no* se debe 
 *       hacer. Un `sync.WaitGroup` o canales (channels) serían las herramientas 
 *       idóneas en Go para esperar a que terminen múltiples goroutines.
 */
 
package main

import (
    "fmt"
    "sync"
)

/*
Note: this program has a bug for demonstration purposes
We demonstrate how to fix this problem in the next listing
*/
func doWork(cond *sync.Cond) {
    fmt.Println("Work started")
    fmt.Println("Work finished")
    cond.Signal()
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
