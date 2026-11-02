/*
 * Archivo:     selectwithnil.go
 *
 * Propósito:   Ilustrar cómo desactivar un caso de `select` poniendo su canal
 *              en `nil`. Dos generadores envían 50 ventas y 40 gastos
 *              aleatorios; cuando un canal se cierra se le asigna `nil` y
 *              `select` ya no lo atiende.
 *
 * Compilación: go build selectwithnil.go
 * Ejecución:   ./selectwithnil
 *              (Alternativamente: go run selectwithnil.go)
 *
 * Entrada:     Ninguna (montos aleatorios entre 1 y 100).
 * Salida:      Líneas "Sale of: x" y "Expense of: y" mezcladas (unos 5
 *              segundos) y al final "End of day profit and loss: total".
 *
 * Notas:
 *    1. Un canal cerrado siempre está listo para leer; si no se desactivara,
 *       `select` lo elegiría una y otra vez con el valor cero.
 *    2. Un canal `nil` nunca está listo, así que `sales = nil` quita ese caso
 *       del `select` sin afectar al otro.
 *    3. El bucle termina cuando ambos canales son `nil`
 *       (`for sales != nil || expenses != nil`).
 */

package main

import (
    "fmt"
    "math/rand"
    "time"
)

func generateAmounts(n int) <-chan int {
    amounts := make(chan int)
    go func() {
        defer close(amounts)
        for i := 0; i < n; i++ {
            amounts <- rand.Intn(100) + 1
            time.Sleep(100 * time.Millisecond)
        }
    }()
    return amounts
}

func main() {
    sales := generateAmounts(50)
    expenses := generateAmounts(40)
    endOfDayAmount := 0
    for sales != nil || expenses != nil {
        select {
        case sale, moreData := <-sales:
            if moreData {
                fmt.Println("Sale of:", sale)
                endOfDayAmount += sale
            } else {
                sales = nil
            }
        case expense, moreData := <-expenses:
            if moreData {
                fmt.Println("Expense of:", expense)
                endOfDayAmount -= expense
            } else {
                expenses = nil
            }
        }
    }
    fmt.Println("End of day profit and loss:", endOfDayAmount)
}
