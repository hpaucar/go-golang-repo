/*
 * Archivo:     stingyspendy.go
 *
 * Propósito:   Simular una cuenta bancaria con saldo inicial de 100 y dos
 *              goroutines: `stingy` hace un millón de depósitos de 10 y
 *              `spendy` un millón de retiros de 10. Ilustra una condición de
 *              carrera (race condition) sobre el saldo compartido.
 *
 * Compilación: go build stingyspendy.go
 * Ejecución:   ./stingyspendy
 *              (Alternativamente: go run stingyspendy.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Stingy Done", "Spendy Done" y el saldo final. El valor esperado
 *              es 100, pero suele salir otro número distinto en cada ejecución.
 *
 * Notas:
 *    1. `*money += 10` es una secuencia leer-modificar-escribir; si dos
 *       goroutines la intercalan, se pierden actualizaciones.
 *    2. Ejecute `go run -race stingyspendy.go` para que el detector de carreras
 *       lo muestre.
 *    3. La espera fija de 2 segundos con `time.Sleep` no garantiza que ambas
 *       goroutines hayan terminado.
 */

package main

import (
    "fmt"
    "time"
)

/*
Note: this program has a race condition for demonstration purposes
In later chapters we cover how to wait for threads to complete their work
*/
func stingy(money *int) {
    for i := 0; i < 1000000; i++ {
        *money += 10
    }
    fmt.Println("Stingy Done")
}

func spendy(money *int) {
    for i := 0; i < 1000000; i++ {
        *money -= 10
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
