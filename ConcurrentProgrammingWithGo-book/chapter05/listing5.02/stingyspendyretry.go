/*
 * Archivo:     stingyspendyretry.go
 *
 * Propósito:   Evitar el saldo negativo haciendo que `spendy` espere a tener al
 *              menos 50: libera el mutex, duerme 10 ms y vuelve a comprobar
 *              (sondeo, polling). Ilustra una solución que funciona pero
 *              desperdicia tiempo y CPU.
 *
 * Compilación: go build stingyspendyretry.go
 * Ejecución:   ./stingyspendyretry
 *              (Alternativamente: go run stingyspendyretry.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Stingy Done", "Spendy Done" y "Money in bank account: 100" si
 *              ambas goroutines terminan en 2 segundos. El saldo nunca llega a
 *              ser negativo.
 *
 * Notas:
 *    1. Es necesario liberar el mutex antes de dormir; si no, `stingy` nunca
 *       podría depositar y habría un bloqueo.
 *    2. El sondeo introduce retardos de hasta 10 ms aunque el dinero ya esté
 *       disponible. `sync.Cond` (listado 5.3) evita este problema.
 *    3. La espera fija de 2 segundos en `main` sigue sin garantizar la
 *       finalización.
 */

package main

import (
    "fmt"
    "os"
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
    for i := 0; i < 200000; i++ {
        mutex.Lock()
        for *money < 50 {
            mutex.Unlock()
            time.Sleep(10 * time.Millisecond)
            mutex.Lock()
        }
        *money -= 50
        if *money < 0 {
            fmt.Println("Money is negative!")
            os.Exit(1)
        }
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
