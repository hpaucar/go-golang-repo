/*
 * Simula depósitos y retiros concurrentes sobre un saldo compartido, 
 * usando sync.Mutex para proteger sus actualizaciones y su lectura.
 * Evita la carrera de datos, pero la espera fija de dos segundos no 
 * garantiza que ambas goroutines terminen ni que el saldo mostrado sea 100.
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
