/*
 * Archivo:     closing.go
 *
 * Propósito:   Ilustrar el cierre de un canal con `close()` y lo que ve un
 *              receptor que lee de un canal cerrado sin comprobarlo: no se
 *              bloquea y recibe el valor cero del tipo.
 *
 * Compilación: go build closing.go
 * Ejecución:   ./closing
 *              (Alternativamente: go run closing.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Líneas con hora (HH:MM:SS) con "Sending: i" y "Received: i"
 *              para 1, 2 y 3. Después del cierre, el receptor imprime
 *              "Received: 0" una vez por segundo (unas 3 veces) hasta que
 *              `main` termina.
 *
 * Notas:
 *    1. `close(msgChannel)` indica que no se enviarán más valores. Leer de un
 *       canal cerrado nunca bloquea: devuelve de inmediato el valor cero del
 *       tipo (`0` para `int`).
 *    2. `receiver` hace `msg := <-messages` en un bucle infinito, por eso no
 *       distingue un 0 real del cierre del canal.
 *    3. Para detectar el cierre se usa `msg, ok := <-messages` o
 *       `for msg := range messages` (ver listing7.10 y listing7.11).
 *    4. El `time.Sleep(3 * time.Second)` final solo mantiene vivo el programa
 *       para observar este efecto.
 */

package main

import (
    "fmt"
    "time"
)

func receiver(messages <-chan int) {
    for {
        msg := <-messages
        fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
        time.Sleep(1 * time.Second)
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