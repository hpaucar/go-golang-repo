/*
 * Archivo:     selectmanychannels.go
 *
 * Propósito:   Ilustrar `select` para leer de varios canales a la vez. Dos
 *              goroutines escriben "Tick" cada segundo y "Tock" cada 3
 *              segundos; `main` imprime el mensaje del canal que esté listo.
 *
 * Compilación: go build selectmanychannels.go
 * Ejecución:   ./selectmanychannels
 *              (Alternativamente: go run selectmanychannels.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Tick" una vez por segundo y "Tock" cada 3 segundos, sin fin.
 *              Se detiene con Ctrl+C.
 *
 * Notas:
 *    1. `select` espera hasta que alguno de sus casos pueda ejecutarse; si
 *       hay varios listos, elige uno al azar.
 *    2. `writeEvery` devuelve un canal de solo lectura (`<-chan string`) y
 *       oculta la goroutine que escribe en él (patrón generador).
 *    3. El bucle de `main` no tiene condición de salida: el programa nunca
 *       termina solo.
 */

package main

import (
    "fmt"
    "time"
)

func writeEvery(msg string, seconds time.Duration) <-chan string {
    messages := make(chan string)
    go func() {
        for {
            time.Sleep(seconds)
            messages <- msg
        }
    }()
    return messages
}

func main() {
    messagesFromA := writeEvery("Tick", 1 * time.Second)
    messagesFromB := writeEvery("Tock", 3 * time.Second)
    for {
        select {
        case msg1 := <-messagesFromA:
            fmt.Println(msg1)
        case msg2 := <-messagesFromB:
            fmt.Println(msg2)
        }
    }
}
