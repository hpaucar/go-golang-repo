/*
 * Archivo:     stingyspendynegative.go
 *
 * Propósito:   Mostrar que el mutex no basta cuando hay una condición de
 *              negocio: `stingy` deposita 10 un millón de veces y `spendy`
 *              retira 50 doscientas mil veces sin comprobar si hay saldo. Si el
 *              saldo queda negativo, el programa termina.
 *
 * Compilación: go build stingyspendynegative.go
 * Ejecución:   ./stingyspendynegative
 *              (Alternativamente: go run stingyspendynegative.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Normalmente "Money is negative!" en muy poco tiempo y el
 *              programa termina con código 1 (`os.Exit(1)`).
 *
 * Notas:
 *    1. El mutex evita la carrera de datos, pero no hace que `spendy` espere a
 *       que haya dinero suficiente.
 *    2. `os.Exit(1)` termina el proceso de inmediato, aunque el mutex siga
 *       bloqueado y otras goroutines estén activas.
 *    3. Los listados siguientes resuelven la espera: primero con sondeo
 *       (polling) y luego con variables de condición.
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
