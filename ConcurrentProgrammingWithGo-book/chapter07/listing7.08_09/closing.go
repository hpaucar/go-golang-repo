








/*
 * Archivo:     unbuffered_channel_close.go
 *
 * Propósito:   Ilustrar el cierre de un canal sin búfer en Go mediante la 
 *              función `close()` y analizar el comportamiento de una goroutine 
 *              receptora cuando lee de un canal cerrado sin utilizar una 
 *              comprobación de estado (comma-ok) o un bucle `range`.
 *
 * Compilación: go build unbuffered_channel_close.go
 * Ejecución:   ./unbuffered_channel_close
 *              (Alternativamente: go run unbuffered_channel_close.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes con marcas de tiempo que muestran el envío y la 
 *              recepción secuencial de los enteros 1, 2 y 3. Luego de que 
 *              el hilo principal cierra el canal, el receptor continuará 
 *              imprimiendo el valor cero (`Received: 0`) indefinidamente.
 *
 * Notas:       
 *    1. `close(msgChannel)` indica formalmente que no se enviarán más valores 
 *       a través del canal. Las operaciones de lectura en un canal cerrado 
 *       nunca se bloquean; en su lugar, devuelven de inmediato el valor cero 
 *       del tipo de datos subyacente (`0` para `int`).
 *    2. En este ejemplo, la función `receiver` realiza una lectura simple 
 *       `msg := <-messages` dentro de un bucle infinito `for`, por lo que 
 *       no puede detectar cuándo el canal ha sido cerrado.
 *    3. Para un código de producción robusto, se debe verificar el cierre 
 *       utilizando el idiom `msg, ok := <-messages` o recorrer el canal 
 *       directamente con un bucle `for msg := range messages`.
 *    4. La llamada a `time.Sleep(3 * time.Second)` al final de `main` otorga 
 *       tiempo suficiente para que la goroutine receptora procese los datos 
 *       antes de que el programa finalice.
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