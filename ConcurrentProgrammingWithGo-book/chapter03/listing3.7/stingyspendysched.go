/*
 * Archivo:     stingyspendysched.go
 *
 * Propósito:   Repetir la simulación de depósitos y retiros concurrentes, pero
 *              llamando a `runtime.Gosched()` después de cada operación.
 *              Ilustra que ceder el turno al planificador no evita la condición
 *              de carrera, solo cambia la forma en que se intercalan las
 *              goroutines.
 *
 * Compilación: go build stingyspendysched.go
 * Ejecución:   ./stingyspendysched
 *              (Alternativamente: go run stingyspendysched.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El saldo final tras 2 segundos, normalmente distinto de 100.
 *              Como cada iteración cede el procesador, las goroutines van más
 *              lentas y puede que "Stingy Done" o "Spendy Done" no aparezcan
 *              antes de que termine `main`.
 *
 * Notas:
 *    1. `runtime.Gosched()` significa "puedo continuar, pero cedo mi turno para
 *       que otras goroutines avancen". No ofrece exclusión mutua.
 *    2. Con varios núcleos las dos goroutines se ejecutan en paralelo, así que
 *       la carrera sigue presente; `go run -race` la detecta.
 *    3. La espera fija de 2 segundos no garantiza la finalización de las
 *       goroutines.
 */

package main

import (
    "fmt"
    "runtime"
    "time"
)

/*
Note: this program has a race condition for demonstration purposes
In later chapters we cover how to wait for threads to complete their work
*/
func stingy(money *int) {
    for i := 0; i < 1000000; i++ {
        *money += 10
        runtime.Gosched()
    }
    fmt.Println("Stingy Done")
}

func spendy(money *int) {
    for i := 0; i < 1000000; i++ {
        *money -= 10
        runtime.Gosched()
    }
    fmt.Println("Spendy Done")
}

func main() {
    money := 100
    go stingy(&money)
    go spendy(&money)
    time.Sleep(2 * time.Second)
    fmt.Println("Money in bank account: ", money)
}
