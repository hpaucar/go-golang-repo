/*
 * Archivo:     stingyspendysignalon50.go
 *
 * Propósito:   Solución del ejercicio 5.1 del libro *Learn Concurrent
 *              Programming with Go*, que pide que `stingy` avise con
 *              `cond.Signal()` solo cuando el saldo llegue a 50 o más, en lugar
 *              de avisar en cada depósito. `spendy` espera con `cond.Wait()`
 *              hasta tener al menos 50 para retirar.
 *
 * Compilación: go build stingyspendysignalon50.go
 * Ejecución:   ./stingyspendysignalon50
 *              (Alternativamente: go run stingyspendysignalon50.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Stingy Done", "Spendy Done" y "Money in bank account: 100" si
 *              ambas goroutines terminan en menos de 2 segundos (10 millones
 *              depositados y retirados en total).
 *
 * Notas:
 *    1. Basta con `Signal` porque solo hay una goroutine esperando (`spendy`);
 *       con varias habría que usar `Broadcast`.
 *    2. `cond.Wait()` se llama dentro de un `for`, no de un `if`, porque la
 *       condición debe volver a verificarse al despertar.
 *    3. `main` lee el saldo después de `time.Sleep(2 * time.Second)`; en una
 *       máquina lenta puede imprimir un valor intermedio.
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
		if *money >= 50 {
			cond.Signal()
		}
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
