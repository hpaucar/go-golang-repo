/*
 * Archivo:     ledgermutex.go
 *
 * Propósito:   Mostrar un deadlock en un escenario realista: cuatro cajeros
 *              (goroutines) hacen hasta 999 transferencias aleatorias cada
 *              uno entre cuatro cuentas, usando `Transfer` de
 *              `bankaccountmutex.go` (11.03_04), que bloquea las cuentas en
 *              el orden origen-destino.
 *
 * Compilación: go build ledgermutex.go
 * Ejecución:   ./ledgermutex
 *              (Alternativamente: go run ledgermutex.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes "Locking" y "Unlocked" de cada cajero. En algún
 *              momento los mensajes se detienen porque dos o más cajeros se
 *              bloquean entre sí; casi nunca aparecen los cuatro "COMPLETE".
 *              El programa termina a los 60 segundos.
 *
 * Notas:
 *    1. El deadlock aparece cuando, por ejemplo, un cajero transfiere de
 *       Sam a Paul y otro de Paul a Sam al mismo tiempo (espera circular).
 *    2. El runtime no lo detecta porque `main` sigue en `time.Sleep`.
 *    3. La ruta de importación `.../chapter11/listing11.3_4` no coincide con
 *       la carpeta real `listing11.03_04`; hay que corregirla para compilar.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter11/listing11.3_4"
    "math/rand"
    "time"
)

func main() {
    accounts := []listing11_3_4.BankAccount{
        *listing11_3_4.NewBankAccount("Sam"),
        *listing11_3_4.NewBankAccount("Paul"),
        *listing11_3_4.NewBankAccount("Amy"),
        *listing11_3_4.NewBankAccount("Mia"),
    }
    total := len(accounts)
    for i := 0; i < 4; i++ {
        go func(eId int) {
            for j := 1; j < 1000; j++ {
                from, to := rand.Intn(total), rand.Intn(total)
                for from == to {
                    to = rand.Intn(total)
                }
                accounts[from].Transfer(&accounts[to], 10, eId)
            }
            fmt.Println(eId, "COMPLETE")
        }(i)
    }
    time.Sleep(60 * time.Second)
}
