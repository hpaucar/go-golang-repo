/*
 * Archivo:     ledgermutexorder.go
 *
 * Propósito:   Corregir el deadlock del libro contable (11.05) aplicando un
 *              orden fijo de bloqueo. `Transfer` ordena las dos cuentas por
 *              su `id` con `sort.Slice` y bloquea siempre primero la de
 *              menor `id`, sin importar cuál es el origen o el destino.
 *
 * Compilación: go build ledgermutexorder.go
 * Ejecución:   ./ledgermutexorder
 *              (Alternativamente: go run ledgermutexorder.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes "Locking" y "Unlocked" de los cuatro cajeros y un
 *              "COMPLETE" por cada uno. El programa termina a los 60
 *              segundos aunque el trabajo acabe antes.
 *
 * Notas:
 *    1. Con un orden global (aquí, alfabético por nombre) no puede formarse
 *       un ciclo de espera entre cajeros.
 *    2. El orden de desbloqueo no importa para evitar el deadlock; solo
 *       importa el orden de bloqueo.
 *    3. Esta técnica exige conocer de antemano todos los recursos que se van
 *       a bloquear, lo que no siempre es posible.
 */

package main

import (
    "fmt"
    "math/rand"
    "sort"
    "sync"
    "time"
)

type BankAccount struct {
    id      string
    balance int
    mutex   sync.Mutex
}

func NewBankAccount(id string) *BankAccount {
    return &BankAccount{
        id:      id,
        balance: 100,
        mutex:   sync.Mutex{},
    }
}

func (src *BankAccount) Transfer(to *BankAccount, amount int, tellerId int) {
    accounts := []*BankAccount{src, to}
    sort.Slice(accounts, func(a, b int) bool {
        return accounts[a].id < accounts[b].id
    })
    fmt.Printf("%d Locking %s's account\n", tellerId, accounts[0].id)
    accounts[0].mutex.Lock()
    fmt.Printf("%d Locking %s's account\n", tellerId, accounts[1].id)
    accounts[1].mutex.Lock()
    src.balance -= amount
    to.balance += amount
    to.mutex.Unlock()
    src.mutex.Unlock()
    fmt.Printf("%d Unlocked %s and %s\n", tellerId, src.id, to.id)
}

func main() {
    accounts := []BankAccount{
        *NewBankAccount("Sam"),
        *NewBankAccount("Paul"),
        *NewBankAccount("Amy"),
        *NewBankAccount("Mia"),
    }
    for i := 0; i < 4; i++ {
        go func(tellerId int) {
            for i := 1; i < 1000; i++ {
                from, to := rand.Intn(len(accounts)), rand.Intn(len(accounts))
                for from == to {
                    to = rand.Intn(len(accounts))
                }
                accounts[from].Transfer(&accounts[to], 10, tellerId)
            }
            fmt.Println(tellerId,"COMPLETE")
        }(i)
    }
    time.Sleep(60 * time.Second)
}