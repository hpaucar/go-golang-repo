/*
 * Archivo:     selecttimer.go
 *
 * Propósito:   Ilustrar un tiempo límite (timeout) con `select` y
 *              `time.After`. Se espera un mensaje que llega a los 3 segundos;
 *              el tiempo de espera máximo se da en la línea de comandos.
 *
 * Compilación: go build selecttimer.go
 * Ejecución:   ./selecttimer 5
 *              (Alternativamente: go run selecttimer.go 5)
 *
 * Entrada:     Un argumento: los segundos de espera máxima (p. ej. 5 o 1).
 * Salida:      "Waiting for message for N seconds..." y luego "Message
 *              received: Hello" si N > 3, o "Timed out. Waited until:
 *              HH:MM:SS" si N < 3. Si N = 3, cualquiera de los dos.
 *
 * Notas:
 *    1. `time.After(d)` devuelve un canal que recibe la hora actual después
 *       de `d`; en un `select` funciona como tiempo límite.
 *    2. Sin argumento, `os.Args[1]` provoca un pánico (index out of range); y
 *       el error de `strconv.Atoi` se ignora (un texto no numérico vale 0).
 */

package main

import (
    "fmt"
    "os"
    "strconv"
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
    t, _ := strconv.Atoi(os.Args[1])
    messages := sendMsgAfter(3 * time.Second)
    timeoutDuration := time.Duration(t) * time.Second
    fmt.Printf("Waiting for message for %d seconds...\n", t)
    select {
    case msg := <-messages:
        fmt.Println("Message received:", msg)
    case tNow := <-time.After(timeoutDuration):
        fmt.Println("Timed out. Waited until:", tNow.Format("15:04:05"))
    }
}
