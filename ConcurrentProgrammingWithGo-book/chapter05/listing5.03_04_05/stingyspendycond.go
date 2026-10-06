/*
 * Archivo:     stingy_spendy_cond.go
 *
 * Propósito:   Ilustrar la sincronización avanzada de memoria compartida en
 *              Go utilizando Variables de Condición (sync.Cond) junto con 
 *              Mutexes. Este programa mejora una versión anterior eliminando 
 *              la "espera activa" (busy-waiting). La goroutine "spendy" 
 *              ahora se suspende de forma eficiente hasta que "stingy" 
 *              deposita fondos suficientes y emite una señal.
 *
 * Compilación: go build stingy_spendy_cond.go
 * Ejecución:   ./stingy_spendy_cond
 *              (Alternativamente: go run stingy_spendy_cond.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes de confirmación cuando cada goroutine termina su 
 *              ejecución y el saldo final de la cuenta bancaria. Si ocurre 
 *              un error lógico, imprimirá "Money is negative!".
 *
 * Notas:       
 *    1. `cond.Wait()` libera automáticamente el candado (lock) subyacente y 
 *       suspende la ejecución de la goroutine "spendy". Al ser despertada, 
 *       vuelve a adquirir el candado antes de continuar. Esto es mucho más 
 *       eficiente para el uso de la CPU que el bucle con `time.Sleep`.
 *    2. `cond.Signal()` se utiliza en la función "stingy" para notificar y 
 *       despertar a una goroutine que esté esperando en `cond.Wait()`.
 *    3. La comprobación `for *money < 50` alrededor del `Wait()` previene 
 *       despertares espurios (spurious wakeups), lo cual es una práctica 
 *       estándar y crítica al usar variables de condición.
 *    4. Al igual que en la versión anterior, la función `main` utiliza 
 *       `time.Sleep` de forma rudimentaria para esperar la finalización. 
 *       Un `sync.WaitGroup` sería el enfoque idiomático adecuado en Go.
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
