/*
 * Archivo:     nonblocking.go
 *
 * Propósito:   Ilustrar `select` con `default` para hacer una lectura no
 *              bloqueante: mientras no llega el mensaje (enviado después de 3
 *              segundos), `main` hace otra cosa y vuelve a intentarlo.
 *
 * Compilación: go build nonblocking.go
 * Ejecución:   ./nonblocking
 *              (Alternativamente: go run nonblocking.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "No messages waiting" una vez por segundo (unas 3 veces) y
 *              luego "Message received: Hello"; el programa termina.
 *
 * Notas:
 *    1. Si ningún caso de `select` está listo, se ejecuta `default` de
 *       inmediato, sin bloquear.
 *    2. Es útil para hacer trabajo mientras se espera un dato, pero sin el
 *       `time.Sleep` se convertiría en espera activa (busy waiting) y
 *       consumiría CPU.
 */

package main

import (
    "fmt"
    "time"
)

func sendMsgAfter(seconds time.Duration) <-chan string {
    messages := make(chan string)
    go func() {
        time.Sleep(seconds)
        messages <- "Hello"
    }()
    return messages
}

func main() {
    messages := sendMsgAfter(3 * time.Second)
    for {
        select {
        case msg := <-messages:
            fmt.Println("Message received:", msg)
            return
        default:
            fmt.Println("No messages waiting")
            time.Sleep(1 * time.Second)
        }
    }
}
