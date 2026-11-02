/*
 * Archivo:     bufferchannel.go
 *
 * Propósito:   Ilustrar los canales con búfer (buffered channels): `main`
 *              envía los números 1 a 6 por un canal de capacidad 3 a un
 *              receptor lento (1 segundo por lectura) y termina con el valor
 *              centinela `-1`. Un `sync.WaitGroup` espera al receptor.
 *
 * Compilación: go build bufferchannel.go
 * Ejecución:   ./bufferchannel
 *              (Alternativamente: go run bufferchannel.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Líneas con hora (HH:MM:SS): "Sending: i. Buffer Size: n" con el
 *              número de elementos en el búfer antes de cada envío, y
 *              "Received: i" una vez por segundo, hasta "Received: -1".
 *
 * Notas:
 *    1. `make(chan int, 3)` crea un canal con búfer de capacidad 3: `main`
 *       puede enviar sin bloquearse mientras haya espacio; cuando el búfer se
 *       llena, el envío espera a que el receptor saque un elemento.
 *    2. `len(msgChannel)` devuelve cuántos elementos hay en el búfer en ese
 *       instante; así se observa cómo se llena.
 *    3. El valor centinela `-1` indica a `receiver` que no hay más datos;
 *       `wGroup.Wait()` evita que `main` termine antes de que el receptor
 *       llame a `Done()`.
 */

package main

import (
    "fmt"
    "sync"
    "time"
)

func receiver(messages chan int, wGroup *sync.WaitGroup) {
    msg := 0
    for msg != -1 {
        time.Sleep(1 * time.Second)
        msg = <-messages
        fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
    }
    wGroup.Done()
}

func main() {
    msgChannel := make(chan int, 3)
    wGroup := sync.WaitGroup{}
    wGroup.Add(1)
    
    go receiver(msgChannel, &wGroup)
    
    for i := 1; i <= 6; i++ {
        size := len(msgChannel)
        fmt.Printf("%s Sending: %d. Buffer Size: %d\n",
            time.Now().Format("15:04:05"), i, size)
        msgChannel <- i
    }
    
    msgChannel <- -1
    wGroup.Wait()
}