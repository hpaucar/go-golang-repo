/*
 * Archivo:     closingFlag.go
 *
 * Propósito:   Ilustrar la lectura con dos valores (comma-ok,
 *              `msg, more := <-messages`) para detectar si un canal fue
 *              cerrado, de modo que el receptor termine de forma ordenada.
 *
 * Compilación: go build closingFlag.go
 * Ejecución:   ./closingFlag
 *              (Alternativamente: go run closingFlag.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Líneas con hora (HH:MM:SS) con "Sending: i" y
 *              "Received: i true" para 1, 2 y 3. Tras el cierre se imprime
 *              "Received: 0 false" y `receiver` termina.
 *
 * Notas:
 *    1. `msg, more := <-messages` devuelve el valor y un booleano: `more` es
 *       `true` si el valor vino de un envío y `false` si el canal está
 *       cerrado y vacío (entonces `msg` es el valor cero, `0`).
 *    2. `if !more { return }` hace que el receptor salga del bucle al detectar
 *       el cierre, en lugar de leer ceros sin fin (compárese con
 *       listing7.08_09).
 *    3. Es la versión explícita de `for range` sobre un canal; útil cuando se
 *       necesita hacer algo justo al detectar el cierre.
 */

package main

import (
    "fmt"
    "time"
)

func receiver(messages <-chan int) {
    for {
        msg, more := <-messages
        fmt.Println(time.Now().Format("15:04:05"), "Received:", msg, more)
        time.Sleep(1 * time.Second)
        if !more {
            return
        }
    }
}

func main() {
    msgChannel := make(chan int)
    go receiver(msgChannel)
    
    for i := 1; i <= 3; i++ {
        fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
        msgChannel <- i
        time.Sleep(1 * time.Second)
    }
    
    close(msgChannel)
    time.Sleep(3 * time.Second)
}