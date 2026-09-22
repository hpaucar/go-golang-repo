/*
 * Simula depósitos y retiros concurrentes sobre un saldo compartido, 
 * cediendo la ejecución con runtime.Gosched() después de cada operación.
 * Ilustra que ceder la ejecución no evita la carrera de datos y que una 
 * espera fija no garantiza que ambas goroutines terminen.
 * "runtime.Gosched()": Puedo continuar, pero cedo mi turno para que otras tareas tengan oportunidad de avanzar
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
