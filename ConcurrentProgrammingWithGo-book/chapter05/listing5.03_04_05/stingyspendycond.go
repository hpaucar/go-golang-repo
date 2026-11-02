/*
 * Archivo:     stingyspendycond.go
 *
 * Propósito:   Ilustrar las variables de condición (`sync.Cond`) junto con un
 *              mutex. Esta versión elimina el sondeo del listado anterior:
 *              `spendy` se suspende con `cond.Wait()` hasta que `stingy`
 *              deposita fondos y emite `cond.Signal()`.
 *
 * Compilación: go build stingyspendycond.go
 * Ejecución:   ./stingyspendycond
 *              (Alternativamente: go run stingyspendycond.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Stingy Done", "Spendy Done" y el saldo final (100 si ambas
 *              goroutines terminan en 2 segundos). Si hubiera un error lógico,
 *              se imprimiría "Money is negative!".
 *
 * Notas:
 *    1. `cond.Wait()` libera el candado y suspende la goroutine; al despertar,
 *       vuelve a adquirir el candado antes de continuar. Así no se gasta CPU
 *       esperando.
 *    2. `cond.Signal()` despierta a una goroutine que espera en `cond.Wait()`;
 *       si nadie espera, la señal simplemente se ignora.
 *    3. La comprobación `for *money < 50` alrededor de `Wait()` es obligatoria:
 *       un depósito de 10 no siempre basta y la condición debe revisarse al
 *       despertar.
 *    4. `main` sigue usando `time.Sleep` para esperar; `sync.WaitGroup` sería
 *       la forma idiomática.
 */

package main

import (
    "fmt"
    "os"
    "sync"
    "time"
)

func main() {
    money := 100
    mutex := sync.Mutex{}
    cond := sync.NewCond(&mutex)
    go stingy(&money, cond)
    go spendy(&money, cond)
    time.Sleep(2 * time.Second)
    mutex.Lock()
    fmt.Println("Money in bank account: ", money)
    mutex.Unlock()
}

func stingy(money *int, cond *sync.Cond) {
    for i := 0; i < 1000000; i++ {
        cond.L.Lock()
        *money += 10
        cond.Signal()
        cond.L.Unlock()
    }
    fmt.Println("Stingy Done")
}

func spendy(money *int, cond *sync.Cond) {
    for i := 0; i < 200000; i++ {
        cond.L.Lock()
        for *money < 50 {
            cond.Wait()
        }
        *money -= 50
        if *money < 0 {
            fmt.Println("Money is negative!")
            os.Exit(1)
        }
        cond.L.Unlock()
    }
    fmt.Println("Spendy Done")
}
