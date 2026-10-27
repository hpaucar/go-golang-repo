/*
 * Archivo:     channel_comma_ok.go
 *
 * Propósito:   Ilustrar el uso del idiom "comma-ok" (`msg, more := <-messages`) 
 *              en Go para detectar de forma explícita y segura si un canal 
 *              ha sido cerrado. Demuestra cómo una goroutine receptora puede 
 *              finalizar su ciclo de ejecución de manera ordenada tras procesar 
 *              todos los elementos enviados.
 *
 * Compilación: go build channel_comma_ok.go
 * Ejecución:   ./channel_comma_ok
 *              (Alternativamente: go run channel_comma_ok.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes con marcas de tiempo que muestran el envío y la 
 *              recepción de los valores 1, 2 y 3 acompañados de `true`. Tras 
 *              cerrar el canal, se recibe el valor cero acompañado de `false`, 
 *              lo que provoca que la función `receiver` termine limpiamente.
 *
 * Notas:       
 *    1. La sintaxis `msg, more := <-messages` devuelve el valor leído y un 
 *       booleano `more`. Si `more` es `true`, el canal sigue abierto o tiene 
 *       datos; si es `false`, el canal está cerrado y el valor recibido es 
 *       el valor cero del tipo de dato (`0` para `int`).
 *    2. La condición `if !more { return }` permite que el receptor reaccione 
 *       al cierre del canal y salga de su bucle infinito sin consumir 
 *       indefinidamente valores cero.
 *    3. Este enfoque es el equivalente manual y explícito al uso del bucle 
 *       `for range`, ideal cuando se requiere realizar lógica adicional justo 
 *       en el momento en que se detecta el cierre del canal.
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