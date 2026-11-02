/*
 * Archivo:     ledgermutexarb.go
 *
 * Propósito:   Evitar el deadlock del libro contable con un árbitro. En
 *              lugar de un mutex por cuenta, el `Arbitrator` lleva un mapa
 *              de cuentas en uso protegido por un `sync.Cond`: un cajero
 *              reserva todas las cuentas que necesita de una sola vez o
 *              espera (`cond.Wait`) hasta que estén libres.
 *
 * Compilación: go build ledgermutexarb.go
 * Ejecución:   ./ledgermutexarb
 *              (Alternativamente: go run ledgermutexarb.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes "Locking" y "Unlocked" de los cuatro cajeros y,
 *              al final, "COMPLETE" para cada uno. El programa termina a los
 *              60 segundos aunque el trabajo acabe antes.
 *
 * Notas:
 *    1. Se rompe la condición de retención y espera (hold and wait): nadie
 *       retiene una cuenta mientras espera otra.
 *    2. `UnlockAccounts` usa `cond.Broadcast()` para despertar a todos los
 *       cajeros en espera, que vuelven a comprobar si sus cuentas quedaron
 *       libres.
 *    3. El costo es menos concurrencia: todas las reservas pasan por el
 *       mismo mutex del árbitro.
 */

package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func main() {
    accounts := []BankAccount{
        *NewBankAccount("Sam"),
        *NewBankAccount("Paul"),
        *NewBankAccount("Amy"),
        *NewBankAccount("Mia"),
    }
    total := len(accounts)
    arb := NewArbitrator()
    for i := 0; i < 4; i++ {
        go func(tellerId int) {
            for i := 1; i < 1000; i++ {
                from, to := rand.Intn(total), rand.Intn(total)
                for from == to {
                    to = rand.Intn(total)
                }
                accounts[from].Transfer(&accounts[to], 10, tellerId, arb)
            }
            fmt.Println(tellerId,"COMPLETE")
        }(i)
    }
    time.Sleep(60 * time.Second)
}

type BankAccount struct {
    id      string
    balance int
}

func NewBankAccount(id string) *BankAccount {
    return &BankAccount{
        id:      id,
        balance: 100,
    }
}

type Arbitrator struct {
    accountsInUse map[string]bool
    cond *sync.Cond
}

func NewArbitrator() *Arbitrator{
    return &Arbitrator{
        accountsInUse: map[string]bool{},
        cond:          sync.NewCond(&sync.Mutex{}),
    }
}

func (a *Arbitrator) LockAccounts(ids... string) {
    a.cond.L.Lock()
    for allAvailable := false; !allAvailable; {
        allAvailable = true
        for _, id := range ids {
            if a.accountsInUse[id] {
                allAvailable = false
                a.cond.Wait()
            }
        }
    }
    for _, id := range ids {
        a.accountsInUse[id] = true
    }
    a.cond.L.Unlock()
}

func (a *Arbitrator) UnlockAccounts(ids... string) {
    a.cond.L.Lock()
    for _, id := range ids {
        a.accountsInUse[id] = false
    }
    a.cond.Broadcast()
    a.cond.L.Unlock()
}

func (src *BankAccount) Transfer(to *BankAccount, amount int, tellerId int,
    arb *Arbitrator) {
    fmt.Printf("%d Locking %s and %s\n", tellerId, src.id, to.id)
    arb.LockAccounts(src.id, to.id)
    src.balance -= amount
    to.balance += amount
    arb.UnlockAccounts(src.id, to.id)
    fmt.Printf("%d Unlocked %s and %s\n", tellerId, src.id, to.id)
}

