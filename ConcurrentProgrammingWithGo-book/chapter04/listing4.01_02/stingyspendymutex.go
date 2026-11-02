/*
 * Archivo:     stingyspendymutex.go
 *
 * Propósito:   Corregir la carrera de datos de la cuenta bancaria protegiendo
 *              cada depósito, cada retiro y la lectura final del saldo con un
 *              `sync.Mutex`. Cada goroutine adquiere el candado (lock) solo
 *              durante la actualización.
 *
 * Compilación: go build stingyspendymutex.go
 * Ejecución:   ./stingyspendymutex
 *              (Alternativamente: go run stingyspendymutex.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Stingy Done", "Spendy Done" y "Money in bank account: 100" si
 *              ambas goroutines terminan en 2 segundos.
 *
 * Notas:
 *    1. `mutex.Lock()` y `mutex.Unlock()` delimitan la sección crítica: solo
 *       una goroutine a la vez modifica `money`.
 *    2. La lectura del saldo en `main` también se protege con el mutex para no
 *       leer un valor a medio actualizar.
 *    3. La espera fija de 2 segundos no garantiza que ambas goroutines
 *       terminen; en una máquina lenta el saldo mostrado podría no ser 100.
 */

package main

import (
    "fmt"
    "sync"
    "time"
)

func stingy(money *int, mutex *sync.Mutex) {
    for i := 0; i < 1000000; i++ {
        mutex.Lock()
        *money += 10
        mutex.Unlock()
    }
    fmt.Println("Stingy Done")
}

func spendy(money *int, mutex *sync.Mutex) {
    for i := 0; i < 1000000; i++ {
        mutex.Lock()
        *money -= 10
        mutex.Unlock()
    }
    fmt.Println("Spendy Done")
}

func main() {
    money := 100
    mutex := sync.Mutex{}
    go stingy(&money, &mutex)
    go spendy(&money, &mutex)
    time.Sleep(2 * time.Second)
    mutex.Lock()
    fmt.Println("Money in bank account: ", money)
    mutex.Unlock()
}
